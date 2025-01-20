package registry

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

const ServicePort = ":5678"
const ServiceURL = "http://localhost" + ServicePort + "/services"

type registry struct {
	registrations []Registration
	mutex         sync.Mutex
}

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.registrations = append(r.registrations, reg)
	return nil
}

func (r *registry) remove(url string) error {
	for i, reg := range r.registrations {
		if reg.ServiceURL == url {
			r.mutex.Lock()
			r.registrations = append(r.registrations[:i], r.registrations[i+1:]...)
			r.mutex.Unlock()
			return nil
		}
	}
	return fmt.Errorf("registry service not found url: %s\n", url)
}

var reg = registry{
	registrations: make([]Registration, 0),
}

type RegistryService struct{}

// RegisterService is a http handler
func (rs RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("request received")
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding Service %s with Url %s", r.ServiceName, r.ServiceURL)
		err = reg.add(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = reg.remove(r.ServiceURL)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Removing Service %s with Url %s", r.ServiceName, r.ServiceURL)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
