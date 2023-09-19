package main

import (
	"context"
	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/team"
)

type Opsgenie struct {
	alertClient *alert.Client
	teamClient  *team.Client
}

func OpsgenieClient(apiKey string) (*Opsgenie, error) {
	alertClient, err := alert.NewClient(&client.Config{ApiKey: apiKey})
	if err != nil {
		return nil, err
	}
	teamClient, err := team.NewClient(&client.Config{ApiKey: apiKey})
	if err != nil {
		return nil, err
	}

	return &Opsgenie{
		alertClient: alertClient,
		teamClient:  teamClient,
	}, nil
}

func (cl *Opsgenie) GetOpsgenieAlertMetricsCreatedTotal(query string) (float64, error) {
	counterRequest := alert.CountAlertsRequest{Query: ""}
	if query != "" {
		counterRequest = alert.CountAlertsRequest{Query: query}
	}
	response, err := cl.alertClient.CountAlerts(context.Background(), &counterRequest)
	if err != nil {
		return 0.0, err
	}
	return float64(response.Count), nil
}

func (cl *Opsgenie) GetOpsgenieTeams() ([]string, error) {
	var teamList []string

	response, err := cl.teamClient.List(context.Background(), &team.ListTeamRequest{})
	if err != nil {
		return teamList, err
	}
	for _, team := range response.Teams {
		teamList = append(teamList, team.Name)
	}
	return teamList, nil
}
