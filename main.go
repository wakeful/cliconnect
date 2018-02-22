package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/AlecAivazis/survey.v1"
)

var (
	uri         = flag.String("url", "http://127.0.0.1:28083", "kafka connect url")
	showVersion = flag.Bool("version", false, "show version and exit")
	url         = "https://github.com/wakeful/cliconnect"
	version     = "dev"
)

type connectorList []string

type connectorStatus struct {
	Name      string `json:"name"`
	Connector struct {
		State string `json:"state"`
	} `json:"connector"`
}

type Client struct {
	Url        string
	HTTPClient *http.Client
}

func (c Client) Get(endpoint string) ([]byte, error) {

	response, err := c.HTTPClient.Get(*uri + endpoint)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	output, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (c Client) Call(method, connector, action string) error {

	req, err := http.NewRequest(method, c.Url+"/connectors/"+connector+"/"+action, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Close = true

	_, err = c.HTTPClient.Do(req)

	return err
}

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("cliconnect\n url: %s\n version: %s\n", url, version)
		os.Exit(2)
	}

	client := &Client{
		Url:        *uri,
		HTTPClient: &http.Client{},
	}

	outputReqConnectors, err := client.Get("/connectors")
	if err != nil {
		log.Fatalln(err)
	}

	var dataConnectorList connectorList
	if err := json.Unmarshal(outputReqConnectors, &dataConnectorList); err != nil {
		log.Fatalln(err)
	}

	if len(dataConnectorList) < 1 {
		log.Println("could not find any connectors")
		os.Exit(2)
	}

	var selectConnector = []*survey.Question{
		{
			Name: "connector",
			Prompt: &survey.Select{
				Message: "Select connector:",
				Options: dataConnectorList,
				Default: "",
			},
		},
	}

	selectedConnector := struct {
		Connector string `survey:"connector"`
	}{}

	if err = survey.Ask(selectConnector, &selectedConnector); err != nil {
		fmt.Println(err.Error())
	}

	outputConnectorStatus, err := client.Get("/connectors/" + *&selectedConnector.Connector + "/status")
	if err != nil {
		log.Fatalln(err)
	}

	var dataConnectorStatus connectorStatus
	if err := json.Unmarshal(outputConnectorStatus, &dataConnectorStatus); err != nil {
		log.Fatalln(err)
	}

	var connectorActionOptions = []string{"restart"}
	if strings.ToLower(dataConnectorStatus.Connector.State) == "running" {
		connectorActionOptions = append(connectorActionOptions, "pause")
	} else {
		connectorActionOptions = append(connectorActionOptions, "start")
	}

	var selectAction = []*survey.Question{
		{
			Name: "action",
			Prompt: &survey.Select{
				Message: "Select action for connector:" + selectedConnector.Connector,
				Options: connectorActionOptions,
				Default: "",
			},
		},
	}

	selectedAction := struct {
		Action string `survey:"action"`
	}{}

	if err = survey.Ask(selectAction, &selectedAction); err != nil {
		log.Fatalln(err)
	}

	switch selectedAction.Action {
	case "pause":
		err = client.Call("PUT", selectedConnector.Connector, "pause")
		break
	case "start":
		err = client.Call("PUT", selectedConnector.Connector, "resume")
		break
	case "restart":
		err = client.Call("POST", selectedConnector.Connector, "restart")
		break
	}

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("All done")

}
