package iapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// CreateKeyValuePairs ...
func CreateKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=%s,", key, value)
	}
	return b.String()
}

// GetHost ...
func (server *Server) GetHost(hostname string) ([]HostStruct, error) {

	var hosts []HostStruct

	results, err := server.NewAPIRequest("GET", "/objects/hosts/"+hostname, nil)
	if err != nil {
		return nil, err
	}

	// Contents of the results is an interface object. Need to convert it to json first.
	jsonStr, marshalErr := json.Marshal(results.Results)
	if marshalErr != nil {
		return nil, marshalErr
	}

	// then the JSON can be pushed into the appropriate struct.
	// Note : Results is a slice so much push into a slice.

	if unmarshalErr := json.Unmarshal(jsonStr, &hosts); unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return hosts, err

}

// GetHostTr ...
func (server *Server) GetHostTr(hostname string) ([]HostStructRead, error) {

	var hosts []HostStructRead

	results, err := server.NewAPIRequest("GET", "/objects/hosts/"+hostname, nil)
	if err != nil {
		return nil, err
	}

	// Contents of the results is an interface object. Need to convert it to json first.
	jsonStr, marshalErr := json.Marshal(results.Results)
	if marshalErr != nil {
		return nil, marshalErr
	}

	// then the JSON can be pushed into the appropriate struct.
	// Note : Results is a slice so much push into a slice.

	if unmarshalErr := json.Unmarshal(jsonStr, &hosts); unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return hosts, err

}

// SetNotifications ...
func (server *Server) SetNotifications(hostname, enableNotifications string) error {
	eNotification, err := strconv.ParseBool(enableNotifications)

	strBody := fmt.Sprintf(`{ "filter": "regex(hostname, host.name)", "filter_vars": { "hostname": "%s" }, "attrs": { "enable_notifications": %t } }`, hostname, eNotification)

	// Make the API request to create the hosts.
	notificationResp, err := server.NewAPIRequest("POST", "/objects/services", []byte(strBody))
	if err != nil {
		return err
	}

	if notificationResp.Code != 200 {
		return err
	}

	return nil
}

// CreateHost ...
func (server *Server) CreateHost(hostname, address, zone, checkCommand, checkPeriod, actionURL, notesURL, notes, volatile, enableActiveChecks, enableNotifications string, variables map[string]string, templates []string) ([]HostStruct, error) {

	var newAttrs HostAttrs
	newAttrs.Address = address
	newAttrs.Zone = zone
	newAttrs.CheckCommand = checkCommand

	if variables != nil {
		newAttrs.Vars = Flatten(variables)
	}

	if actionURL != "" {
		newAttrs.ActionURL = actionURL
	}

	if notesURL != "" {
		newAttrs.NotesURL = notesURL
	}

	if notes != "" {
		newAttrs.Notes = notes
	}

	newAttrs.EnableNotifications = enableNotifications

	if enableActiveChecks != "" {
		newAttrs.EnableActiveChecks, _ = strconv.ParseBool(enableActiveChecks)
	}

	if volatile != "" {
		newAttrs.Volatile, _ = strconv.ParseBool(volatile)
	}

	if checkPeriod != "" {
		newAttrs.CheckPeriod = checkPeriod
	}

	newAttrsMarshalled, marshalErr := json.Marshal(newAttrs)
	if marshalErr != nil {
		return nil, marshalErr
	}

	cleanedHostAttrs := make(map[string]interface{})

	unmarshalErr := json.Unmarshal(newAttrsMarshalled, &cleanedHostAttrs)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	var newHost HostStruct
	newHost.Name = hostname
	newHost.Type = "Host"
	newHost.Templates = templates
	newHost.Attrs = Flatten(cleanedHostAttrs)

	// Create JSON from completed struct
	payloadJSON, marshalErr := json.Marshal(newHost)
	if marshalErr != nil {
		return nil, marshalErr
	}

	fmt.Printf("<payload> %s\n", payloadJSON)

	// Make the API request to create the hosts.
	results, err := server.NewAPIRequest("PUT", "/objects/hosts/"+hostname, []byte(payloadJSON))
	if err != nil {
		return nil, err
	}

	if results.Code == 200 {
		hosts, err := server.GetHost(hostname)
		return hosts, err
	}

	return nil, fmt.Errorf("%s", results.ErrorString)

}

// UpdateHost ...
func (server *Server) UpdateHost(hostname, address, zone, checkCommand, checkPeriod, actionURL, notesURL, notes, volatile, enableActiveChecks, enableNotifications string, variables map[string]string, templates []string) ([]HostStruct, error) {

	var newAttrs HostAttrs
	newAttrs.Address = address
	newAttrs.Zone = zone
	newAttrs.CheckCommand = checkCommand

	if variables != nil {
		newAttrs.Vars = Flatten(variables)
	}

	if actionURL != "" {
		newAttrs.ActionURL = actionURL
	}

	if notesURL != "" {
		newAttrs.NotesURL = notesURL
	}

	if notes != "" {
		newAttrs.Notes = notes
	}

	newAttrs.EnableNotifications = enableNotifications

	if enableActiveChecks != "" {
		newAttrs.EnableActiveChecks, _ = strconv.ParseBool(enableActiveChecks)
	}

	if volatile != "" {
		newAttrs.Volatile, _ = strconv.ParseBool(volatile)
	}

	if checkPeriod != "" {
		newAttrs.CheckPeriod = checkPeriod
	}

	newAttrsMarshalled, marshalErr := json.Marshal(newAttrs)
	if marshalErr != nil {
		return nil, marshalErr
	}

	cleanedHostAttrs := make(map[string]interface{})

	unmarshalErr := json.Unmarshal(newAttrsMarshalled, &cleanedHostAttrs)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	var newHost HostStruct
	newHost.Name = hostname
	newHost.Type = "Host"
	newHost.Templates = templates
	newHost.Attrs = Flatten(cleanedHostAttrs)

	// Create JSON from completed struct
	payloadJSON, marshalErr := json.Marshal(newHost)
	if marshalErr != nil {
		return nil, marshalErr
	}

	//fmt.Printf("<payload> %s\n", payloadJSON) // for debugging purposes

	// Make the API request to create the hosts.
	results, err := server.NewAPIRequest("POST", "/objects/hosts/"+hostname, []byte(payloadJSON))
	if err != nil {
		return nil, err
	}

	if results.Code == 200 {
		hosts, err := server.GetHost(hostname)
		return hosts, err
	}

	return nil, fmt.Errorf("%s", results.ErrorString)

}

// DeleteHost ...
func (server *Server) DeleteHost(hostname string) error {

	results, err := server.NewAPIRequest("DELETE", "/objects/hosts/"+hostname+"?cascade=1", nil)
	if err != nil {
		return err
	}

	if results.Code == 200 {
		return nil
	} else {
		return fmt.Errorf("%s", results.ErrorString)
	}
}
