


# A fizzbuzz REST server
REST server that reproduces the original fizzbuzz algorithm but with parameters and metrics.
  

## Informations

### Version

1.0.0

## Content negotiation

### URI Schemes
  * http

### Consumes
  * application/json

### Produces
  * application/json

## All endpoints

###  operations

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| GET | /fizzbuzz | [get fizzbuzz](#get-fizzbuzz) | Returns a list of strings with numbers from 1 to limit, where all multiples of int1 are replaced by str1, all multiples of int2 are replaced by str2, all multiples of int1 and int2 are replaced by str1str2. |
| GET | /metrics | [get metrics](#get-metrics) | Returns prometheus metrics. |
| GET | /popular | [get popular](#get-popular) | Returns the parameters corresponding to the most used request, as well as the number of hits for this request. |
  


## Paths

### <span id="get-fizzbuzz"></span> Returns a list of strings with numbers from 1 to limit, where all multiples of int1 are replaced by str1, all multiples of int2 are replaced by str2, all multiples of int1 and int2 are replaced by str1str2. (*GetFizzbuzz*)

```
GET /fizzbuzz
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| int1 | `query` | int (formatted integer) | `int64` |  |  | `3` | The first int which multiples will be replaced by str1 or str1str2 if it's also a multiple of int2. |
| int2 | `query` | int (formatted integer) | `int64` |  |  | `5` | The second int which multiples will be replaced by str2 or str1str2 if it's also a multiple of int1. |
| limit | `query` | int (formatted integer) | `int64` |  |  | `50` | The number to which the fizzbuzz algorithm will end and the size of the final array. |
| str1 | `query` | string | `string` |  |  | `"fizz"` | The first string that will replace multiples of int1 and be part of multiples of int1 and int2. |
| str2 | `query` | string | `string` |  |  | `"buzz"` | The second string that will replace multiples of int2 and be part of multiples of int1 and int2. |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-fizzbuzz-200) | OK | An array of strings with numbers or str1, str2, str1str2 |  | [schema](#get-fizzbuzz-200-schema) |
| [400](#get-fizzbuzz-400) | Bad Request | Bad or missing query parameters or negative limit. |  | [schema](#get-fizzbuzz-400-schema) |
| [500](#get-fizzbuzz-500) | Internal Server Error | An internal error happened while creating the response. |  | [schema](#get-fizzbuzz-500-schema) |

#### Responses


##### <span id="get-fizzbuzz-200"></span> 200 - An array of strings with numbers or str1, str2, str1str2
Status: OK

###### <span id="get-fizzbuzz-200-schema"></span> Schema
   
  

[]string

##### <span id="get-fizzbuzz-400"></span> 400 - Bad or missing query parameters or negative limit.
Status: Bad Request

###### <span id="get-fizzbuzz-400-schema"></span> Schema

##### <span id="get-fizzbuzz-500"></span> 500 - An internal error happened while creating the response.
Status: Internal Server Error

###### <span id="get-fizzbuzz-500-schema"></span> Schema

### <span id="get-metrics"></span> Returns prometheus metrics. (*GetMetrics*)

```
GET /metrics
```

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-metrics-200) | OK | Prometheus metrics |  | [schema](#get-metrics-200-schema) |
| [500](#get-metrics-500) | Internal Server Error | An internal error happened while fetching prometheus. |  | [schema](#get-metrics-500-schema) |

#### Responses


##### <span id="get-metrics-200"></span> 200 - Prometheus metrics
Status: OK

###### <span id="get-metrics-200-schema"></span> Schema

##### <span id="get-metrics-500"></span> 500 - An internal error happened while fetching prometheus.
Status: Internal Server Error

###### <span id="get-metrics-500-schema"></span> Schema

### <span id="get-popular"></span> Returns the parameters corresponding to the most used request, as well as the number of hits for this request. (*GetPopular*)

```
GET /popular
```

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-popular-200) | OK | A json containing the parameters of the most used request and its count. |  | [schema](#get-popular-200-schema) |
| [500](#get-popular-500) | Internal Server Error | An internal error happened while fetching the requests data. |  | [schema](#get-popular-500-schema) |

#### Responses


##### <span id="get-popular-200"></span> 200 - A json containing the parameters of the most used request and its count.
Status: OK

###### <span id="get-popular-200-schema"></span> Schema

##### <span id="get-popular-500"></span> 500 - An internal error happened while fetching the requests data.
Status: Internal Server Error

###### <span id="get-popular-500-schema"></span> Schema

## Models
