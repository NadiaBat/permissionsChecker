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
type AdditionalParams struct{
	region  int
	project int
}

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
	server := http.Server{Addr: ":8888", Handler: handler}

	defer s.Shutdown(server)

	server.ListenAndServe()
}

func (s *Server) Shutdown(server http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Server shutdown error: %s", err)
	}
}

func handlerCheck(r *http.Request) ([]byte, error) {
	userId, actions, additionalParams, err := getParams(r)
	if err != nil {
		return nil, errors.Wrap(err, "Can`t get AdditionalParams.")
	}

	permissions, err := BulkCheck(userId, actions, additionalParams)
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

func getParams(r *http.Request) (int, []string, AdditionalParams, error) {
	queryParams := r.URL.Query()
	var additionalParams AdditionalParams

	actions := queryParams["actions[]"]
	if actions == nil || len(actions) == 0 {
		return 0, nil, additionalParams, errors.New("Обязательный параметр actions должен быть не пустым массивом.")
	}

	userFromQuery := queryParams["userId"]
	if userFromQuery == nil {
		return 0, actions, additionalParams, errors.New("Не задан обязательный параметр userId.")
	}

	userId, err := strconv.Atoi(userFromQuery[0])
	if err != nil {
		return 0, actions, additionalParams, errors.Wrap(err, "Can`t get userId.")
	}

	if userId == 0 {
		return 0, actions, additionalParams, errors.New("User can`t be 0.")
	}

	regionId, ok := queryParams["params[region]"]
	if ok {
		additionalParams.region, err = strconv.Atoi(regionId[0])
		if err != nil {
			return userId, actions, additionalParams, errors.Wrapf(err, "Region must be integer.")
		}
	}

	projectId, ok := queryParams["params[project]"]
	if ok {
		additionalParams.project, err = strconv.Atoi(projectId[0])
		if err != nil {
			return userId, actions, additionalParams, errors.Wrapf(err, "Region must be integer.")
		}
	}

	return userId, actions, additionalParams, nil
}
