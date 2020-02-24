package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

/**
Supporting models & utility code
**/

// RateDescriptor describes a rate in requests per second and duration.
type RateDescriptor struct {
	Rate     uint          `json:"rate"`
	Duration time.Duration `json:"duration"`
}

func (rs RateDescriptor) String() string {
	return fmt.Sprintf("%v req/s for %v", rs.Rate, rs.Duration)
}

// AttackDescriptor describes an attack by name and series of rates.
type AttackDescriptor struct {
	Name  string           `json:"name"`
	Rates []RateDescriptor `json:"rates"`
}

// Duration returns the aggregate time in an attack by summing all the duration of the rates.
func (attack AttackDescriptor) Duration() time.Duration {
	duration := time.Second
	for _, rate := range attack.Rates {
		duration += rate.Duration
	}
	return duration
}

// dynamicPacer defines the abstract interface for pacers the can build attack plans from CSV files or strings.
type dynamicPacer interface {
	vegeta.Pacer
	setAttack(attack AttackDescriptor)
	parsePacingCSV(csv *csv.Reader) []RateDescriptor
	parsePacingStr(pacing string) []RateDescriptor
}

// pacerState maintains state for a Vegeta pacer.
type pacerState struct {
	Rate    RateDescriptor
	Metrics vegeta.Metrics
}

// activePacerState maintains state for the actively executing Vegeta pacer
// This is only necessary because the `Pace` function is called by value
// rather than by reference.
var activePacerState pacerState

// Rounding support lifted from Vegeta reporters since it is private.
var durations = [...]time.Duration{
	time.Hour,
	time.Minute,
	time.Second,
	time.Millisecond,
	time.Microsecond,
	time.Nanosecond,
}

// round to the next most precise unit.
func round(d time.Duration) time.Duration {
	for i, unit := range durations {
		if d >= unit && i < len(durations)-1 {
			return d.Round(durations[i+1])
		}
	}
	return d
}

/**
Dynamic Pacer Implementations
**/

// StepFunctionPacer paces an attack with specific request rates for specific durations.
type StepFunctionPacer struct {
	Attack AttackDescriptor
}

func (pacer StepFunctionPacer) String() string {
	return fmt.Sprintf("StepFunctionPacer Rates{%s: %d rates}", pacer.Attack.Name, len(pacer.Attack.Rates))
}

// Pace determines the length of time to sleep until the next hit is sent.
func (pacer StepFunctionPacer) Pace(elapsed time.Duration, hits uint64) (time.Duration, bool) {
	// Determine which Rate is active and accumulate and expected number of hits
	var activeRate RateDescriptor
	aggregateDuration := time.Second
	expectedHits := pacer.hits(elapsed)
	for _, rate := range pacer.Attack.Rates {
		aggregateDuration += rate.Duration
		if elapsed <= aggregateDuration {
			activeRate = rate
			break
		}
	}

	// Use the last rate if we didn't find one
	if activeRate == (RateDescriptor{}) {
		activeRate = pacer.Attack.Rates[len(pacer.Attack.Rates)-1]
		fmt.Printf("🔥  Setting default rate of %d req/sec for remainder of attack\n", activeRate.Rate)
	}

	// Report when the rate changes
	if activePacerState.Rate != activeRate {
		if activePacerState.Metrics.Requests > 0 {
			activePacerState.Metrics.Close()
			reporter := vegeta.NewTextReporter(&activePacerState.Metrics)
			reporter.Report(os.Stdout)
			activePacerState.Metrics = vegeta.Metrics{}
		}

		activePacerState.Rate = activeRate
		elapsedSummary := func() string {
			if uint64(elapsed.Seconds()) > 0 {
				return fmt.Sprintf(" (%v elapsed)", round(elapsed))
			}
			return ""
		}
		fmt.Printf("💥  Attacking at a rate of %v%s\n", activeRate, elapsedSummary())
	}

	// Calculate when to send the next hit based on the active rate
	if hits == 0 || hits < uint64(expectedHits) {
		// Running behind, send next hit immediately.
		return 0, false
	}

	nsPerHit := math.Round(1 / pacer.hitsPerNs(activeRate))
	hitsToWait := float64(hits+1) - float64(expectedHits)
	nextHitIn := time.Duration(nsPerHit * hitsToWait)
	return nextHitIn, false
}

