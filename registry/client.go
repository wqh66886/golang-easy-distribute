package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
)

// RegisterService registers a service with the registry
func RegisterService(reg Registration) error {
	serviceUpdateUrl, err := url.Parse(reg.ServiceUpdateURL)
	if err != nil {
		return err
	}
	//Register a processor to handle dependent services
	http.Handle(serviceUpdateUrl.Path, &serviceUpdateHandler{})
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	err = enc.Encode(reg)
	if err != nil {
		return err
	}

	// send request to the registry
	resp, err := http.Post(ServiceURL, "application/json", &buf)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register service,Registry Service response with code %v\n", resp.Status)
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
		return fmt.Errorf("failed to deregister service,Registry Service response with code %v\n", resp.StatusCode)
	}
	log.Printf("%s is removed", reg.ServiceName)
	return nil
}

type serviceUpdateHandler struct{}

func (suh *serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var pat patch
		err := json.NewDecoder(r.Body).Decode(&pat)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("updated received %v\n", pat)
		prov.update(pat)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

type providers struct {
	services map[ServiceName][]string
	mutex    sync.RWMutex
}

var prov = providers{
	services: make(map[ServiceName][]string),
}

func (p *providers) update(pat patch) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	for _, entry := range pat.Added {
		if _, exist := p.services[entry.Name]; !exist {
			p.services[entry.Name] = make([]string, 0)
		}
		p.services[entry.Name] = append(p.services[entry.Name], entry.URL)
	}

	for _, entry := range pat.Removed {
		if urls, exist := p.services[entry.Name]; exist {
			for i, _url := range urls {
				if _url == entry.URL {
					p.services[entry.Name] = append(urls[:i], urls[i+1:]...)
				}
			}
		}
	}
}

func (p *providers) getProviders(names ...ServiceName) ([]string, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	proves := make([]string, 0)
	for _, name := range names {
		if prov, exist := p.services[name]; !exist {
			return nil, fmt.Errorf("%s is not registered", name)
		} else {
			proves = append(proves, prov...)
		}
	}
	return proves, nil
}

func GetProviders(names ...ServiceName) ([]string, error) {
	return prov.getProviders(names...)
}
