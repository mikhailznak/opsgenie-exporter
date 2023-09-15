package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"strings"
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
			[]string{"team", "status", "priority"},
			nil,
		),
	}
}

// Implement Prometheus Interfaces and getting metrics values for MetricsCollector
func (col *MetricsCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- col.opsgenieAlertMetricsCount
}
func (col *MetricsCollector) Collect(metrics chan<- prometheus.Metric) {
	// Get Opsgenie Teams
	teamList := []string{}
	if *teams == "all" {
		response, err := col.client.GetOpsgenieTeams()
		if err != nil {
			log.Printf("Error get opsgenie team list: %s", err)
			teamList = []string{}
		} else {
			teamList = response
		}
	} else {
		teamList = strings.Split(*teams, ",")
	}

	// Get priorities
	priorityList := []string{}
	if *priorities == "all" {
		priorityList = []string{"all", "P1", "P2", "P3", "P4", "P5"}
	} else {
		priorityList = strings.Split(*priorities, ",")
	}

	// Get status
	statusList := []string{}
	if *statuses == "all" {
		statusList = []string{"all", "open", "closed"}
	} else {
		statusList = strings.Split(*statuses, ",")
	}

	for _, team := range teamList {
		for _, status := range statusList {
			for _, priority := range priorityList {
				labels := []string{team, status, priority}
				metrics <- prometheus.MustNewConstMetric(
					col.opsgenieAlertMetricsCount,
					prometheus.CounterValue,
					col.procOpsgenieAlertMetricsTotal(team, status, priority),
					labels...,
				)
			}
		}
	}
}
func (col *MetricsCollector) procOpsgenieAlertMetricsTotal(responders string, status string, priority string) float64 {
	// Configure reponders query parameters
	queryResponders := fmt.Sprintf("")
	if responders != "all" && responders != "" {
		queryResponders = fmt.Sprintf("responders:%s", responders)
	}

	// Configure status query parameters
	queryStatus := fmt.Sprintf("")
	if status != "all" && status != "" {
		queryStatus = fmt.Sprintf("status:%s", status)
	}

	// Configure priority query parameters
	queryPriority := fmt.Sprintf("")
	if priority != "all" && priority != "" {
		queryPriority = fmt.Sprintf("priority:%s", priority)
	}

	query := fmt.Sprintf("%s %s %s", queryResponders, queryStatus, queryPriority)
	value, err := col.client.GetOpsgenieAlertMetricsCreatedTotal(query)
	if err != nil {
		log.Printf("Error during getting opsgenie_alert_created_total. Team: %s, Status: %s", err, responders, status)
		return 0.0
	}
	return value
}
