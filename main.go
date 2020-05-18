package main

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/hedonhermdev/go-psql-msg/db"
	"github.com/hedonhermdev/go-psql-msg/db/listener"
)

func main() {

	conn_info := db.LoadFromEnv()
	conn_string := conn_info.ConnString()

	if len(os.Args) == 1 {
		log.Fatal("No channel names provided to listen on")
	}

	events := make(chan listener.Event)
	errors := make(chan error)

	for _, chName := range os.Args[1:] {
		fmt.Printf("Listening on channel %s...\n", chName)
		l := listener.NewJSONListener(conn_string)
		go l.Listen(chName, events, errors)
	}

	for {
		select {
		case err := <-errors:
			log.Print(err)
		case event := <-events:
			fmt.Printf("%+v\n", event)
			for k, v := range event.Data {
				fmt.Printf("%s: %v\n", k, reflect.TypeOf(v))
			}
		}
	}

}
