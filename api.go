package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//api runs the http server for api endpoints
func api() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No handler"))
		return
	})

	http.HandleFunc("/api/track/sigfox", func(w http.ResponseWriter, r *http.Request) {
		req, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		publish(req)
	})

	port := os.Getenv("API_PORT")
	log.Printf("listening on port %s", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), logHandler(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}

//logHandler logs request infos
func logHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
