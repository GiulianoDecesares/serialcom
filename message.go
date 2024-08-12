package serialcom

import (
	"fmt"
)

type Message struct {
	Id   string            `json:"id"`
	Data map[string]string `json:"data"`
}

func (message Message) String() string {
	str := fmt.Sprintf("Id: %s", message.Id)

	for key, value := range message.Data {
		str += fmt.Sprintf("\n[%s -- %s]", key, value)
	}

	return str
}
