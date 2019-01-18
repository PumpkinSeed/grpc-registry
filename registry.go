package registry

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	consul "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type logger interface {
	Print(v ...interface{})
}

type Registry struct {
	conns  map[string]*grpc.ClientConn
	consul *consul.Client
	block  bool
	log    logger

	creds credentials.TransportCredentials
	mut   *sync.Mutex
}

type Config struct {
	Address string
}

func New(c *consul.Config, log logger) (*Registry, error) {
	cc, err := consul.NewClient(c)
	if err != nil {
		return nil, err
	}

	return &Registry{
		consul: cc,
		conns:  make(map[string]*grpc.ClientConn),
		block:  false,
		mut:    new(sync.Mutex),
		log:    log,
	}, nil
}

func (r *Registry) SetCredentials(creds credentials.TransportCredentials) {
	r.creds = creds
}

/*func (r *Registry) Add(tag string, id string, c *grpc.ClientConn) error {
	if _, ok := r.conns[tag][id]; ok && r.conns[tag][id].health {
		return errors.New("client exists and healthy")
	}

	r.conns[tag][id] = client{
		client: c,
		health: true,
	}

	return nil
}*/

func (r *Registry) Get(tag string) (*grpc.ClientConn, error) {
	for r.block {
		fmt.Println("wait")
		time.Sleep(10 * time.Millisecond)
	}

	if _, ok := r.conns[tag]; !ok {
		return nil, errors.New("invalid tag")
	}

	return r.conns[tag], nil
}

func (r *Registry) PeriodicCheck(name string, tags []string) error {
	var services = make(map[string]*consul.ServiceEntry)
	for _, tag := range tags {
		s, _, err := r.consul.Health().Service(name, tag, true, nil)
		if err != nil {
			return err
		}
		if len(s) > 0 {
			services[tag] = s[0]
		}

	}

	r.block = true
	defer func() {
		r.block = false
	}()
	var tagsLog []string
	for tag, service := range services {
		target := fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port)
		conn, err := grpc.Dial(target, grpc.WithTransportCredentials(r.creds))
		if err != nil {
			return err
		}

		r.mut.Lock()
		r.conns[tag] = conn
		r.mut.Unlock()
		tagsLog = append(tagsLog, tag)
	}

	r.log.Print("PeriodicCheck refresh the registry, tags currently: ", strings.Join(tagsLog, ", "))
	return nil
}
