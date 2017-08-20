package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type Server struct{}
type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/check/":
		w.Header().Set("Content-Type", "application/json")
		result, err := handlerCheck(r)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		w.Write(result)
	default:
	}
}

func (s *Server) Serve() {
	handler := handler{}
	server := http.Server{Addr: ":9999", Handler: handler}

	defer s.Shutdown(server)

	server.ListenAndServe()
}

func (s *Server) Shutdown(server http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err == nil {
		log.Fatalf("Server shutdown error: %s", err)
	}
}

func handlerCheck(r *http.Request) ([]byte, error) {
	userId, actions, err := getParams(r)
	if err != nil {
		return nil, errors.Wrap(err, "Can`t get params.")
	}

	permissions, err := BulkCheck(userId, actions, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Can`t execute bulk checking.")
	}

	result, err := json.Marshal(permissions)
	if err != nil {
		return result, errors.Wrap(err, "Permissions object "+
			"to json marshal failed")
	}

	return result, err
}

func getParams(r *http.Request) (int, []string, error) {
	params := r.URL.Query()

	actions := params["actions[]"]
	if actions == nil || len(actions) == 0 {
		return 0, nil, errors.New("Обязательный параметр actions должен быть не пустым массивом.")
	}

	userId, err := strconv.Atoi(params["userId"][0])
	if err != nil {
		return 0, actions, errors.Wrap(err, "Can`t get userId.")
	}

	if userId == 0 {
		return 0, actions, errors.New("User can`t be 0.")
	}

	return userId, actions, nil
}
