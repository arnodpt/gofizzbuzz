# gofizzbuzz

A fizzbuzz REST server.

## Summary

The original fizz-buzz consists in writing all numbers from 1 to 100, and just replacing all multiples of 3 by "fizz", all multiples of 5 by "buzz", and all multiples of 15 by "fizzbuzz". 

The output would look like this: "1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16,...".

This server exposes two REST API endpoints :

**/fizzbuzz** :
  - Accepts five parameters: three integers int1, int2 and limit, and two strings str1 and str2.
  - Returns a list of strings with numbers from 1 to limit, where: all multiples of int1 are replaced by str1, all multiples of int2 are replaced by str2, all multiples of int1 and int2 are replaced by str1str2.


**/popular** :
  - Accept no parameter.
  - Return the parameters corresponding to the most used request, as well as the number of hits for this request.

## Prerequisites

To run this server you should install the following prerequisites :

- [make](https://linux.die.net/man/1/make) *v4.2.1*
- [docker](https://docs.docker.com/get-docker/) *v20.10.12*
- [docker-compose](https://docs.docker.com/compose/install/) *1.29.2*

To develop, build and test this server you should install Go :

- [go](https://go.dev/doc/install) *v1.16.7*

The versions described here are the one I personally used for this project.

## Setup

First clone the project :

```bash
  git clone https://github.com/arnodpt/gofizzbuzz
  cd gofizzbuzz
```

Modify the [.env.sh](./.env.sh) environment file to set the config (ports etc) or use the default one :

## Run

To run the server simply use :

```bash
make
```

To stop the server simply use :

```bash
make stop
```

## Usage

To use this server you can use GET requests on those two endpoints :

**/fizzbuzz** :
  - Accepts five parameters: three integers int1, int2 and limit, and two strings str1 and str2.
  - Returns a list of strings with numbers from 1 to limit, where: all multiples of int1 are replaced by str1, all multiples of int2 are replaced by str2, all multiples of int1 and int2 are replaced by str1str2.
  - You can use the [curl_fizzbuzz.sh](./curl_fizzbuzz.sh) curl script to test the endpoint

**/popular** :
  - Accept no parameter.
  - Return the parameters corresponding to the most used request, as well as the number of hits for this request.
  - You can use the [curl_popular.sh](./curl_popular.sh) curl script to test the endpoint

For more informations about the REST API you can check the documentation generated with go-swagger in [doc/markdown.md](doc/markdown.md).

## Tests

Test suite is made using built-in golang test command `go test`. You should have installed golang to be able to run them.

To run them simply use :

```bash
make test
```

## Monitoring

This project launches a Prometheus and a Grafana to monitor the server, in particular, the metrics about the number of requests by query parameters for the `/fizzbuzz` endpoint, that are used by the `/popular` endpoint.

Prometheus is accessible at `http://localhost:9090` (or `PROMETHEUS_PORT` in config).

Grafana is accessible at `http://localhost:3000` (or `GRAFANA_PORT` in config). User is `admin` and the credential is available in [docker-compose.yml](./docker-compose.yml).

## How it was made

This server has 3 containers, all described in the [docker-compose.yml](./docker-compose.yml) :

- A golang server exposing the REST API
- A prometheus that scrapes the metrics from the server
- A grafana that displays prometheus metrics

For the server, I used the golang library [gofiber](https://github.com/gofiber/fiber).

The code for the two endpoints is available in the [api/](./api/) folder.

For the `/popular` endpoint, I chose to only use prometheus metrics and not use a separate DB since this kind of metrics is clearly made for monitoring purpose, to know what parameters the users use the most.\
To push a metric at each api usage, I implemented a middleware using the [golang prometheus instrumentation library](https://github.com/prometheus/client_golang) available in the [middleware/](./middleware/) folder.\
Since the metrics were not persistent, I had to implement my own persistence using a goroutine that gathers metrics from prometheus and writes them in the `server_data` folder.
