# Vegeta Variable Load Testing

This repository contains an example of running a multi-rate load test using [Vegeta](https://github.com/tsenart/vegeta)
as a library. This implements a variable load testing utility without any forking or modification of existing project code bases.

The utility will run an attack on a target URL (default of http://localhost:8080/) using the rate patterns defined in the `attack.csv` file.

This code was adapted from a 2017 article on Josh Barrats blog: https://serialized.net/2017/06/load-testing-with-vegeta-and-python/

## Example Output

```
$ go run vegeta-varload.go https://www.opsani.com/
🚀  Start variable load test against https://www.opsani.com/ with 6 load profiles for 44 total seconds
💥  Attacking at rate of 10 req/sec for 5s (0s elapsed)
Requests      [total, rate, throughput]         59, 9.99, 9.88
Duration      [total, attack, wait]             5.974s, 5.904s, 69.698ms
Latencies     [min, mean, 50, 90, 95, 99, max]  68.801ms, 234.386ms, 260.643ms, 278.544ms, 282.593ms, 386.85ms, 397.14ms
Bytes In      [total, mean]                     653189, 11071.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:59
Error Set:
💥  Attacking at rate of 20 req/sec for 5s (6s elapsed)
Requests      [total, rate, throughput]         159, 32.46, 32.00
Duration      [total, attack, wait]             4.969s, 4.898s, 71.609ms
Latencies     [min, mean, 50, 90, 95, 99, max]  66.987ms, 108.663ms, 92.837ms, 186.03ms, 189.781ms, 190.907ms, 191.036ms
Bytes In      [total, mean]                     1760289, 11071.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:159
Error Set:
💥  Attacking at rate of 30 req/sec for 10s (11s elapsed)
Requests      [total, rate, throughput]         409, 41.11, 40.81
Duration      [total, attack, wait]             10.021s, 9.949s, 71.706ms
Latencies     [min, mean, 50, 90, 95, 99, max]  67.787ms, 178.439ms, 132.946ms, 366.304ms, 391.335ms, 393.376ms, 394.914ms
Bytes In      [total, mean]                     4528039, 11071.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:409
Error Set:
💥  Attacking at rate of 40 req/sec for 10s (21s elapsed)
Requests      [total, rate, throughput]         610, 61.06, 60.65
Duration      [total, attack, wait]             10.058s, 9.99s, 68.371ms
Latencies     [min, mean, 50, 90, 95, 99, max]  66.283ms, 94.093ms, 72.092ms, 167.423ms, 174.303ms, 178.791ms, 179.701ms
Bytes In      [total, mean]                     6753310, 11071.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:610
Error Set:
💥  Attacking at rate of 50 req/sec for 5s (31s elapsed)
Requests      [total, rate, throughput]         559, 112.47, 110.97
Duration      [total, attack, wait]             5.037s, 4.97s, 67.161ms
Latencies     [min, mean, 50, 90, 95, 99, max]  66.221ms, 119.805ms, 113.698ms, 186.62ms, 233.64ms, 257.057ms, 332.034ms
Bytes In      [total, mean]                     6188689, 11071.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:559
Error Set:
💥  Attacking at rate of 60 req/sec for 8s (36s elapsed)
Requests      [total, rate, throughput]         844, 104.69, 103.80
Duration      [total, attack, wait]             8.131s, 8.062s, 68.763ms
Latencies     [min, mean, 50, 90, 95, 99, max]  64.79ms, 119.58ms, 75.722ms, 223.715ms, 233.403ms, 238.961ms, 241.214ms
Bytes In      [total, mean]                     9343924, 11071.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:844
Error Set:
✨  Attack completed in 44.071653509s
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
