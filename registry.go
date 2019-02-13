package registry

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	consul "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var ErrConsulNotAvailable = errors.New("consul not available")

const heatlhCheckTimeout = 500

type logger interface {
	Printf(string, ...interface{})
}

type Registry struct {
	conns        map[string]*grpc.ClientConn
	consul       *consul.Client
	consulConfig *consul.Config
	block        bool
	log          logger
	available    bool

	creds credentials.TransportCredentials
	mut   *sync.Mutex
}

type Config struct {
	Address string
}

func New(c *consul.Config, log logger) (*Registry, error) {
	reg := &Registry{
		consulConfig: c,
		conns:        make(map[string]*grpc.ClientConn),
		block:        false,
		mut:          new(sync.Mutex),
		log:          log,
	}
	if !reg.HealthCheck() {
		reg.available = false
		return reg, ErrConsulNotAvailable
	}

	cc, err := consul.NewClient(c)
	if err != nil {
		return nil, err
	}
	reg.consul = cc

	return reg, nil
}

func (r *Registry) Available() []string {
	var result []string
	if !r.available {
		return result
	}

	for tag := range r.conns {
		result = append(result, tag)
	}
	return result
}

func (r *Registry) SetCredentials(creds credentials.TransportCredentials) {
	r.creds = creds
}

func (r *Registry) Get(tag string) (*grpc.ClientConn, error) {
	for r.block {
		time.Sleep(10 * time.Millisecond)
	}

	if value, ok := r.conns[tag]; !ok || value == nil {
		return nil, errors.New("invalid tag, requested service not available")
	}

	return r.conns[tag], nil
}

func (r *Registry) HealthCheck() bool {
	conn, err := net.DialTimeout("tcp", r.consulConfig.Address, time.Duration(heatlhCheckTimeout)*time.Millisecond)
	if err != nil || conn == nil {
		return false
	}
	defer conn.Close()
	return true
}

func (r *Registry) PeriodicCheck(name string, tags []string) error {
	if !r.HealthCheck() {
		r.available = false
		return ErrConsulNotAvailable
	}

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
		var conn *grpc.ClientConn
		var err error
		if r.creds == nil {
			conn, err = grpc.Dial(target, grpc.WithInsecure())
		} else {
			conn, err = grpc.Dial(target, grpc.WithTransportCredentials(r.creds))
		}
		if err != nil {
			return err
		}

		r.mut.Lock()
		r.conns[tag] = conn
		r.mut.Unlock()
		tagsLog = append(tagsLog, tag)
	}

	if r.log != nil {
		r.log.Printf("PeriodicCheck refresh the registry, tags currently: ", strings.Join(tagsLog, ", "))
	}
	return nil
}

func (r *Registry) Close() error {
	for _, conn := range r.conns {
		if conn != nil {
			if err := conn.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}
