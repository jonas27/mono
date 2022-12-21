package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	unavailable = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "image_checker",
			Name:      "unavailable",
			Help:      "The image is not available.",
		})
	withDigesAndtTag = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "image_checker",
			Name:      "with_digest_and_tag",
			Help:      "This shows the number of images with a tag and digest.",
		})
	withDigestOrTag = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "image_checker",
			Name:      "with_digest_or_tag",
			Help:      "This shows the number of images either missing a tag or the digest.",
		})
	wrongDigest = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "image_checker",
			Name:      "wrong_digest",
			Help:      "The digest and the tag are not matching.",
		})
)

func serveMetrics() {
	// https://stackoverflow.com/questions/35117993/how-to-disable-go-collector-metrics-in-prometheus-client-golang
	r := prometheus.NewRegistry()
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})

	r.MustRegister(unavailable)
	r.MustRegister(withDigesAndtTag)
	r.MustRegister(withDigestOrTag)
	r.MustRegister(wrongDigest)

	http.Handle("/metrics", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
