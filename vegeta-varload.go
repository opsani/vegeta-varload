package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

// RateDescriptor describes a rate in requests per second and duration.
type RateDescriptor struct {
	Rate     uint          `json:"rate"`
	Duration time.Duration `json:"duration"`
}

func (rs RateDescriptor) String() string {
	return fmt.Sprintf("%vreq/s for %v", rs.Rate, rs.Duration)
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

// StepFunctionPacer paces an attack with specific request rates for specific durations.
type StepFunctionPacer struct {
	Attack AttackDescriptor
}

func (vrp StepFunctionPacer) String() string {
	return fmt.Sprintf("Variable Rates{%s: %d rates}", vrp.Attack.Name, len(vrp.Attack.Rates))
}

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

// pacerState maintains state for a Vegeta pacer.
type pacerState struct {
	Rate    RateDescriptor
	Metrics vegeta.Metrics
}

// activePacerState maintains state for the actively executing Vegeta pacer
// This is only necessary because the `Pace` function is called by value
// rather than by reference.
var activePacerState pacerState

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
		fmt.Printf("🔥  Setting default rate of %dreq/sec for remainder of attack\n", activeRate.Rate)
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
	if hits < uint64(expectedHits) {
		// Running behind, send next hit immediately.
		return 0, false
	}

	nsPerHit := math.Round(1 / pacer.hitsPerNs(activeRate))
	hitsToWait := float64(hits+1) - float64(expectedHits)
	nextHitIn := time.Duration(nsPerHit * hitsToWait)
	return nextHitIn, false
}

func (pacer StepFunctionPacer) hitsPerNs(rate RateDescriptor) float64 {
	return float64(rate.Rate) / float64(time.Second)
}

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

func main() {
	// Load the attack CSV
	csvFile, _ := os.Open("attack.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var attack AttackDescriptor
	attack.Name = "Variable Load Test"
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		rate, _ := strconv.Atoi(line[0])
		duration, _ := time.ParseDuration(strings.TrimSpace(line[1]))
		attack.Rates = append(attack.Rates, RateDescriptor{
			Rate:     uint(rate),
			Duration: duration,
		})
	}

	// Run the attack
	targetURL := "http://localhost:8080/"
	if len(os.Args) == 2 {
		targetURL = os.Args[1]
	}
	fmt.Printf("🚀  Starting variable load test against %s with %d load profiles for %v\n", targetURL, len(attack.Rates), round(attack.Duration()))
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    targetURL,
	})
	pacer := StepFunctionPacer{Attack: attack}
	attacker := vegeta.NewAttacker()
	startedAt := time.Now()
	for res := range attacker.Attack(targeter, pacer, attack.Duration(), attack.Name) {
		activePacerState.Metrics.Add(res)
	}

	activePacerState.Metrics.Close()

	reporter := vegeta.NewTextReporter(&activePacerState.Metrics)
	reporter.Report(os.Stdout)

	attackDuration := time.Since(startedAt)
	fmt.Printf("✨  Variable load test against %s completed in %v\n", targetURL, round(attackDuration))
}
