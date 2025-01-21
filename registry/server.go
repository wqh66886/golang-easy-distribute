package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

const (
	ServicePort = ":5678"                                        // registry service port
	ServiceURL  = "http://localhost" + ServicePort + "/services" // registry service url
)

var reg = registry{
	registrations: make([]Registration, 0),
}

type RegistryService struct{}

type registry struct {
	registrations []Registration
	mutex         sync.RWMutex
}

func (r *registry) add(reg Registration) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	r.registrations = append(r.registrations, reg)
	err := r.sendRequiredServices(reg) // send required services to all registered services
	if err != nil {
		return err
	}
	r.notify(patch{
		Added: []patchEntry{
			{
				Name: reg.ServiceName,
				URL:  reg.ServiceURL,
			},
		},
	})
	return nil
}

func (r *registry) remove(url string) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for i, reg := range r.registrations {
		if reg.ServiceURL == url {
			r.notify(patch{
				Removed: []patchEntry{
					{
						Name: reg.ServiceName,
						URL:  reg.ServiceURL,
					},
				},
			})
			r.registrations = append(r.registrations[:i], r.registrations[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("registry service not found url: %s\n", url)
}

// sendRequiredServices sends required services to all registered services
func (r *registry) sendRequiredServices(registration Registration) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	var p patch
	for _, serviceReg := range r.registrations { // iterate over all registered services
		for _, reqService := range registration.RequiredServices { // iterate over all required services
			if reqService == serviceReg.ServiceName { // if required service is found
				p.Added = append(p.Added, patchEntry{ // add to patch
					Name: serviceReg.ServiceName,
					URL:  serviceReg.ServiceURL,
				})
			}
		}
	}
	err := r.sendPatch(p, registration.ServiceUpdateURL)
	if err != nil {
		return err
	}
	return nil
}

// sendPatch sends patch to service
func (r *registry) sendPatch(p patch, url string) error {
	d, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return err
	}
	return nil
}

func (r *registry) notify(fullPatch patch) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, reg := range r.registrations {
		go func(reg Registration) {
			for _, reqService := range reg.RequiredServices {
				p := patch{Added: []patchEntry{}, Removed: []patchEntry{}}
				sendUpdate := false
				for _, add := range fullPatch.Added {
					if add.Name == reqService {
						p.Added = append(p.Added, add)
						sendUpdate = true
					}
				}
				for _, remove := range fullPatch.Removed {
					if remove.Name == reqService {
						p.Removed = append(p.Removed, remove)
						sendUpdate = true
					}
				}
				if sendUpdate {
					err := r.sendPatch(p, reg.ServiceUpdateURL)
					if err != nil {
						log.Println(err)
						return
					}
				}
			}
		}(reg)
	}
}

// RegistryService is a http handler, which handles POST and DELETE requests
func (rs RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Printf("response body closed failed: %v\n", err)
			}
		}()
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("adding service %s with url %s", r.ServiceName, r.ServiceURL)
		err = reg.add(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		dec := json.NewDecoder(r.Body)
		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Printf("response body closed failed: %v\n", err)
			}
		}()
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
		log.Printf("removing service %s with url %s", r.ServiceName, r.ServiceURL)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
