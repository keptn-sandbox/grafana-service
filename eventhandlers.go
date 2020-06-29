package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	keptn "github.com/keptn/go-utils/pkg/lib"
)

/**
* Here are all the handler functions for the individual event
  See https://github.com/keptn/spec/blob/0.1.3/cloudevents.md for details on the payload

  -> "sh.keptn.event.configuration.change"
  -> "sh.keptn.events.deployment-finished"
  -> "sh.keptn.events.tests-finished"
  -> "sh.keptn.event.start-evaluation"
  -> "sh.keptn.events.evaluation-done"
  -> "sh.keptn.event.problem.open"
	-> "sh.keptn.events.problem"
	-> "sh.keptn.event.action.triggered"
	-> "sh.keptn.event.monitoring.configure"
*/

var grafanaAPIDatasourceURL = os.Getenv("GRAFANA_URL") + "/api/datasources"
var grafanaAPIDashboardURL = os.Getenv("GRAFANA_URL") + "/api/dashboards/db"
var bearer = "Bearer " + os.Getenv("GRAFANA_TOKEN")

// Handles ConfigureMonitoringEventType = "sh.keptn.event.monitoring.configure"
func HandleConfigureMonitoringEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.ConfigureMonitoringEventData) error {
	log.Printf("Handling Configure Monitoring Event: %s", incomingEvent.Context.GetID())

	dataprovider := data.Type

	if os.Getenv("GRAFANA_URL") == "" {
		return errors.New("Grafana endpoint not set")
	}
	if os.Getenv("GRAFANA_TOKEN") == "" {
		return errors.New("Grafana API token not set")
	}

	keptnHandler, err := keptn.NewKeptn(&incomingEvent, keptn.KeptnOpts{ConfigurationServiceURL: "http://localhost:9090"})
	shipyard, err := keptnHandler.GetShipyard()
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	// create grafana datasource for prometheus
	if os.Getenv("PROMETHEUS_URL") == "" {
		log.Printf("PROMETHEUS_URL not set, therefore not creating datasource in Grafana")
	} else if dataprovider == "prometheus" {
		log.Printf("PROMETHEUS_URL set. Configuring datasource in Grafana.")

		values := map[string]string{"name": "Prometheus", "type": "prometheus", "url": os.Getenv("PROMETHEUS_URL"), "access": "proxy"}
		jsonValue, _ := json.Marshal(values)

		log.Printf("Sending request to Grafana: %s", grafanaAPIDatasourceURL)
		req, err := http.NewRequest("POST", grafanaAPIDatasourceURL, bytes.NewBuffer(jsonValue))
		req.Header.Add("Authorization", bearer)
		req.Header.Add("Content-Type", "application/json")

		if err != nil {
			log.Printf("Could not create request: %s", err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error on response.\n[ERRO] -", err)
		}
		log.Printf("Creating of datasource HTTP response: %d", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println(string([]byte(body)))

	}

	// CREATE DASHBOARD
	if dataprovider == "prometheus" {
		// TODO GENERATE JSON PAYLOAD

		var s []string
		for _, stage := range shipyard.Stages {
			s = append(s, `"`+stage.Name+`"`)
		}
		var myStages = "[" + strings.Join(s, ", ") + "]"
		log.Println(myStages)

		// output, err := exec.Command("jsonnet", "-J", "./grafonnet-lib/", "keptn.jsonnet", ">", "gen/keptn.json").Output()
		output, err := exec.Command("jsonnet",
			"-J", "./grafonnet-lib/",
			"--ext-code", "stages="+myStages+"",
			"--ext-str", "service="+data.Service+"",
			"--ext-str", "project="+data.Project+"",
			"keptn.jsonnet").Output()
		if err != nil {
			log.Printf("error executing jsonnet: %s", err.Error())
			return err
		}
		//stringOutput := string(output)
		var data map[string]interface{}
		err = json.Unmarshal(output, &data)
		if err != nil {
			log.Println("could not unmarshal json")
			return err
		}
		// fmt.Println(data)

		values := map[string]interface{}{"dashboard": data, "overwrite": true}
		jsonValue, _ := json.Marshal(values)
		log.Println(string(jsonValue))

		log.Printf("Sending request to Grafana: %s", grafanaAPIDashboardURL)
		req, err := http.NewRequest("POST", grafanaAPIDashboardURL, bytes.NewBuffer(jsonValue))
		req.Header.Add("Authorization", bearer)
		req.Header.Add("Content-Type", "application/json")

		if err != nil {
			log.Printf("Could not create request: %s", err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error on response.\n[ERRO] -", err)
		}
		log.Printf("Creating of dashboad HTTP response: %d", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println(string([]byte(body)))

	} else {
		log.Printf("type %s not supported", dataprovider)
	}

	return nil
}

//
// Handles ConfigurationChangeEventType = "sh.keptn.event.configuration.change"
// TODO: add in your handler code
//
func HandleConfigurationChangeEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.ConfigurationChangeEventData) error {
	log.Printf("Handling Configuration Changed Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles DeploymentFinishedEventType = "sh.keptn.events.deployment-finished"
// TODO: add in your handler code
//
func HandleDeploymentFinishedEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.DeploymentFinishedEventData) error {
	log.Printf("Handling Deployment Finished Event: %s", incomingEvent.Context.GetID())

	// capture start time for tests
	// startTime := time.Now()

	// run tests
	// ToDo: Implement your tests here

	// Send Test Finished Event
	// return myKeptn.SendTestsFinishedEvent(&incomingEvent, "", "", startTime, "pass", nil, "grafana-service")
	return nil
}

//
// Handles TestsFinishedEventType = "sh.keptn.events.tests-finished"
// TODO: add in your handler code
//
func HandleTestsFinishedEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.TestsFinishedEventData) error {
	log.Printf("Handling Tests Finished Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles EvaluationDoneEventType = "sh.keptn.events.evaluation-done"
// TODO: add in your handler code
//
func HandleStartEvaluationEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.StartEvaluationEventData) error {
	log.Printf("Handling Start Evaluation Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles DeploymentFinishedEventType = "sh.keptn.events.deployment-finished"
// TODO: add in your handler code
//
func HandleEvaluationDoneEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.EvaluationDoneEventData) error {
	log.Printf("Handling Evaluation Done Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles ProblemOpenEventType = "sh.keptn.event.problem.open"
// Handles ProblemEventType = "sh.keptn.events.problem"
// TODO: add in your handler code
//
func HandleProblemEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.ProblemEventData) error {
	log.Printf("Handling Problem Event: %s", incomingEvent.Context.GetID())

	// Deprecated since Keptn 0.7.0 - use the HandleActionTriggeredEvent instead

	return nil
}

//
// Handles ActionTriggeredEventType = "sh.keptn.event.action.triggered"
// TODO: add in your handler code
//
func HandleActionTriggeredEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.ActionTriggeredEventData) error {
	log.Printf("Handling Action Triggered Event: %s", incomingEvent.Context.GetID())

	// check if action is supported
	if data.Action.Action == "action-xyz" {
		//myKeptn.SendActionStartedEvent() TODO: implement the SendActionStartedEvent in keptn/go-utils/pkg/lib/events.go

		// Implement your remediation action here

		//myKeptn.SendActionFinishedEvent() TODO: implement the SendActionFinishedEvent in keptn/go-utils/pkg/lib/events.go
	}
	return nil
}
