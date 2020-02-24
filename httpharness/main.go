package main

// HTTP Test Harness example.  Lifted from the Golang docs:
// https://golang.org/pkg/net/http/#example_ListenAndServe
//
// Original Author: Jack Neely <jjneely@42lines.net>
// 2020/02/24

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Set a session cookie that is just the current time
func setCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:    "HSESSIONID",
		Value:   time.Now().String(),
		Expires: time.Now().AddDate(0, 0, 1),
	}
	http.SetCookie(w, &cookie)
}

func main() {
	running := true
	world := flag.String("name", "world",
		"The name to greet")
	bind := flag.String("bind", ":8080",
		"The HOST:PORT combination to listen on")
	flag.Parse()

	// Define some basic content that can identify itself via --name
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		if _, err := req.Cookie("HSESSIONID"); err == http.ErrNoCookie {
			setCookie(w)
		} else if err != nil {
			http.Error(w, "Bad HSESSIONID cookie", http.StatusBadRequest)
		} else {
			fmt.Fprintf(w, "Hello, %s!\n", *world)
		}
	}

	// Report if this process in "healthy" which it should be unless
	// shutdown is called
	healthHandler := func(w http.ResponseWriter, req *http.Request) {
		if !running {
			http.Error(w, "Healthy: NO!", http.StatusServiceUnavailable)
		} else {
			fmt.Fprintf(w, "Healthy: %v\n", running)
		}
	}

	// Handle requests to gracefully shutdown
	shutdownHandler := func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "running: %v!\n", running)
		running = false
	}

	// All undefined paths in our tree are, well, not found
	notFoundHandler := func(w http.ResponseWriter, req *http.Request) {
		s := fmt.Sprintf("404: %s", req.URL.String())
		http.Error(w, s, http.StatusNotFound)
	}

	http.HandleFunc("/", notFoundHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/rest/healthCheck", healthHandler)
	http.HandleFunc("/rest/system/shutdown", shutdownHandler)
	log.Printf("Running webserver for \"%s\" on %s", *world, *bind)
	log.Fatal(http.ListenAndServe(*bind, nil))
}