// TODO: Move to RateDescriptor type
func (pacer StepFunctionPacer) hitsPerNs(rate RateDescriptor) float64 {
	return float64(rate.Rate) / float64(time.Second)
}

// TODO: Move to RateDescriptor type
func (pacer StepFunctionPacer) hits(duration time.Duration) float64 {
	hits := float64(0)
	aggregateDuration := time.Second
	for _, rate := range pacer.Attack.Rates {
		hits += (float64(rate.Rate) * rate.Duration.Seconds())
		aggregateDuration += rate.Duration
		if duration <= aggregateDuration {
			break
		}
	}
	return hits
}

func (pacer StepFunctionPacer) parsePacingCSV(csv *csv.Reader) []RateDescriptor {
	var rates []RateDescriptor
	for {
		line, err := csv.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		rate, err := strconv.Atoi(line[0])
		if err != nil {
			log.Fatal(err)
		}
		duration, err := time.ParseDuration(strings.TrimSpace(line[1]))
		if err != nil {
			log.Fatal(err)
		}
		rates = append(rates, RateDescriptor{
			Rate:     uint(rate),
			Duration: duration,
		})
	}
	return rates
}

// parsePacingStr parses a string of the form "duration1@rate1, duration2@rate2"... into an array of rate descriptors
func (pacer StepFunctionPacer) parsePacingStr(pacing string) []RateDescriptor {
	var rates []RateDescriptor
	descriptors := strings.SplitN(pacing, ",", -1)
	for _, descriptor := range descriptors {
		components := strings.SplitN(descriptor, "@", 2)
		if components[0] == "" || components[1] == "" {
			msg := fmt.Errorf("invalid pacing descriptor %q", pacing)
			log.Fatal(msg)
		}

		duration, err := time.ParseDuration(strings.TrimSpace(components[0]))
		if err != nil {
			msg := fmt.Errorf("invalid pacing descriptor %q: %s", pacing, err)
			log.Fatal(msg)
		}
		rate, err := strconv.Atoi(strings.TrimSpace(components[1]))
		if err != nil {
			msg := fmt.Errorf("invalid pacing descriptor %q: %s", pacing, err)
			log.Fatal(msg)
		}
		rates = append(rates, RateDescriptor{
			Rate:     uint(rate),
			Duration: duration,
		})
	}
	return rates
}

func (pacer *StepFunctionPacer) setAttack(attack AttackDescriptor) {
	pacer.Attack = attack
}

//------------------------------------------------------------------

// CurveFittingPacer chases a set of target rates spread out of a total duration.
type CurveFittingPacer struct {
	Duration time.Duration
	Attack   AttackDescriptor
	Slope    float64
}

func (pacer CurveFittingPacer) String() string {
	return fmt.Sprintf("CurveFittingPacer Rates{%s: %d rates}", pacer.Attack.Name, len(pacer.Attack.Rates))
}

