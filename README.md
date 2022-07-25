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

```
  git clone https://github.com/arnodpt/gofizzbuzz
  cd gofizzbuzz
```sh

Modify the [.env.sh](./.env.sh) environment file to set the config (ports etc) or use the default one :

## Run

To run the server simply use :

```
make
```sh

To stop the server simply use :

```
make stop
```sh

## 
