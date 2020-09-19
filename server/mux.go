package server

import "strings"

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
	for i, s := range m.servers {
		if err := s.Start(); err != nil {
			return err
		}
		// ここでループが止まらないか？
		errs[i] = <-s.AsyncErr()
	}
	return &multiErrors{
		errors: errs,
	}
}
