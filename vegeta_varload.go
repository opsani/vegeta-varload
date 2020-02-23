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

type RateDescriptor struct {
	Rate     uint          `json:"rate"`
	Duration time.Duration `json:"duration"`
}

type AttackDescriptor struct {
	Name  string           `json:"name"`
	Rates []RateDescriptor `json:"rates"`
}

func (attack AttackDescriptor) Duration() time.Duration {
	duration := time.Second
	for _, rate := range attack.Rates {
		duration += rate.Duration
	}
	return duration
}

type MultiRatePacer struct {
	Attack AttackDescriptor
}

func (mrp MultiRatePacer) String() string {
	return fmt.Sprintf("Multi Rate{%s: %d rates}", mrp.Attack.Name, len(mrp.Attack.Rates))
}

var CurrentRate RateDescriptor

// Pace determines the length of time to sleep until the next hit is sent.
func (mrp MultiRatePacer) Pace(elapsed time.Duration, hits uint64) (time.Duration, bool) {
	// Determine which Rate is active and accumulate and expected number of hits
	var activeRate RateDescriptor
	expectedHits := uint64(0)
	aggregateDuration := time.Second
	for _, rate := range mrp.Attack.Rates {
		expectedHits += (uint64(rate.Rate) * uint64(rate.Duration/time.Second))
		aggregateDuration += rate.Duration
		if elapsed <= aggregateDuration {
			activeRate = rate
			break
		}
	}

	// Use the last rate if we didn't find one
	if activeRate == (RateDescriptor{}) {
		activeRate = mrp.Attack.Rates[len(mrp.Attack.Rates)-1]
		fmt.Printf("🔥  Setting default rate of %d req/sec for remainder of attack\n", activeRate.Rate)
	}

	// Report when the rate changes
	if CurrentRate != activeRate {
		CurrentRate = activeRate
		fmt.Printf("💥  Attacking at rate of %d req/sec for %v (%ds elapsed)\n", activeRate.Rate, activeRate.Duration, uint64(elapsed.Seconds()))
	}

	// Calculate when to send the next hit based on the active rate
	if hits < expectedHits {
		// Running behind, send next hit immediately.
		return 0, false
	}
	per := 1 * time.Second
	interval := uint64(per.Nanoseconds() / int64(activeRate.Rate))
	if math.MaxInt64/interval < hits {
		// We would overflow delta if we continued, so stop the attack.
		return 0, true
	}
	delta := time.Duration((hits + 1) * interval)

	// Zero or negative durations cause time.Sleep to return immediately.
	return delta - elapsed, false
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
	fmt.Printf("🚀  Start variable load test against %s with %d load profiles for %d total seconds\n", targetURL, len(attack.Rates), uint64(attack.Duration().Seconds()))
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    targetURL,
	})
	pacer := MultiRatePacer{Attack: attack}
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, pacer, attack.Duration(), attack.Name) {
		metrics.Add(res)
	}
	metrics.Close()

	latency := metrics.Latencies.P95
	fmt.Printf("✨  Attack completed (latency %s, %d requests sent)\n", latency, metrics.Requests)

	reporter := vegeta.NewTextReporter(&metrics)
	reporter.Report(os.Stdout)
}
