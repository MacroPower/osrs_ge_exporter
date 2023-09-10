package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MacroPower/osrs_ge_exporter/internal/collector"
	"github.com/MacroPower/osrs_ge_exporter/internal/log"
	"github.com/MacroPower/osrs_ge_exporter/internal/version"
	"github.com/MacroPower/osrs_ge_exporter/pkg/client"

	"github.com/alecthomas/kong"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const appName = "osrs_ge_exporter"

var cli struct {
	Address     string        `help:"Address to listen on for metrics." env:"ADDRESS" default:":8080"`
	MetricsPath string        `help:"Path under which to expose metrics." env:"METRICS_PATH" default:"/metrics"`
	Timeout     time.Duration `help:"HTTP timeout." type:"time.Duration" env:"TIMEOUT" default:"30s"`
	Log         struct {
		Level  string `help:"Log level." default:"info"`
		Format string `help:"Log format. One of: [logfmt, json]" default:"logfmt"`
	} `prefix:"log." embed:""`
}

func main() {
	cliCtx := kong.Parse(&cli, kong.Name(appName))

	logLevel := &log.AllowedLevel{}
	if err := logLevel.Set(cli.Log.Level); err != nil {
		cliCtx.FatalIfErrorf(err)
	}

	logFormat := &log.AllowedFormat{}
	if err := logFormat.Set(cli.Log.Format); err != nil {
		cliCtx.FatalIfErrorf(err)
	}

	logger := log.New(&log.Config{
		Level:  logLevel,
		Format: logFormat,
	})

	err := log.Info(logger).Log("msg", fmt.Sprintf("Starting %s", appName))
	cliCtx.FatalIfErrorf(err)
	err = version.LogInfo(logger)
	cliCtx.FatalIfErrorf(err)
	err = version.LogBuildContext(logger)
	cliCtx.FatalIfErrorf(err)

	mux := http.NewServeMux()

	c := client.NewOSRSPriceClient()
	metricExporter := collector.NewExporter(c, cli.Timeout, logger)
	prometheus.MustRegister(metricExporter)
	mux.Handle(cli.MetricsPath, promhttp.Handler())

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html>
			<head><title>OSRS GE Exporter</title></head>
			<body>
			<h1>OSRS GE Exporter</h1>
			<p><a href="` + cli.MetricsPath + `">Metrics</a></p>
			</body>
			</html>`))
		if err != nil {
			log.Error(logger).Log("msg", "Failed writing response", "err", err)
		}
	})

	server := &http.Server{
		Addr:              cli.Address,
		ReadTimeout:       cli.Timeout,
		ReadHeaderTimeout: cli.Timeout,
		WriteTimeout:      cli.Timeout,
		Handler:           mux,
	}

	log.Info(logger).Log("msg", "Listening", "address", cli.Address)
	if err := server.ListenAndServe(); err != nil {
		log.Error(logger).Log("msg", "HTTP server error", "err", err)
		cliCtx.Exit(1)
	}
}
