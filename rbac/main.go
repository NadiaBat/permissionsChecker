package main

import (
	"log"
	"github.com/pkg/errors"
	"strconv"
	"fmt"
)

func main()  {
	err := InitMysql()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Couldn`t init mysql connection."))
	}

	err = InitDicts()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Couldn`t init dictionaries."))
	}

	params := make(InputParams)
	//params["type"] = "video_of_the_day"
	userId := 217188183
	params["userId"] = strconv.Itoa(userId)

	checkingSet1 := CheckingSet{
		217188183,
		"2_1_maineditor",
		params,
		false,
	}

	params["type"] = "video_of_the_day"
	checkingSet2 := CheckingSet{
		217188183,
		"ncc.recordsoftheday.delete.access",
		params,
		false,
	}

	checkingSets := CheckingSets{&checkingSet1, &checkingSet2}
	err = BulkCheckAccess(checkingSets)

	fmt.Println(checkingSets[0], checkingSets[1])
}
