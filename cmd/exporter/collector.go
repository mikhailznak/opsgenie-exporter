package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"strings"
	"time"
)

type MetricsCollector struct {
	client                    *Opsgenie
	opsgenieAlertMetricsCount *prometheus.Desc
}

// Configure metrics
func newMetricsCollector(opsgenieClient *Opsgenie) *MetricsCollector {
	return &MetricsCollector{
		client: opsgenieClient,
		opsgenieAlertMetricsCount: prometheus.NewDesc(
			"opsgenie_alert_created_total",
			"How many alerts were created during all time",
			[]string{"team", "status", "priority", "type"},
			nil,
		),
	}
}

// Implement Prometheus Interfaces and getting metrics values for MetricsCollector
func (col *MetricsCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- col.opsgenieAlertMetricsCount
}
func (col *MetricsCollector) Collect(metrics chan<- prometheus.Metric) {
	metricsList := col.procAlertCount()
	for _, metric := range metricsList {
		metrics <- metric
	}
}
func (col *MetricsCollector) procAlertCount() []prometheus.Metric {
	// Get Opsgenie Teams
	teamList := strings.Split(*teams, ",")

	// Get priorities
	priorityList := strings.Split(*priorities, ",")

	// Get status
	statusList := strings.Split(*statuses, ",")

	// Get types
	typeList := strings.Split(*filterByType, ",")

	var promMetrics []prometheus.Metric
	for _, team := range teamList {
		for _, status := range statusList {
			for _, priority := range priorityList {
				for _, typeParam := range typeList {
					labels := []string{team, status, priority, typeParam}
					promMetrics = append(promMetrics, prometheus.MustNewConstMetric(
						col.opsgenieAlertMetricsCount,
						prometheus.CounterValue,
						col.getOpsgenieAlertCount(team, status, priority, typeParam),
						labels...,
					))
					time.Sleep(time.Millisecond * time.Duration(*pauseBetweenOpsgenieRequests))
				}
			}
		}
	}
	return promMetrics
}
func (col *MetricsCollector) getOpsgenieAlertCount(teams string, status string, priority string, typeParam string) float64 {
	// Configure query parameters
	queryResponders := getOpsgenieQueryParameter("teams", teams)
	queryStatus := getOpsgenieQueryParameter("status", status)
	queryPriority := getOpsgenieQueryParameter("priority", priority)
	queryType := getOpsgenieQueryParameter("type", typeParam)
	query := fmt.Sprintf("%s %s %s %s", queryResponders, queryStatus, queryPriority, queryType)
	value, err := col.client.GetOpsgenieAlertMetricsCreatedTotal(query)
	if err != nil {
		log.Printf(
			"Error during getting opsgenie_alert_created_total %s. Team: %s, Status: %s, Priority: %s, LabelType: %s",
			err, teams, status, priority, typeParam)
		return 0.0
	}
	return value
}
func getOpsgenieQueryParameter(name string, arg string) string {
	queryArg := ""
	if arg != "all" && arg != "" {
		queryArg = fmt.Sprintf("%s:%s", name, arg)
	} else if arg == "all" {
		queryArg = ""
	}
	return queryArg
}
