# Vegeta Variable Load Testing

This repository contains an example of running a multi-rate load test using [Vegeta](https://github.com/tsenart/vegeta)
as a library. This implements a variable load testing utility without any forking or modification of existing project code bases.

The utility will run an attack on a target URL (default of http://localhost:8080/) using the rate patterns defined in the `attack.csv` file.

This code was adapted from a 2017 article on Josh Barrats blog: https://serialized.net/2017/06/load-testing-with-vegeta-and-python/

## Example Output

```
$ go run vegeta-varload.go https://www.opsani.com/
🚀  Start variable load test against https://golang.org/ with 6 load profiles for 44 total seconds
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
✨  Attack completed in 44.073s
```

## Building & Running on macOS

All the dependencies necessary to build are available via Homebrew.:

1. Install golang: `$ brew install golang`
2. Install Vegeta dependency: `$ go get -u github.com/tsenart/vegeta`
3. Build and run with Go: `$ go run vegeta-varload.go`

## Running via Docker

A Dockerfile is provided that can be used to run the load test

```bash
$ docker build -t vegeta-varload .
$ docker run -ti vegeta-varload https://www.opsani.com/
```

## Running via Docker Compose

A Docker Compose assembly is provided that will run Nginx in one container and load test it with Vegeta Varload in another.

```bash
$ docker-compose up -d
$ docker-compose logs -f vegeta
```

### Running ad-hoc load tests through Compose

```bash
$ docker-compose run vegeta https://www.opsani.com/
```

## License

This code is licensed under the terms of the MIT Open Source license just as Vegeta is.
