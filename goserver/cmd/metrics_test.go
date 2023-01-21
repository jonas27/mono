package main

import (
	"testing"
)

func TestCustomMetrics(t *testing.T) {
	metricsServer()
	serve(":8080", "../static")
}
