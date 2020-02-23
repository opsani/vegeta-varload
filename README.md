# Vegeta Variable Load Testing

This repository contains an example of running a multi-rate load test using [Vegeta](https://github.com/tsenart/vegeta)
as a library. This implements a variable load testing utility without any forking or modification of existing project code bases.

The utility will run an attack on a targetr URL (default of http://localhost:8080/) using the rate patterns defined in the `attack.csv` file.

This code was adapted from a 2017 article on Josh Barrats blog: https://serialized.net/2017/06/load-testing-with-vegeta-and-python/

## Example Output

```
$ go run vegeta_varload.go https://www.opsani.com/
🚀  Start variable load test against https://www.opsani.com/ with 6 load profiles for 44 total seconds
💥  Attacking at rate of 10 req/sec for 5s (0s elapsed)
💥  Attacking at rate of 20 req/sec for 5s (6s elapsed)
💥  Attacking at rate of 30 req/sec for 10s (11s elapsed)
💥  Attacking at rate of 40 req/sec for 10s (21s elapsed)
💥  Attacking at rate of 50 req/sec for 5s (31s elapsed)
💥  Attacking at rate of 60 req/sec for 8s (36s elapsed)
✨  Attack completed (latency 764.511853ms, 2640 requests sent)
Requests      [total, rate, throughput]         2640, 60.00, 32.98
Duration      [total, attack, wait]             44.083s, 44.001s, 81.816ms
Latencies     [min, mean, 50, 90, 95, 99, max]  77.216ms, 272.827ms, 160.355ms, 679.439ms, 764.512ms, 866.838ms, 1.042s
Bytes In      [total, mean]                     154226160, 58419.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           55.08%
Status Codes  [code:count]                      200:1454  403:1186
Error Set:
403 Forbidden
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