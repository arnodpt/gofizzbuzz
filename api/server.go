package http

import (
	fiber "github.com/gofiber/fiber/v2"
	fiberprometheus "gofizzbuzz/middleware"
)

// NewServer initialize a server
func NewServer(persistence bool, requestDataPath string) *fiber.App {
	// initialize backend
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	// init prometheus middleware
	prometheus := fiberprometheus.New("fizzbuzz", persistence, requestDataPath)
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.MiddlewareQueryParamsStats)

	// initialize api with prometheus
	api := newServer(prometheus)

	if persistence {
		// store data routine
		go func() {
			prometheus.StoreMetricsData(requestDataPath)
		}()
	}
	// handlers
	app.Get("/fizzbuzz", api.Fizzbuzz)
	app.Get("/popular", api.PopularRequest)
	return app
}
