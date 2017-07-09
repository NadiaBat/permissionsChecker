package httpServer

import (
	"net/http"
	"sync"
	"github.com/NadiaBat/permissionsChecker/models/permission"
	"strconv"
	"os"
)

type Http struct { }
type httpHandler struct {wg *sync.WaitGroup}

func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	switch r.URL.Path {
	case "/check/":
		routeCheck(r)
	default:
	}
}

func (h Http) ServeHttp(wg *sync.WaitGroup)  {
	handler := httpHandler{wg: wg}
	http.ListenAndServe(":9999", handler)
}

func routeCheck(r *http.Request)  {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	params := r.URL.Query()
	actions := getActionsOrFail(params)
	wg.Add(len(actions))

	userId := getUserIdOrFail(params)

	assignments := []*permission.Assignment{}
	for _, actionName := range actions {
		assignments = append(assignments, getAssignment(userId, actionName))
	}

	permission.BulkCheck(&wg, assignments)
}

func getUserIdOrFail(params map[string][]string) (int) {
	userId, err := strconv.Atoi(params["userId"][0])
	if (err != nil) {
		os.Exit(2)
	}

	return userId
}

func getActionsOrFail(params map[string][]string) ([]string) {
	actions := params["actions[]"]
	if (actions == nil || len(actions) == 0) {
		os.Exit(2)
	}

	return actions
}

func getAssignment(userId int, actionName string) (*permission.Assignment) {
	action := permission.Action{Name: actionName}
	assignment := permission.Assignment{UserId: userId, Action: action}
	assignment.HasAccess = true

	return &assignment
}
