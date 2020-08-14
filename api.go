package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Position struct {
	Lat    string `json:"lat"`
	Lon    string `json:"lon"`
	Alt    string `json:"alt"`
	At     string `json:"at"`
	Origin string `json:"origin"`
	Batt   string `json:"batt"`
}

//api runs the http server for api endpoints
func api() {

	http.HandleFunc("/", nullHandler)
	http.HandleFunc("/api/track/sigfox", sigfoxHandler)
	http.HandleFunc("/api/tracks", tracksHandler)

	port := os.Getenv("API_PORT")
	log.Printf("listening on port %s", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), logHandler(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}

//nullHandler is just a defaut endpoint for / leading to nowhere
func nullHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("No handler"))
	return
}

//sigfoxHandler parses sigfox requests and publishes it to mqtt
func sigfoxHandler(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	publish(sigfoxTopic, req)
}

//logHandler logs request infos
func logHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

//tracksHandler returns the recorded tracks for the given day
func tracksHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet && r.Method != http.MethodOptions {
		http.Error(w, "not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS"))

	day := r.URL.Query().Get("day")
	if day == "" {
		err := errors.New(fmt.Sprintf("url parameter '%s' is missing", "day"))
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	match := regexp.MustCompile(`^[\d]{4}-[\d]{2}-[\d]{2}$`).FindString(day)
	if match == "" {
		err := errors.New(fmt.Sprintf("url parameter '%s' is not well formed. YYYY-MM-DD required", "day"))
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("getting geo datas for day %s", day)

	stmt, err := db.Prepare("SELECT lat, lon, alt, at, origin, batt FROM position WHERE DATE(at) = ? ORDER BY id ASC")
	if err != nil {
		log.Println(err)
		return
	}

	res, err := stmt.Query(day)
	if err != nil {
		log.Println(err)
		return
	}

	positions := make([]Position, 0)
	for res.Next() {
		var pos Position
		err := res.Scan(&pos.Lat, &pos.Lon, &pos.Alt, &pos.At, &pos.Origin, &pos.Batt)
		if err != nil {
			log.Println(err.Error())
		} else {
			positions = append(positions, pos)
		}
	}

	json, err := json.Marshal(positions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