// Pace determines the length of time to sleep until the next hit is sent.
func (pacer CurveFittingPacer) Pace(elapsed time.Duration, hits uint64) (time.Duration, bool) {
	// Determine which Rate is active and accumulate and expected number of hits
	var activeRate RateDescriptor
	aggregateDuration := time.Second
	for _, rate := range pacer.Attack.Rates {
		aggregateDuration += rate.Duration
		if elapsed <= aggregateDuration {
			activeRate = rate
			break
		}
	}

	// Use the last rate if we didn't find one
	if activeRate == (RateDescriptor{}) {
		activeRate = pacer.Attack.Rates[len(pacer.Attack.Rates)-1]
		fmt.Printf("🔥  Setting default rate of %d req/sec for remainder of attack\n", activeRate.Rate)
	}

	// Report when the rate changes
	if activePacerState.Rate != activeRate {
		if activePacerState.Metrics.Requests > 0 {
			activePacerState.Metrics.Close()
			reporter := vegeta.NewTextReporter(&activePacerState.Metrics)
			reporter.Report(os.Stdout)
			activePacerState.Metrics = vegeta.Metrics{}
		}

		activePacerState.Rate = activeRate
		elapsedSummary := func() string {
			if uint64(elapsed.Seconds()) > 0 {
				return fmt.Sprintf(" (%v elapsed)", round(elapsed))
			}
			return ""
		}
		fmt.Printf("💥  Attacking at a rate of %v%s\n", activeRate, elapsedSummary())
	}
	expectedHits := pacer.hits(elapsed, activeRate)

	// Calculate when to send the next hit based on the active rate
	if hits == 0 || hits < uint64(expectedHits) {
		// Running behind, send next hit immediately.
		return 0, false
	}

	// rate := p.rate(elapsed)
	interval := math.Round(1e9 / float64(pacer.rate(elapsed, activeRate)))

	if n := uint64(interval); n != 0 && math.MaxInt64/n < hits {
		// We would overflow wait if we continued, so stop the attack.
		return 0, true
	}

	delta := float64(hits+1) - expectedHits
	wait := time.Duration(interval * delta)

	return wait, false

	// nsPerHit := math.Round(1 / pacer.hitsPerNs(activeRate))
	// hitsToWait := float64(hits+1) - float64(expectedHits)
	// nextHitIn := time.Duration(nsPerHit * hitsToWait)
	// return nextHitIn, false

	// TODO:
	// Find your current rate
	// Find that rate you might be chasing
	// Check if duration is >= your aggregate duration, then start chasing the next one
}

func (rs RateDescriptor) hitsPerNs() float64 {
	return float64(rs.Rate) / float64(rs.Duration)
}

// hits returns the number of hits that have been sent during an attack
// lasting t nanoseconds. It returns a float so we can tell exactly how
// much we've missed our target by when solving numerically in Pace.
func (pacer CurveFittingPacer) hits(t time.Duration, rate RateDescriptor) float64 {
	if t < 0 {
		return 0
	}

	// TODO: Iterate across the rates and accumulate
	a := pacer.Slope
	// TODO: Average out the number hits across the windows?
	b := rate.hitsPerNs() * 1e9
	x := t.Seconds()

	return (a*math.Pow(x, 2))/2 + b*x
}

// rate calculates the instantaneous rate of attack at
// t nanoseconds after the attack began.
func (pacer CurveFittingPacer) rate(t time.Duration, rate RateDescriptor) float64 {
	a := pacer.Slope
	x := t.Seconds()
	b := rate.hitsPerNs() * 1e9
	return a*x + b
}

// // TODO: Move to RateDescriptor type
// func (pacer CurveFittingPacer) hitsPerNs(rate RateDescriptor) float64 {
// 	return float64(rate.Rate) / float64(time.Second)
// }

// // TODO: Move to RateDescriptor type
// func (pacer CurveFittingPacer) hits(duration time.Duration) float64 {
// 	hits := float64(0)
// 	aggregateDuration := time.Second
// 	for _, rate := range pacer.Attack.Rates {
// 		hits += (float64(rate.Rate) * rate.Duration.Seconds())
// 		aggregateDuration += rate.Duration
// 		if duration <= aggregateDuration {
// 			break
// 		}
// 	}
// 	return hits
// }

func (pacer CurveFittingPacer) durationForRatesLen(l int) time.Duration {
	return time.Second * time.Duration(math.Ceil(float64(pacer.Duration)/float64(l)))
}

func (pacer CurveFittingPacer) parsePacingCSV(csv *csv.Reader) []RateDescriptor {
	var rates []RateDescriptor
	for {
		line, err := csv.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		rate, err := strconv.Atoi(line[0])
		if err != nil {
			log.Fatal(err)
		}
		rates = append(rates, RateDescriptor{
			Rate: uint(rate),
		})
	}

	// Split the total duration amongst our rate descriptors
	duration := pacer.durationForRatesLen(len(rates))
	for _, rate := range rates {
		rate.Duration = duration
	}

	return rates
}

