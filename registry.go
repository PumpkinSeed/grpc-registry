package registry

import (
	"errors"

	consul "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

type Registry struct {
	conns map[string]map[string]client

	consul *consul.Client
}

type client struct {
	client *grpc.ClientConn
	health bool
}

type Config struct {
	Address string
}

func New(c *consul.Config) (*Registry, error) {
	cc, err := consul.NewClient(c)
	if err != nil {
		return nil, err
	}

	return &Registry{
		consul: cc,
		conns:  make(map[string]map[string]client),
	}, nil
}

func (r *Registry) Add(tag string, id string, c *grpc.ClientConn) error {
	if _, ok := r.conns[tag][id]; ok && r.conns[tag][id].health {
		return errors.New("client exists and healthy")
	}

	r.conns[tag][id] = client{
		client: c,
		health: true,
	}

	return nil
}

func (r *Registry) Get(tag string) (*grpc.ClientConn, error) {
	if _, ok := r.conns[tag]; !ok {
		return nil, errors.New("invalid tag")
	}
	for _, conn := range r.conns[tag] {
		if conn.health {
			return conn.client, nil
		}
	}

	return nil, errors.New("no connection found")
}

func (r *Registry) PeriodicCheck(name string, tags []string) error {
	services, _, err := r.consul.Health().ServiceMultipleTags(name, tags, true, nil)
	if err != nil {
		return err
	}

	for _, service := range services {

	}
	return nil
}
