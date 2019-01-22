package registry

import (
	"log"
	"testing"

	consul "github.com/hashicorp/consul/api"
	"github.com/rs/xid"
)

type service struct {
	id      string
	tags    []string
	port    int
	address string
	name    string
}

var (
	defaultName = "registry_test"

	services = []service{
		{
			id:      xid.New().String(),
			tags:    []string{"test1"},
			port:    4445,
			address: "192.168.1.1",
			name:    defaultName,
		},
		{
			id:      xid.New().String(),
			tags:    []string{"test2"},
			port:    4445,
			address: "192.168.1.2",
			name:    defaultName,
		},
		{
			id:      xid.New().String(),
			tags:    []string{"test3"},
			port:    4445,
			address: "192.168.1.3",
			name:    defaultName,
		},
		{
			id:      xid.New().String(),
			tags:    []string{"test1", "test2"},
			port:    4445,
			address: "192.168.1.12",
			name:    defaultName,
		},
	}
)

func register(t *testing.T) {
	config := consul.DefaultConfig()
	config.Address = "localhost:8500"
	c, err := consul.NewClient(config)
	if err != nil {
		t.Error(err)
	}

	for _, service := range services {
		reg := &consul.AgentServiceRegistration{
			ID:      service.id,
			Name:    service.name,
			Tags:    service.tags,
			Port:    service.port,
			Address: service.address,
		}

		err := c.Agent().ServiceRegister(reg)
		if err != nil {
			t.Error(err)
		}
	}
}

type l struct {
}

func (l) Print(v ...interface{}) {
	log.Print(v...)
}

func TestRegistryOn(t *testing.T) {
	register(t)

	conf := consul.DefaultConfig()
	conf.Address = "localhost:8500"
	registry, err := New(conf, l{})
	if err != nil {
		t.Error(err)
		return
	}

	err = registry.PeriodicCheck(defaultName, []string{"test1", "test2", "test3"})
	if err != nil {
		t.Error(err)
		return
	}
}