// parsePacingStr parses a string of the form "rate1, rate2"... into an array of rate descriptors
func (pacer CurveFittingPacer) parsePacingStr(pacing string) []RateDescriptor {
	var rates []RateDescriptor
	descriptors := strings.SplitN(pacing, ",", -1)
	duration := pacer.durationForRatesLen(len(descriptors))
	for _, descriptor := range descriptors {
		rate, err := strconv.Atoi(strings.TrimSpace(descriptor))
		if err != nil {
			msg := fmt.Errorf("invalid pacing descriptor %q: %s", pacing, err)
			log.Fatal(msg)
		}
		rates = append(rates, RateDescriptor{
			Rate:     uint(rate),
			Duration: duration,
		})
	}
	return rates
}

func (pacer *CurveFittingPacer) setAttack(attack AttackDescriptor) {
	pacer.Attack = attack
}

/**
CLI interface
**/

// paceOpts aggregates the pacing command line options
type paceOpts struct {
	url      string
	file     string
	pacer    string
	pacing   string
	duration time.Duration
}

func main() {
	const StepFunctionArg = "step-function"
	const CurveFittingArg = "curve-fitting"

	var PacerArgs = []string{StepFunctionArg, CurveFittingArg}

	// Parse the commandline options
	opts := paceOpts{}
	flag.StringVar(&opts.url, "url", "http://localhost:8080/", "The URL to attack")
	flag.StringVar(&opts.pacer, "pacer", "",
		fmt.Sprintf("Pacer to use for governing load rate [%s]", strings.Join(PacerArgs, ", ")))
	flag.StringVar(&opts.pacing, "pacing", "", "String describing the pace")
	flag.StringVar(&opts.file, "file", "", "CSV file describing the pace")
	flag.DurationVar(&opts.duration, "duration", 0, fmt.Sprintf("Duration of the test. Required when pacer is %q", CurveFittingArg))
	flag.Parse()

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	_, err := url.ParseRequestURI(opts.url)
	if err != nil {
		msg := fmt.Errorf("invalid URL %q: %s", opts.url, err)
		log.Fatal(msg)
	}

	if opts.file == "" && opts.pacing == "" {
		err := fmt.Errorf("-file or -pacing must be provided")
		log.Fatal(err)
	} else if opts.file != "" && opts.pacing != "" {
		err := fmt.Errorf("-file and -pacing cannot both be provided")
		log.Fatal(err)
	}

	var pacer dynamicPacer
	switch opts.pacer {
	case StepFunctionArg:
		pacer = &StepFunctionPacer{}
	case CurveFittingArg:
		if opts.duration <= 0 {
			err := fmt.Errorf("%q pacer requires a -duration be provided", opts.pacer)
			log.Fatal(err)
		}
		fmt.Printf("Configuring with duration: %v", opts.duration)
		pacer = &CurveFittingPacer{Duration: opts.duration, Slope: 1}
	default:
		err := fmt.Errorf("unknown pacer type: %q", opts.pacer)
		log.Fatal(err)
	}

	// Build the attack
	var attack AttackDescriptor
	attack.Name = "Variable Load Test"
	if opts.file != "" {
		csvFile, err := os.Open(opts.file)
		if err != nil {
			log.Fatal(err)
		}
		csv := csv.NewReader(bufio.NewReader(csvFile))
		attack.Rates = pacer.parsePacingCSV(csv)
	} else if opts.pacing != "" {
		attack.Rates = pacer.parsePacingStr(opts.pacing)
	}
	pacer.setAttack(attack)

	// Run the attack
	fmt.Printf("🚀  Starting variable load test against %q with %d load profiles for %v\n", opts.url, len(attack.Rates), round(attack.Duration()))
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    opts.url,
	})
	attacker := vegeta.NewAttacker()
	startedAt := time.Now()
	for res := range attacker.Attack(targeter, pacer, attack.Duration(), attack.Name) {
		activePacerState.Metrics.Add(res)
	}

	activePacerState.Metrics.Close()

	reporter := vegeta.NewTextReporter(&activePacerState.Metrics)
	reporter.Report(os.Stdout)

	attackDuration := time.Since(startedAt)
	fmt.Printf("✨  Variable load test against %q completed in %v\n", opts.url, round(attackDuration))
}
