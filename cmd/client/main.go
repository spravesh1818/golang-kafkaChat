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

	"github.com/fatih/color"
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
		messageColor := color.New(color.FgGreen, color.Bold)
		messageColor.Println(m.Message)
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

	messageColor := color.New(color.FgCyan)

	messageColor.Println("What do you want to be called while using the app?")
	_, _ = fmt.Scanf("%s\n", &user)

	if err := joinRoom(host); err != nil {
		log.Fatal(err)
	}

	messageColor.Println("Joined Room Successfully.Type your message on the empty prompt to continue")
	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)

			msg, _ := reader.ReadString('\n')
			if err := publishMessage(host, msg); err != nil {
				log.Fatal(err)
			}
		}
	}()

	chMsg := make(chan chatapp.Message)
	chErr := make(chan error)
	consumer := medium.NewConsumer(strings.Split(brokers, ","), topic)

	go func() {
		consumer.Read(context.Background(), chMsg, chErr)
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

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
