package listener

import (
	"encoding/json"

	"github.com/lib/pq"
)

const (
	INSERT = "insert"
	UPDATE = "update"
	DELETE = "delete"
)

type field map[string]interface{}

type Event struct {
	Channel   string
	EventType string `json:"action"`
	TableName string `table:"table"`
	Data      field  `json:"data"`
}

func ParseEvent(n *pq.Notification) (*Event, error) {
	e := new(Event)
	e.Channel = n.Channel
	err := json.Unmarshal([]byte(n.Extra), e)
	if err != nil {
		return nil, err
	}
	return e, nil
}
