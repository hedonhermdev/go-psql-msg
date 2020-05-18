package main

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/hedonhermdev/go-psql-msg/config"
	"github.com/hedonhermdev/go-psql-msg/db"
)

func main() {

	conf_file, err := os.Open("config.yaml")
	if err != nil {
		log.Fatal("Could not open config file")
	}
	conf, err := config.ParseConfig(conf_file)
	if err != nil {
		log.Fatal("Could not parse config from config.yaml")
	}

	events := make(chan db.Event)
	errors := make(chan error)

	if len(conf.Channels) == 0 {
		log.Fatal("No channels to listen on. Exiting...")
	}

	for _, chName := range conf.Channels {
		fmt.Printf("Listening on channel %s...\n", chName)
		l := db.NewJSONListener(conf.Database)
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
