# Vegeta Variable Load Testing

This repository contains an example of running a multi-rate load test using [Vegeta](https://github.com/tsenart/vegeta)
as a library. This implements a variable load testing utility without any forking or modification of existing project code bases.

The utility will run an attack on a targetr URL (default of http://localhost:8080/) using the rate patterns defined in the `attack.csv` file.

This code was adapted from a 2017 article on Josh Barrats blog: https://serialized.net/2017/06/load-testing-with-vegeta-and-python/

## Example Output

```
$ go run vegeta_varload.go https://www.opsani.com/
🚀  Start variable load test against https://www.opsani.com/ with 6 load profiles for 44 total seconds
💥  Attacking at rate of 100 req/sec (0 seconds elapsed)
💥  Attacking at rate of 200 req/sec (6 seconds elapsed)
💥  Attacking at rate of 300 req/sec (11 seconds elapsed)
💥  Attacking at rate of 350 req/sec (21 seconds elapsed)
💥  Attacking at rate of 200 req/sec (31 seconds elapsed)
✨  Attack completed (latency 30.001291979s, 10851 requests sent)
Requests      [total, rate, throughput]         10851, 199.99, 19.60
Duration      [total, attack, wait]             1m1s, 54.256s, 7.136s
Latencies     [min, mean, 50, 90, 95, 99, max]  77.84ms, 21.728s, 30s, 30.001s, 30.001s, 30.002s, 30.025s
Bytes In      [total, mean]                     164683161, 15176.77
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           11.09%
Status Codes  [code:count]                      0:8032  200:1203  403:1616
Error Set:
403 Forbidden
Get https://www.opsani.com/: EOF
Get https://www.opsani.com/: net/http: request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)
net/http: request canceled (Client.Timeout exceeded while reading body)
Get https://www.opsani.com/: net/http: request canceled (Client.Timeout exceeded while awaiting headers)
```

## Building & Running on macOS

All the dependencies necessary to build are available via Homebrew.:

1. Install golang: `$ brew install golang`
2. Install Vegeta dependency: `$ go get -u github.com/tsenart/vegeta`
3. Build and run with Go: `$ go run vegeta_breaker.go`

## Running via Docker

A Dockerfile is provided that can be used to run the load test

```bash
$ docker build -t vegeta_varload .
$ docker run -ti vegeta_varload https://www.opsani.com/
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