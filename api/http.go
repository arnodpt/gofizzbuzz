package http

import (
	"encoding/json"
	"strconv"

	"github.com/rs/zerolog/log"
	fiber "github.com/gofiber/fiber/v2"
	fiberprometheus "gofizzbuzz/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

type Server struct {
	Prom *fiberprometheus.FiberPrometheus
}

var (
	invalidQueryParameters      = []byte(`{"error":"invalid query parameters"}`)
	negativeLimit               = []byte(`{"error":"limit should be positive"}`)
	internalJsonError           = []byte(`{"error":"internal json error"}`)
	internalPrometheusError     = []byte(`{"error":"cannot gather prometheus"}`)
	notRequestsFound            = []byte(`{"msg":"no requests found"}`)
)

// newServer creates the server object handling the api calls
func newServer(prom *fiberprometheus.FiberPrometheus) *Server {
	return &Server{Prom : prom}
}

// Fizzbuzz returns a list of strings with numbers from 1 to limit, where:
// all multiples of int1 are replaced by str1
// all multiples of int2 are replaced by str2
// all multiples of int1 and int2 are replaced by str1str2.
func (s *Server) Fizzbuzz(c *fiber.Ctx) error {
	//Parsing parameters
	var queryParameters FizzbuzzQueryParams
	err := c.QueryParser(&queryParameters)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Fizzbuzz error")
		c.Context().SetStatusCode(400)
		return c.Send(invalidQueryParameters)
	}

	// str1 has fizz as a default value
	var str1 string
	if queryParameters.Str1 != nil {
		str1 = *queryParameters.Str1
	} else {
		str1 = "fizz"
	}

	// str2 has buzz as a default value
	var str2 string
	if queryParameters.Str2 != nil {
		str2 = *queryParameters.Str2
	} else {
		str2 = "buzz"
	}

	// int1 has 3 as a default value
	var int1 int
	if queryParameters.Int1 != nil {
		int1 = *queryParameters.Int1
	} else {
		int1 = 3
	}

	// int2 has 3 as a default value
	var int2 int
	if queryParameters.Int2 != nil {
		int2 = *queryParameters.Int2
	} else {
		int2 = 5
	}

	// limit has 50 as a default value
	var limit int
	if queryParameters.Limit != nil {
		limit = *queryParameters.Limit
	} else {
		limit = 50
	}
	if limit < 0 {
		c.Context().SetStatusCode(400)
		return c.Send(negativeLimit)
	}

	// Actual fizzbuzz algorithm
	var res []string
	for i := 1; i <= limit; i++ {
		if i%int1 == 0 && i%int2 == 0 {
			res = append(res, str1+str2)
		} else if i%int1 == 0 {
			res = append(res, str1)
		} else if i%int2 == 0 {
			res = append(res, str2)
		} else {
			res = append(res, strconv.Itoa(i))
		}
	}
	
	// Filling result
	var fizzbuzzResult FizzbuzzResult
	fizzbuzzResult.Result = res
	resBytes, err := json.Marshal(fizzbuzzResult)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Fizzbuzz error")
		c.Context().SetStatusCode(500)
		return c.Send([]byte(internalJsonError))
	}

	return c.Send(resBytes)
}

// statistics endpoint allowing users to know what the most frequent request has been. This endpoint should:
// Return the parameters corresponding to the most used request, as well as the number of hits for this request"
func (s *Server) PopularRequest(c *fiber.Ctx) error {
	// Gather prometheus requests_total metrics
	reg := prometheus.NewRegistry()
	reg.MustRegister(s.Prom.RequestsTotal)
	mf, err := reg.Gather()
	if (err != nil) {
		log.Error().Stack().Err(err).Msg("Prometheus gather error")
		c.Context().SetStatusCode(500)
		return c.Send([]byte(internalPrometheusError))
	} else {
		// find most used requests
		var result PopularResult
		result.Count = 0
		for _, v := range mf {
			metrics := v.GetMetric()
			for _, m := range metrics {
				str1_present := false
				str2_present := false
				int1_present := false
				int2_present := false
				limit_present := false
				str1 := ""
				str2 := ""
				int1 := 0
				int2 := 0
				limit := 0
				for _, label := range m.Label {
					if (label.Name == nil || label.Value == nil) {
						continue
					}
					if (*label.Name == "str1") {
						str1_present = true
						str1 = *label.Value
					} else if (*label.Name == "str2") {
						str2_present = true
						str2 = *label.Value
					} else if (*label.Name == "int1") {
						int1_present = true
						int1, err = strconv.Atoi(*label.Value)
						if err != nil {
							continue
						}
					} else if (*label.Name == "int2") {
						int2_present = true
						int2, err = strconv.Atoi(*label.Value)
						if err != nil {
							continue
						}
					} else if (*label.Name == "limit") {
						limit_present = true
						limit, err = strconv.Atoi(*label.Value)
						if err != nil {
							continue
						}
					}
				}
				if (str1_present && str2_present && int1_present && int2_present && limit_present && m.Counter != nil && m.Counter.Value != nil && int64(*m.Counter.Value) > result.Count) {
					result.Count = int64(*m.Counter.Value)
					result.Str1 = str1
					result.Str2 = str2
					result.Int1 = int1
					result.Int2 = int2
					result.Limit = limit
				}
			}
		}
		if result.Count > 0 {
			resBytes, err := json.Marshal(result)
			if err != nil {
				log.Error().Stack().Err(err).Msg("Popular request error")
				c.Context().SetStatusCode(500)
				return c.Send([]byte(internalJsonError))
			}
			return c.Send(resBytes)		
		} else {
			c.Context().SetStatusCode(200)
			return c.Send([]byte(notRequestsFound))
		}
	}
}