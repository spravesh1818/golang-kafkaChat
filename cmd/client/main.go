package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	chatapp "golang_kafka_chat_application/pkg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	medium "golang_kafka_chat_application/pkg/medium"

	_ "github.com/joho/godotenv/autoload"
)

var user string

type Request struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func printMessage(m chatapp.Message) {
	switch {
	case m.Username == user:
		return
	case m.Username == chatapp.SystemID:
		fmt.Println(m.Message)
	default:
		fmt.Printf("%s: %s", m.Username, m.Message)
	}
}

func processRequest(endpoint string, data Request) error {
	requestBody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	res, err := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusAccepted {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}

	return nil
}

func joinRoom(host string) error {
	d := Request{Username: user}
	endpoint := fmt.Sprintf("%s/join", host)
	return processRequest(endpoint, d)
}

func publishMessage(host string, message string) error {
	d := Request{Username: user, Message: message}
	endpoint := fmt.Sprintf("%s/publish", host)
	return processRequest(endpoint, d)
}

func main() {
	var (
		host    = os.Getenv("HOST")
		brokers = os.Getenv("KAFKA_BROKERS")
		topic   = os.Getenv("KAFKA_TOPIC")
	)

	fmt.Println("What do you want to be called while using the app?")
	_, _ = fmt.Scanf("%s\n", &user)

	if err := joinRoom(host); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Joined Room Successfully")
	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)

			// Print the prompt message to the console

			msg, _ := reader.ReadString('\n')
			if err := publishMessage(host, msg); err != nil {
				log.Fatal(err)
			}
		}
	}()

	chMsg := make(chan chatapp.Message)
	chErr := make(chan error)
	consumer := medium.NewConsumer(strings.Split(brokers, ","), topic)

	// Create a buffered channel with a buffer size of 1 for handling signals
	quit := make(chan os.Signal, 1)

	// Notify the quit channel for SIGINT and SIGTERM signals
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Your goroutine for consuming messages
	go func() {
		consumer.Read(context.Background(), chMsg, chErr)
	}()

	// Wait for a signal
	<-quit

	for {
		select {
		case <-quit:
			goto end
		case m := <-chMsg:
			printMessage(m)
		case err := <-chErr:
			log.Println(err)
		}
	}
end:

	fmt.Println("\nyou have abandoned the room")
}
