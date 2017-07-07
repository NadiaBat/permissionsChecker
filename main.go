package main

import (
"fmt"
"time"
"sync"
"github.com/NadiaBat/permissionsChecker/httpServer"
)

type message struct {
	isActive chan int
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	h := httpServer.Http{}
	h.ServeHttp(&wg)
}






func watersExamples()  {
	channel := message{}
	var wg sync.WaitGroup
	wg.Add(2)

	channel.isActive = make(chan int)

	go func() {
		defer wg.Done()
		time.Sleep(time.Second)
		channel.isActive <- 1200
	}()

	go func() {
		defer wg.Done()
		fmt.Println(<- channel.isActive)
	}()

	wg.Wait()
	fmt.Println("-----------")
}




//func getMessage(i int) (message) {
//	currentTime := time.Now().Format("Y-m-d H-i-s")
//
//	header := fmt.Sprintf("Header %g by %s", i, currentTime)
//	text := fmt.Sprintf("Text %g by %s", i, currentTime)
//
//	return message{header, text, <- i}
//}



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
