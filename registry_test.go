package registry

import (
	fmt "fmt"
	"log"
	"net"
	"testing"

	consul "github.com/hashicorp/consul/api"
	"github.com/rs/xid"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

type service struct {
	id        string
	tags      []string
	port      int
	address   string
	name      string
	innerName string
}

var (
	defaultName = "registry_test"

	services = []service{
		{
			id:        xid.New().String(),
			tags:      []string{"test1"},
			port:      4445,
			address:   "localhost",
			name:      defaultName,
			innerName: "test1-1-1",
		},
		{
			id:        xid.New().String(),
			tags:      []string{"test2"},
			port:      4446,
			address:   "localhost",
			name:      defaultName,
			innerName: "test2-1-2",
		},
		{
			id:        xid.New().String(),
			tags:      []string{"test3"},
			port:      4447,
			address:   "localhost",
			name:      defaultName,
			innerName: "test3-1-3",
		},
		{
			id:        xid.New().String(),
			tags:      []string{"test1", "test2"},
			port:      4448,
			address:   "localhost",
			name:      defaultName,
			innerName: "test1-2-12",
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

	var errChan = make(chan error)
	for _, service := range services {
		go func(errChan chan error) {
			err := startRPC(fmt.Sprintf("%s:%d", service.address, service.port), service.innerName)
			if err != nil {
				errChan <- err
			}
		}(errChan)
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
		select {
		case err := <-errChan:
			if err != nil {
				t.Error(err)
			}
		default:
		}
	}

}

type l struct {
}

func (l) Printf(s string, args ...interface{}) {
	log.Printf(s, args...)
}

func TestRegistryOn(t *testing.T) {
	//register(t)

	conf := consul.DefaultConfig()
	conf.Address = "localhost:8500"
	registry, err := New(conf, l{})
	if err == ErrConsulNotAvailable {
		t.Log(err)
		return
	}
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

func startRPC(bindAddress string, name string) error {
	lis, err := net.Listen("tcp", bindAddress)
	if err != nil {
		return err
	}

	imp := grpc.NewServer()
	s := &testServer{
		name: name,
	}

	log.Printf("Service listening on %s", bindAddress)
	RegisterHandlerServer(imp, s)
	if err := imp.Serve(lis); err != nil {
		return err
	}

	return nil
}

type testServer struct {
	name string
}

func (t *testServer) Process(context.Context, *Request) (*Response, error) {
	return &Response{
		NumOf: 10,
		Name:  t.name,
	}, nil
}
