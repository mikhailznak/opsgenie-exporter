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
	apiKey = flag.String("opsgenie.apiKey", "", "Opsgenie Api Key. Permissions: Read and Configuration access")
	port   = flag.Int("port", 9110, "Exporter exposed port")
	teams  = flag.String(
		"teams",
		"",
		"List of teams to get metrics. "+
			"`all` value means that will be received metric without querying `teams` label")
	priorities = flag.String(
		"priorities",
		"P1,P2",
		"List of priorities to capture. "+
			"`all` value means that will be received metric without querying `priority` label")
	statuses = flag.String(
		"statuses",
		"open",
		"List of statuses to capture. "+
			"`all` value means that will be received metric without querying `status` label")
	filterByType = flag.String(
		"types",
		"",
		"List of extra properties `type` to count alerts by specific type. "+
			"`all` value means that will be received metric without querying `type` label")
	pauseBetweenOpsgenieRequests = flag.Int(
		"pause",
		1,
		"Pause between Opsgenie requests which exporter makes to get metrics count."+
			"Time is in milliseconds")
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
