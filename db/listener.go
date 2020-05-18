package db

import (
	"fmt"
	"log"
	"time"

	"github.com/hedonhermdev/go-psql-msg/config"
	"github.com/lib/pq"
)

const (
	SSLMODE = "disable"
)

type JSONListener struct {
	listener *pq.Listener
}

func reportProblem(ev pq.ListenerEventType, err error) {
	if err != nil {
		log.Printf("ERROR: %q", err)
	}
}

func makeConnString(conf config.DBConfig) string {
	conn_string := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s", conf.Host, conf.Port, conf.DBName, conf.User, conf.Password, SSLMODE)
	return conn_string
}

func NewJSONListener(conf config.DBConfig) *JSONListener {
	conn_string := makeConnString(conf)
	return &JSONListener{pq.NewListener(conn_string, 5*time.Second, 100*time.Second, reportProblem)}
}

func (l *JSONListener) Listen(pg_channel string, eventChan chan Event, errorChan chan error) error {
	err := l.listener.Listen(pg_channel)
	if err != nil {
		return err
	}
	for {
		select {
		case n := <-l.listener.Notify:
			event, err := ParseEvent(n)
			if err != nil {
				errorChan <- err
				continue
			}
			eventChan <- *event
		case <-time.After(90 * time.Second):
			log.Print("INFO: No new notifications in 90 seconds, checking connection")
			err := l.listener.Ping()
			if err != nil {
				return err
			}
		}
	}
}
