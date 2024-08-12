package serialcom

func IsOK(message Message) bool {
	return message.Id == "ok"
}

func Ping() Message {
	return Message{
		Id:   "ping",
		Data: make(map[string]string),
	}
}

func Pong() Message {
	return Message{
		Id:   "pong",
		Data: make(map[string]string),
	}
}
