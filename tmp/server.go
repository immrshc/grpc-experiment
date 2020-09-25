package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type muxErrs struct {
	errs []error
}

func (m *muxErrs) Error() string {
	messages := make([]string, 0)
	for _, err := range m.errs {
		if err != nil {
			messages = append(messages, err.Error())
		}
	}
	return strings.Join(messages, ", ")
}

type server interface {
	Start() error
	ErrChan() <-chan error
}

type mux struct {
	servers []server
}

func (m *mux) Serve() error {
	errs := make([]error, len(m.servers))
	cases := make([]reflect.SelectCase, len(m.servers))
	for i, s := range m.servers {
		if err := s.Start(); err != nil {
			errs[i] = err
			break
		}
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(s.ErrChan())}
	}
	chosen, value, ok := reflect.Select(cases)
	if ok {
		if err, ok := value.Interface().(error); ok {
			errs[chosen] = err
		}
	}
	return &muxErrs{
		errs: errs,
	}
}

type webServer struct {
	Addr         string
	SurvivalTime time.Duration
	errCh        chan error
}

func (w *webServer) Start() error {
	mux := http.NewServeMux()
	server := &http.Server{Addr: w.Addr, Handler: mux}

	w.errCh = make(chan error, 1)
	go func() {
		fmt.Printf("----ListenAndServe %s----\n", w.Addr)
		w.errCh <- server.ListenAndServe()
		close(w.errCh)
	}()
	go func() {
		time.Sleep(w.SurvivalTime)
		fmt.Printf("-------Shutdown %s-------\n", w.Addr)
		if err := server.Shutdown(context.Background()); err != nil {
			w.errCh <- err
			close(w.errCh)
		}
	}()
	return nil
}

func (w *webServer) ErrChan() <-chan error {
	return w.errCh
}

func main() {
	s1 := &webServer{Addr: ":8080", SurvivalTime: 10 * time.Second}
	s2 := &webServer{Addr: ":8081", SurvivalTime: 5 * time.Second}
	m := mux{
		servers: []server{s1, s2},
	}
	if err := m.Serve(); err != nil {
		log.Printf("main exited because %s", err)
	}
}
