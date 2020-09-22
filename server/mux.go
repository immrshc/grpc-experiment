package server

import (
	"reflect"
	"strings"
)

type Server interface {
	Start() error
	AsyncErr() <-chan error
}

type mux struct {
	servers []Server
}

type multiErrors struct {
	errors []error
}

func (m *multiErrors) Error() string {
	messages := make([]string, 0)
	for _, err := range m.errors {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, ", ")
}

func NewMux(servers []Server) *mux {
	return &mux{servers: servers}
}

func (m *mux) Serve() error {
	errs := make([]error, len(m.servers))
	cases := make([]reflect.SelectCase, len(m.servers))
	for i, s := range m.servers {
		if err := s.Start(); err != nil {
			errs[i] = err
			break
		}
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(s.AsyncErr())}
	}
	chosen, value, ok := reflect.Select(cases)
	if chosen < len(cases) && ok {
		if err, ok := value.Interface().(error); ok {
			errs[chosen] = err
		}
	}
	return &multiErrors{
		errors: errs,
	}
}
