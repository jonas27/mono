package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	port    = flag.String("port", ":8080", "The port")
	metrics = flag.Bool("metrics", true, "Enable Prometheus metrics.")
	secure  = flag.Bool("secure", false, "Should the server use TLS.")

	httpRequestsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Count of all HTTP requests",
	})
)

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile)
	metricsServer()
	httpRequestsTotal.Add(10)

	log.Println("Start unencrypted server")
	log.Fatal(http.ListenAndServe(*port, nil))
}

func metricsServer() {
	log.Println("registering metrics")
	r := prometheus.NewRegistry()
	r.MustRegister(httpRequestsTotal)
	http.Handle("/customMetrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	prometheus.MustRegister(httpRequestsTotal)
	http.Handle("/metrics", promhttp.Handler())
}
