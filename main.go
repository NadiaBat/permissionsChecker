package main

import (
	//"github.com/NadiaBat/permissionsChecker/http"
	"github.com/NadiaBat/permissionsChecker/storage"
)

type message struct {
	isActive chan int
}

type user struct {
	id  int
	fio string
}

func main() {
	assignments := storage.GetAllAssignments(true)
	userAssignments := assignments[200132743]
	println(userAssignments.UserId)

	for _, assignment := range userAssignments.Items {
		println(assignment.ItemName)
	}
	//
	//h := http.Server{}
	//h.Serve()
}

// bulk checker by user and operations

// cli
// --host
// --port
// --dbDsn
// docopt.go (for cli generations) @see git.rn/projects/PORTAL/repos/paged/browse

// for http use standart golang http (/net/http) @see https://cryptic.io/go-http/

// use channels for updating data from db

// handlers for bulk operations

// use Notify for hardware interruptions
// use logger instead of panic, panic as an exclusion

// teeeestsss!!!!!!!!!!!!!!! don`t run? only tests
