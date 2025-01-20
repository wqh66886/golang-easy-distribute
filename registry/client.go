package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func RegisterService(reg Registration) error {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	err := enc.Encode(reg)
	if err != nil {
		return err
	}
	resp, err := http.Post(ServiceURL, "application/json", &buf)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register service,Registry Service response with code %s\n", resp.StatusCode)
	}

	return nil
}

func DeregisterService(reg Registration) error {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	err := enc.Encode(reg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodDelete, ServiceURL, &buf)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to deregister service,Registry Service response with code %s\n", resp.StatusCode)
	}
	log.Printf("%s Service is removed", reg.ServiceName)
	return nil
}
