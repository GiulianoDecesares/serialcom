package serialcom

import (
	"bufio"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/tarm/serial"
)

type Device struct {
	port     *serial.Port
	receiver *bufio.Scanner

	responseTimeOut time.Duration
}

func NewDevice(usbPort string, baud int, responseTimeOut time.Duration) (*Device, error) {
	port, err := serial.OpenPort(&serial.Config{Name: usbPort, Baud: baud})

	if err != nil {
		return nil, errors.Wrap(err, "error opening USB port")
	}

	// Arduino will auto reset itself after we open the port, so wait for a couple of seconds
	time.Sleep(2 * time.Second)

	return &Device{
		port:            port,
		receiver:        bufio.NewScanner(port),
		responseTimeOut: responseTimeOut,
	}, nil
}

func (device *Device) Release() error {
	return device.port.Close()
}

func (device *Device) Send(message Message) (Message, error) {
	bytes, err := json.Marshal(message)

	if err != nil {
		return Message{}, errors.Wrap(err, "error marshalling message")
	}

	if _, err = device.port.Write(bytes); err != nil {
		return Message{}, errors.Wrap(err, "error writing message to port")
	}

	// Wait for response with timeout
	timeout := time.After(device.responseTimeOut)
	responseChan := make(chan Message, 1)
	errChan := make(chan error, 1)

	go func() {
		if device.receiver.Scan() {
			data := device.receiver.Bytes()
			message := Message{}

			err = json.Unmarshal(data, &message)

			if err != nil {
				errChan <- errors.Wrap(err, "error unmarshalling message")
				return
			}

			responseChan <- message
		}
	}()

	select {
	case <-timeout:
		return Message{}, errors.New("device response timed out")

	case currentErr := <-errChan:
		return Message{}, currentErr

	case currentResponse := <-responseChan:
		return currentResponse, nil
	}
}
