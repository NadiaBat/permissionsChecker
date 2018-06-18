package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type Server struct {
}
type handler struct {
	rbac Rbac
}

type paramsSet struct {
	UserID           int                `json:"user_id,string,omitempty"`
	Action           string             `json:"action,omitempty"`
	AdditionalParams []AdditionalParams `json:"params"`
}

type AdditionalParams struct {
	UserID       int  `json:"user_id,string,omitempty"`
	Region       int  `json:"regionId,string,omitempty"`
	Project      int  `json:"projectId,string,omitempty"`
	IsCommercial bool `json:"isCommercial,omitempty"`
	StringParams map[string]string
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/check/":
		w.Header().Set("Content-Type", "application/json")
		result, err := h.handlerCheck(r)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		w.Write(result)
	default:
	}
}

func (s *Server) Serve(rbac *Rbac) {
	handler := handler{rbac: *rbac}
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

func (h handler) handlerCheck(r *http.Request) ([]byte, error) {
	fmt.Println("HANDLER CHECK")
	sets, err := getParams(r)
	fmt.Println(sets)
	if err != nil {
		return nil, errors.Wrap(err, "Can`t get AdditionalParams.")
	}

	permissions := Permissions{}
	permissionsAll := PermissionsAll{}
	for _, params := range sets {
		for _, addParams := range params.AdditionalParams {
			// res := Permission{}
			res, err := h.rbac.Check(params.UserID, params.Action, addParams)
			fmt.Println(res)
			if err != nil {
				return nil, errors.Wrap(err, "Checking error.")
			}
			permissions = append(permissions, &res)
		}
		permissionsAll = append(permissionsAll, permissions)
		permissions = Permissions{}
	}

	result, err := json.Marshal(permissionsAll)
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
	fmt.Println(params)
	return params, nil
}
