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

	//keptnHandler, err := keptn.NewKeptn(&incomingEvent, keptn.KeptnOpts{ConfigurationServiceURL: "http://localhost:9090"})
	keptnHandler, err := keptn.NewKeptn(&incomingEvent, keptn.KeptnOpts{})
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

		// payload for creating a datasource in Grafana
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

		var s []string
		for _, stage := range shipyard.Stages {
			s = append(s, `"`+stage.Name+`"`)
		}
		var myStages = "[" + strings.Join(s, ", ") + "]"
		log.Println(myStages)

		// output, err := exec.Command("jsonnet", "-J", "./grafonnet-lib-local/", "keptn.jsonnet", ">", "gen/keptn.json").Output()
		output, err := exec.Command("jsonnet",
			"-J", "./grafonnet-lib-local/",
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
