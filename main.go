package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/hidez8891/nature-remo-prometheus/natureremo"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type config struct {
	Port     int    `env:"PORT" envDefault:"8080"`
	Host     string `env:"HOST"`
	ApiToken string `env:"NATURE_REMO_TOKEN,required"`
	ApiUrl   string `env:"NATURE_REMO_API_URL" envDefault:"https://api.nature.global"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v", err)
		return
	}

	registry := prometheus.NewRegistry()

	collector := natureremo.NewCollector(cfg.ApiToken, cfg.ApiUrl)
	registry.MustRegister(collector)

	e := echo.New()

	e.GET("/metrics", echo.WrapHandler(
		promhttp.HandlerFor(registry, promhttp.HandlerOpts{})),
	)

	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello")
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)))
}
