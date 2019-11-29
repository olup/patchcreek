package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// TopicChannel is a channel of string
type TopicChannel chan string

// Subscriber is a subscriber
type Subscriber struct {
	c TopicChannel
}

// Topic is a slice of TopicChannel
type Topic []Subscriber

// Registry is the list of all topics
type Registry map[string]Topic

var myRegistry Registry

func readerToString(r io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	return buf.String()
}

func channelHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topic := vars["topic"]

	if r.Method == "GET" {
		sub := Subscriber{make(chan string)}

		if _, ok := myRegistry[topic]; ok {
			myRegistry[topic] = append(myRegistry[topic], sub)
		} else {
			myRegistry[topic] = make(Topic, 1)
			myRegistry[topic][0] = sub
		}

		response := <-sub.c

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, response)
	}

	if r.Method == "POST" {

		rBody := http.MaxBytesReader(w, r.Body, 100000)
		body := string(readerToString(rBody))

		SubcribersCount := 0

		if _, ok := myRegistry[topic]; ok {

			defer delete(myRegistry, topic)

			SubcribersCount = len(myRegistry[topic])

			for _, sub := range myRegistry[topic] {
				sub.c <- body
				defer close(sub.c)
			}
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{ "success" : true, servedSubscribers : `+strconv.Itoa(SubcribersCount)+" }")

	}
}

func main() {

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},                                               // All origins
		AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}, // Allowing only get, just an example
	})

	myRegistry = make(Registry)
	r := mux.NewRouter()
	r.HandleFunc("/{topic}", channelHandler)
	http.ListenAndServe(":5000", c.Handler(r))
}
