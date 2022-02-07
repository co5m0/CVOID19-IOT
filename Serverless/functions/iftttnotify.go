package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/nuclio/nuclio-sdk-go"
)

func Handler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
	context.Logger.Info("This is an unstrucured %s", "log")

	ENTER_URL := "https://cvoid19-id.firebaseio.com/cvoid-bar/door-enter.json"
	EXIT_URL := "https://cvoid19-id.firebaseio.com/cvoid-bar/door-exit.json"
	IFTTT_URL := "https://maker.ifttt.com/trigger/cvoid_notify/with/key/xxxxxxxxx"
	enter, err := countPeople(ENTER_URL, "enter")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(enter)
	exit, err := countPeople(EXIT_URL, "exit")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(exit)
	limit, err := strconv.Atoi(os.Getenv("LIMIT"))
	if err != nil {
		log.Fatalln("ENV NOT FOUND")
	}
	peopleInside := enter - exit
	if peopleInside >= limit {
		_, err := http.Get(IFTTT_URL)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func countPeople(url string, match string) (int, error) {
	people := 0
	var errorR error
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
		errorR = err
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		errorR = err
	}
	//Convert the body to type string
	sb := string(body)
	people = strings.Count(sb, match)
	return people, errorR
}