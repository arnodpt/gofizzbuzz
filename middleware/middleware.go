package fiberprometheus

import (
	"github.com/gofiber/adaptor/v2"
	fiber "github.com/gofiber/fiber/v2"
	"strconv"
	"github.com/rs/zerolog/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
    dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"time"
	"os"
)

// FiberPrometheus
type FiberPrometheus struct {
	RequestsTotal   *prometheus.CounterVec
	url            string
}

// Initialize prometheus counter
func create(serviceName, namespace, subsystem string, persistence bool, requestDataPath string) *FiberPrometheus {
	constLabels := make(prometheus.Labels)
	if serviceName != "" {
		constLabels["service"] = serviceName
	}

	// Init http requests counter
	counter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:        prometheus.BuildFQName(namespace, subsystem, "requests_total"),
			Help:        "Count all http requests by status code, method, path and query parameters.",
			ConstLabels: constLabels,
		},
		[]string{"status_code", "method", "path", "str1", "str2", "int1", "int2", "limit"},
	)

	// Load metrics from data file storage if persistence is activated
	if persistence {
		mf, err := loadMetricsData(requestDataPath)
		if err == nil {
			for _, v := range mf {
				metrics := v.GetMetric()
				for _, m := range metrics {
					labels := make(map[string]string)
					for _, label := range m.Label {
						if (label.Name != nil && label.Value != nil && *label.Name != "service") {
							labels[*label.Name] = *label.Value
						}
					}
					if m.Counter != nil && m.Counter.Value != nil {
						counter.With(labels).Add(*m.Counter.Value)
					}
				}
			}
		}
	}

	return &FiberPrometheus{
		RequestsTotal:   counter,
	}
}

// New creates a new instance of FiberPrometheus middleware
// serviceName is available as a const label
func New(serviceName string, persistence bool, requestDataPath string) *FiberPrometheus {
	return create(serviceName, "http", "", persistence, requestDataPath)
}

// RegisterAt will register the prometheus handler at a given URL
func (ps *FiberPrometheus) RegisterAt(app *fiber.App, url string) {
	ps.url = url
	app.Get(url, adaptor.HTTPHandler(promhttp.Handler()))
}

// FizzbuzzQueryParams are the accepted query parameters for a Fizzbuzz query
type FizzbuzzQueryParams struct {
	Str1  *string `query:str1`
	Str2  *string `query:str2`
	Int1  *int    `query:int1`
	Int2  *int    `query:int2`
	Limit *int    `query:limit`
}

// Middleware to count requests according to their query parameters
func (ps *FiberPrometheus) MiddlewareQueryParamsStats(c *fiber.Ctx) error {

	method := c.Route().Method

	err := c.Next()
	
	// Don't process prometheus-related requests
	if c.Route().Path == ps.url || c.Route().Path == "/popular" {
		return err
	}

	// initialize with default error code
	status := fiber.StatusInternalServerError
	if err != nil {
		if e, ok := err.(*fiber.Error); ok {
			// Get correct error code from fiber.Error type
			status = e.Code
		}
	} else {
		status = c.Response().StatusCode()
	}

	path := c.Route().Path
	statusCode := strconv.Itoa(status)

	//Parsing parameters
	var queryParameters FizzbuzzQueryParams
	c.QueryParser(&queryParameters)

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

	// Increments a request total metric associated with the query parameters
	ps.RequestsTotal.WithLabelValues(statusCode, method, path, str1, str2, strconv.Itoa(int1), strconv.Itoa(int2), strconv.Itoa(limit)).Inc()

	return err
}

// Store requests total counter every second if persistence is activated
func (ps *FiberPrometheus) StoreMetricsData(path string) {
	for {
    	time.Sleep(1 * time.Second)
		// Gather prometheus requests_total metrics
		reg := prometheus.NewRegistry()
		reg.MustRegister(ps.RequestsTotal)
		mf, err := reg.Gather()
		if err != nil {
			log.Error().Stack().Err(err).Msg("Prometheus gather error")
			continue;
		}

		// Write gathering in storage file
		if (len(mf) > 0) {
			f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0660)

			if err != nil {
				log.Error().Stack().Err(err).Msg("Prometheus store error")
				continue
			}
			defer f.Close()
		
			if _, err := expfmt.MetricFamilyToText(f, mf[0]); err != nil {
				log.Error().Stack().Err(err).Msg("Prometheus store error")
				continue;
			}
		}
	}
}

// Load requests total counter from storage file if persistence is activated
func loadMetricsData(path string) (map[string]*dto.MetricFamily, error) {
	data, err := os.Open(path)
	if err != nil {
		log.Info().Msg("No request data file storage available.")
		return nil, err
	}
	defer data.Close()
    var parser expfmt.TextParser
    mf, err := parser.TextToMetricFamilies(data)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error parsing request data file storage.")
		return nil, err
	}
	return mf, err
}