# Vegeta Varload

This repository contains a variable rate load testing utility built on top of [Vegeta](https://github.com/tsenart/vegeta). It is capable of delivering a dynamic load to an application under test by providing a set of implementations of the [`vegeta.Pacer`](https://github.com/tsenart/vegeta/blob/master/lib/pacer.go) interface. Pacer implementations govern the load rate by determining when the next request should be sent -- allowing for the acceleration or deceleration of the request rate based on arbitrary logic.

### Pacers

There are currently two pacers provided:

* `step-function`: Models a load in which the load oscillates between target rates at specific durations. For example, you may deliver a load of 50 req/s for 10 seconds, followed by 150 req/s for 25 seconds, and so on. The acceleration and deceleration curve between rates is as aggressive as possible and the total test duration is defined by the sum of all step durations.
* `curve-fitting`: Models a load in which the load is described by a set of rates and a total duration of the test. The pacer will accelerate and decelerate between the data points to construct a load curve that is the best fit for the data points.

## Usage

```console
$ ./vegeta-varload --help
Usage of vegeta-varload:
  -duration duration
    	Duration of the test. Required when pacer is "curve-fitting"
  -file string
    	CSV file describing the pace
  -pacer string
    	Pacer to use for governing load rate [step-function, curve-fitting]
  -pacing string
    	String describing the pace
  -url string
    	The URL to attack (default "http://localhost:8080/")
```

## Example Output

```console
$ go run vegeta-varload.go --url https://www.opsani.com/ --pacer step-function --file attack.csv
🚀  Starting variable load test against https://www.opsani.com/ with 6 load profiles for 44 total seconds
💥  Attacking at rate of 10 req/sec for 5s
Requests      [total, rate, throughput]         59, 9.99, 9.87
Duration      [total, attack, wait]             5.977s, 5.904s, 73.411ms
Latencies     [min, mean, 50, 90, 95, 99, max]  68.654ms, 214.712ms, 238.093ms, 253.151ms, 253.714ms, 254.256ms, 254.276ms
Bytes In      [total, mean]                     653189, 11071.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:59
Error Set:
💥  Attacking at rate of 20 req/sec for 5s (6.001s elapsed)
Requests      [total, rate, throughput]         159, 32.43, 31.98
Duration      [total, attack, wait]             4.972s, 4.903s, 69.018ms
Latencies     [min, mean, 50, 90, 95, 99, max]  66.883ms, 108.037ms, 82.134ms, 202.123ms, 202.533ms, 202.859ms, 202.892ms
Bytes In      [total, mean]                     1760289, 11071.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:159
Error Set:
💥  Attacking at rate of 30 req/sec for 10s (11.003s elapsed)
Requests      [total, rate, throughput]         409, 41.10, 40.81
Duration      [total, attack, wait]             10.022s, 9.952s, 70.199ms
Latencies     [min, mean, 50, 90, 95, 99, max]  68.609ms, 158.846ms, 83.582ms, 344.221ms, 368.986ms, 387.947ms, 388.953ms
Bytes In      [total, mean]                     4528039, 11071.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:409
Error Set:
💥  Attacking at rate of 40 req/sec for 10s (21.005s elapsed)
Requests      [total, rate, throughput]         609, 61.10, 60.69
Duration      [total, attack, wait]             10.035s, 9.967s, 68.171ms
Latencies     [min, mean, 50, 90, 95, 99, max]  65.577ms, 87.344ms, 70.846ms, 136.676ms, 144.97ms, 148.751ms, 150.887ms
Bytes In      [total, mean]                     6742239, 11071.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:609
Error Set:
💥  Attacking at rate of 50 req/sec for 5s (31.003s elapsed)
Requests      [total, rate, throughput]         560, 112.16, 110.66
Duration      [total, attack, wait]             5.061s, 4.993s, 67.685ms
Latencies     [min, mean, 50, 90, 95, 99, max]  67.504ms, 132.073ms, 101.869ms, 247.682ms, 266.71ms, 297.917ms, 301.779ms
Bytes In      [total, mean]                     6199760, 11071.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:560
Error Set:
💥  Attacking at rate of 60 req/sec for 8s (36s elapsed)
Requests      [total, rate, throughput]         844, 104.68, 103.79
Duration      [total, attack, wait]             8.132s, 8.063s, 69.164ms
Latencies     [min, mean, 50, 90, 95, 99, max]  67.733ms, 110.846ms, 74.004ms, 193.512ms, 216.077ms, 293.628ms, 305.455ms
Bytes In      [total, mean]                     9343924, 11071.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:844
Error Set:
✨  Variable load test against https://www.opsani.com/ completed in 44.073s
```

## Building & Running on macOS

All the dependencies necessary to build are available via Homebrew.:

1. Install golang: `$ brew install golang`
2. Install Vegeta dependency: `$ go get -u github.com/tsenart/vegeta`
3. Build and run with Go: `$ go run vegeta-varload.go`

## Running via Docker

A Dockerfile is provided that can be used to run the load test

```console
$ docker build -t vegeta-varload .
$ docker run -ti vegeta-varload --url https://www.opsani.com/ --pacer step-function --pacing "10s@15, 1m@20"
```

## Running via Docker Compose

A Docker Compose assembly is provided that will run Nginx in one container and load test it with Vegeta Varload in another.

```console
$ docker-compose up -d
$ docker-compose logs -f vegeta
```

### Running ad-hoc load tests through Compose

```console
$ docker-compose run vegeta https://www.opsani.com/
```

## Acknowledgements

This code was originally adapted from a 2017 article on Josh Barrats blog: https://serialized.net/2017/06/load-testing-with-vegeta-and-python/

## License

This code is licensed under the terms of the MIT Open Source license just as Vegeta is.
