package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
)

var (
	apiKey     = flag.String("opsgenie.apiKey", "", "Opsgenie Api Key. Permissions: Read and Configuration access")
	port       = flag.Int("port", 9110, "Exporter exposed port")
	teams      = flag.String("teams", "", "List of teams to get mertrics")
	priorities = flag.String("priorities", "P1,P2", "List of priorities to capture")
	statuses   = flag.String("statuses", "open", "List of statuses to capture")
)

func main() {
	// Configure parameters
	flag.Parse()

	//Get parameters value from envs
	apiKeyFromEnv := os.Getenv("API_KEY")
	if apiKeyFromEnv != "" {
		*apiKey = apiKeyFromEnv
	}

	// Create opsgenie client
	opsgenieClient, err := OpsgenieClient(*apiKey)
	if err != nil {
		log.Fatal(err)
	}

	// Register metrics
	metrics := newMetricsCollector(opsgenieClient)
	prometheus.MustRegister(metrics)

	// Handle web server
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
