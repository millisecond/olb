package consul

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fabiolb/fabio/config"
	"github.com/hashicorp/consul/api"
)

// register keeps a service registered in consul.
//
// When a value is sent in the dereg channel the service is deregistered from
// consul. To wait for completion the caller should read the next value from
// the dereg channel.
//
//    dereg <- true // trigger deregistration
//    <-dereg       // wait for completion
//
func register(c *api.Client, service *api.AgentServiceRegistration) (dereg chan bool) {
	var serviceID string

	registered := func() bool {
		if serviceID == "" {
			return false
		}
		services, err := c.Agent().Services()
		if err != nil {
			log.Printf("[ERROR] consul: Cannot get service list. %s", err)
			return false
		}
		return services[serviceID] != nil
	}

	register := func() {
		if err := c.Agent().ServiceRegister(service); err != nil {
			log.Printf("[ERROR] consul: Cannot register fabio in consul. %s", err)
			return
		}

		log.Printf("[INFO] consul: Registered fabio with id %q", service.ID)
		log.Printf("[INFO] consul: Registered fabio with address %q", service.Address)
		log.Printf("[INFO] consul: Registered fabio with tags %q", strings.Join(service.Tags, ","))
		log.Printf("[INFO] consul: Registered fabio with health check to %q", service.Check.HTTP)

		serviceID = service.ID
	}

	deregister := func() {
		log.Printf("[INFO] consul: Deregistering fabio")
		c.Agent().ServiceDeregister(serviceID)
	}

	dereg = make(chan bool)
	go func() {
		register()
		for {
			select {
			case <-dereg:
				deregister()
				dereg <- true
				return
			case <-time.After(time.Second):
				if !registered() {
					register()
				}
			}
		}
	}()
	return dereg
}

func serviceRegistration(cfg *config.Consul) (*api.AgentServiceRegistration, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	ipstr, portstr, err := net.SplitHostPort(cfg.ServiceAddr)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(portstr)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(ipstr)
	if ip == nil {
		ip, err = config.LocalIP()
		if err != nil {
			return nil, err
		}
		if ip == nil {
			return nil, errors.New("no local ip")
		}
	}

	serviceID := fmt.Sprintf("%s-%s-%d", cfg.ServiceName, hostname, port)

	checkURL := fmt.Sprintf("%s://%s:%d/health", cfg.CheckScheme, ip, port)
	if ip.To16() != nil {
		checkURL = fmt.Sprintf("%s://[%s]:%d/health", cfg.CheckScheme, ip, port)
	}

	service := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    cfg.ServiceName,
		Address: ip.String(),
		Port:    port,
		Tags:    cfg.ServiceTags,
		Check: &api.AgentServiceCheck{
			HTTP:          checkURL,
			Interval:      cfg.CheckInterval.String(),
			Timeout:       cfg.CheckTimeout.String(),
			TLSSkipVerify: cfg.CheckTLSSkipVerify,
		},
	}

	return service, nil
}
