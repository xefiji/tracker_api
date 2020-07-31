package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

var (
	validate *validator.Validate
)

//decodeRequest runs struc decoding and returns a response with error string if it failed, or request
func decodeRequest(w http.ResponseWriter, r *http.Request, req interface{}) interface{} {
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	return req
}

//validateRequest runs struc validation and returns a response with error string if it failed
func validateRequest(w http.ResponseWriter, req interface{}) {
	validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		for _, err := range err.(validator.ValidationErrors) {
			log.Println(err)
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
			return
		}
	}
}

//parseCoords gets a hexa string, splits it in two and parses coords by unpacking it
func parseCoords(hexa string) (float32, float32) {
	b, err := hex.DecodeString(hexa)
	if err != nil {
		panic(err)
	}

	var lat, lon float32
	buf := bytes.NewReader(b[:4])
	err = binary.Read(buf, binary.LittleEndian, &lat)
	if err != nil {
		log.Println("binary.Read failed:", err)
	}

	buf = bytes.NewReader(b[4:])
	err = binary.Read(buf, binary.LittleEndian, &lon)
	if err != nil {
		log.Println("binary.Read failed:", err)
	}

	return lat, lon
}
