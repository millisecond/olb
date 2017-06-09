package consul

import (
	"errors"
	"log"

	"github.com/fabiolb/fabio/config"
	"github.com/fabiolb/fabio/registry"

	"github.com/hashicorp/consul/api"
)

// be is an implementation of a registry backend for consul.
type be struct {
	c     *api.Client
	dc    string
	cfg   *config.Consul
	dereg chan bool
}

func NewBackend(cfg *config.Consul) (registry.Backend, error) {
	// create a reusable client
	c, err := api.NewClient(&api.Config{Address: cfg.Addr, Scheme: cfg.Scheme, Token: cfg.Token})
	if err != nil {
		return nil, err
	}

	// ping the agent
	dc, err := datacenter(c)
	if err != nil {
		return nil, err
	}

	// we're good
	log.Printf("[INFO] consul: Connecting to %q in datacenter %q", cfg.Addr, dc)
	return &be{c: c, dc: dc, cfg: cfg}, nil
}

func (b *be) Register() error {
	if !b.cfg.Register {
		log.Printf("[INFO] consul: Not registering fabio in consul")
		return nil
	}

	service, err := serviceRegistration(b.cfg)
	if err != nil {
		return err
	}

	b.dereg = register(b.c, service)
	return nil
}

func (b *be) Deregister() error {
	if b.dereg != nil {
		b.dereg <- true // trigger deregistration
		<-b.dereg       // wait for completion
	}
	return nil
}

func (b *be) ReadManual() (value string, version uint64, err error) {
	// we cannot rely on the value provided by WatchManual() since
	// someone has to call that method first to kick off the go routine.
	return getKV(b.c, b.cfg.KVPath, 0)
}

func (b *be) WriteManual(value string, version uint64) (ok bool, err error) {
	// try to create the key first by using version 0
	if ok, err = putKV(b.c, b.cfg.KVPath, value, 0); ok {
		return
	}

	// then try the CAS update
	return putKV(b.c, b.cfg.KVPath, value, version)
}

func (b *be) WatchServices() chan string {
	log.Printf("[INFO] consul: Using dynamic routes")
	log.Printf("[INFO] consul: Using tag prefix %q", b.cfg.TagPrefix)

	svc := make(chan string)
	go watchServices(b.c, b.cfg.TagPrefix, b.cfg.ServiceStatus, svc)
	return svc
}

func (b *be) WatchManual() chan string {
	log.Printf("[INFO] consul: Watching KV path %q", b.cfg.KVPath)

	kv := make(chan string)
	go watchKV(b.c, b.cfg.KVPath, kv)
	return kv
}

// datacenter returns the datacenter of the local agent
func datacenter(c *api.Client) (string, error) {
	self, err := c.Agent().Self()
	if err != nil {
		return "", err
	}

	cfg, ok := self["Config"]
	if !ok {
		return "", errors.New("consul: self.Config not found")
	}
	dc, ok := cfg["Datacenter"].(string)
	if !ok {
		return "", errors.New("consul: self.Datacenter not found")
	}
	return dc, nil
}
