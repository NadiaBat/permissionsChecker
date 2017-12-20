package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"time"
)

type Server struct{}
type handler struct{}

type paramsSet struct {
	UserId           int              `json:"userId"`
	Action           string           `json:"action"`
	AdditionalParams AdditionalParams `json:"params"`
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
	sets, err := getParams(r)

	if err != nil {
		return nil, errors.Wrap(err, "Can`t get AdditionalParams.")
	}

	permissions := Permissions{}
	for _, params := range sets {
		res, err := checkOne(params.UserId, params.Action, params.AdditionalParams)
		if err != nil {
			return nil, errors.Wrap(err, "Checking error.")
		}

		permissions = append(permissions, &res)
	}

	result, err := json.Marshal(permissions)
	if err != nil {
		return nil, errors.Wrap(err, "Response encoding error.")
	}

	return result, err
}

func getParams(r *http.Request) ([]paramsSet, error) {
	var params []paramsSet
	hash := md5.New()
	reader := io.TeeReader(r.Body, hash)
	err := json.NewDecoder(reader).Decode(&params)
	if err != nil {
		return params, errors.Wrap(err, "Request body JSON decoding error.")
	}

	return params, nil
}

func checkOne(userId int, action string, additionalParams AdditionalParams) (Permission, error) {
	permission, err := Check(userId, action, additionalParams)
	if err != nil {
		return Permission{}, errors.Wrap(err, "Can`t execute bulk checking.")
	}

	return permission, nil
}
