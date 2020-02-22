# Vegeta Variable Load Testing

This repository contains an example of running a multi-rate load test using [Vegeta](https://github.com/tsenart/vegeta)
as a library. This implements a variable load testing utility without any forking or modification of existing project code bases.

The utility will run an attack on http://localhost:8080/

This code was adapted from a 2017 article on Josh Barrats blog: https://serialized.net/2017/06/load-testing-with-vegeta-and-python/

## Example Output

```
✨  Success at 20 req/sec (latency 759.834µs)
✨  Success at 40 req/sec (latency 799.734µs)
✨  Success at 80 req/sec (latency 724.74µs)
✨  Success at 160 req/sec (latency 707.743µs)
✨  Success at 320 req/sec (latency 687.584µs)
✨  Success at 640 req/sec (latency 498.001µs)
✨  Success at 1280 req/sec (latency 240.247µs)
✨  Success at 2560 req/sec (latency 246.242µs)
✨  Success at 5120 req/sec (latency 795.901µs)
✨  Success at 10240 req/sec (latency 827.547µs)
✨  Success at 20480 req/sec (latency 3.539196ms)
✨  Success at 40960 req/sec (latency 328.227µs)
💥  Failed at 81920 req/sec (latency 1.253610495s)
✨  Success at 61440 req/sec (latency 15.069307ms)
✨  Success at 71680 req/sec (latency 12.952251ms)
✨  Success at 76800 req/sec (latency 14.35046ms)
✨  Success at 79360 req/sec (latency 7.224106ms)
✨  Success at 80640 req/sec (latency 30.891972ms)
✨  Success at 81280 req/sec (latency 6.531166ms)
✨  Success at 81600 req/sec (latency 50.381644ms)
✨  Success at 81760 req/sec (latency 143.061991ms)
✨  Success at 81840 req/sec (latency 9.716487ms)
```

## Building & Running on macOS

All the dependencies necessary to build are available via Homebrew.:

1. Install golang: `$ brew install golang`
2. Install Vegeta dependency: `$ go get -u github.com/tsenart/vegeta`
3. Build and run with Go: `$ go run vegeta_breaker.go`
