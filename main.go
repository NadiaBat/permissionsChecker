package main

import (
	"github.com/NadiaBat/permissionsChecker/http"
)

type message struct {
	isActive chan int
}

func main() {
	//var wg sync.WaitGroup
	//wg.Add(1)
	//defer wg.Wait()
	h := http.Server{}
	h.Serve()
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
