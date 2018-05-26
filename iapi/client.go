// Package iapi provides a client for interacting with an Icinga2 Server
package iapi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Server ... Use to be ClientConfig
type Server struct {
	Username           string
	Password           string
	BaseURL            string
	AllowUnverifiedSSL bool
	httpClient         *http.Client
}

// func New ...
func New(username, password, url string, allowUnverifiedSSL bool) (*Server, error) {
	return &Server{username, password, url, allowUnverifiedSSL, nil}, nil
}

// func Config ...
func (server *Server) Config(username, password, url string, allowUnverifiedSSL bool) (*Server, error) {

	// TODO : Add code to verify parameters
	return &Server{username, password, url, allowUnverifiedSSL, nil}, nil

}

func (server *Server) Connect() error {

	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: server.AllowUnverifiedSSL,
		},
	}

	server.httpClient = &http.Client{
		Transport: t,
		Timeout:   time.Second * 60,
	}

	request, err := http.NewRequest("GET", server.BaseURL, nil)
	if err != nil {
		server.httpClient = nil
	}

	request.SetBasicAuth(server.Username, server.Password)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	response, err := server.httpClient.Do(request)

	if (err != nil) || (response == nil) {
		server.httpClient = nil
		return err
	}

	defer response.Body.Close()

	return nil

}

// NewAPIRequest ...
func (server *Server) NewAPIRequest(method, APICall string, jsonString []byte) (*APIResult, error) {

	var results APIResult

	fullURL := server.BaseURL + APICall

	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: server.AllowUnverifiedSSL,
		},
	}

	server.httpClient = &http.Client{
		Transport: t,
		Timeout:   time.Second * 60,
	}

	request, requestErr := http.NewRequest(method, fullURL, bytes.NewBuffer(jsonString))
	if requestErr != nil {
		return nil, requestErr
	}

	request.SetBasicAuth(server.Username, server.Password)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	response, doErr := server.httpClient.Do(request)
	if doErr != nil {
		results.Code = 0
		results.Status = "Error : Request to server failed : " + doErr.Error()
		results.ErrorString = doErr.Error()
		return &results, doErr
	}

	defer response.Body.Close()

	if decodeErr := json.NewDecoder(response.Body).Decode(&results); decodeErr != nil {
		return nil, decodeErr
	}

	if results.Code == 0 { // results.Code has default value so set it.
		results.Code = response.StatusCode
	}

	if results.Status == "" { // results.Status has default value, so set it.
		results.Status = response.Status
	}

	switch results.Code {
	case 0:
		results.ErrorString = "Did not get a response code."
	case 404:
		results.ErrorString = results.Status
	case 200:
		results.ErrorString = results.Status
	default:
		theError := strings.Replace(results.Results.([]interface{})[0].(map[string]interface{})["errors"].([]interface{})[0].(string), "\n", " ", -1)
		results.ErrorString = strings.Replace(theError, "Error: ", "", -1)

	}

	return &results, nil

}

// NewAPIRequestFiltered ...
func (server *Server) NewAPIRequestFiltered(method, APICall string, jsonString []byte) (string, error) {

	var hosts FilteredHostResults

	fullURL := server.BaseURL + APICall

	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: server.AllowUnverifiedSSL,
		},
	}

	server.httpClient = &http.Client{
		Transport: t,
		Timeout:   time.Second * 60,
	}

	request, requestErr := http.NewRequest(method, fullURL, bytes.NewBuffer(jsonString))
	if requestErr != nil {
		fmt.Printf("The HTTP request failed with error %s\n", requestErr)
	}

	request.SetBasicAuth(server.Username, server.Password)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-HTTP-Method-Override", "GET")

	response, _ := server.httpClient.Do(request)

	data, _ := ioutil.ReadAll(response.Body)
	getError := json.Unmarshal(data, &hosts)

	// fmt.Println(len(hosts.Results))

	if getError != nil {

		return `Host not found: `, getError
	}

	if len(hosts.Results) != 0 {
		hostName := hosts.Results[0].Name
		return hostName, nil
	}

	return `Hosst not found`, fmt.Errorf("API Error: host not found")

}
