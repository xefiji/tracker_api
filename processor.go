package main

import (
	"api/requests"
	"bytes"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/go-playground/validator.v9"
)

var (
	validate *validator.Validate
	db       *sql.DB
)

func init() {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", user, password, dbname))
	if err != nil {
		panic(err.Error())
	}
}

//messageHandler takes care
var messageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("handling from topic [%s]\n%s\n", msg.Topic(), string(msg.Payload()))

	//todo should switch request according to topic

	sigfoxRequest := &requests.Sigfox{}
	decodeRequest(bytes.NewReader(msg.Payload()), sigfoxRequest)
	if true == validateRequest(sigfoxRequest) {
		lat, lon := parseCoords(sigfoxRequest.Data)
		log.Printf("LAT: %f, LON: %f", lat, lon)
		stmt, err := db.Prepare("INSERT INTO position (lat, lon, at) VALUES(?,?,NOW())")
		if err != nil {
			panic(err.Error())
		}

		_, err = stmt.Exec(lat, lon)
		if err != nil {
			log.Fatal(err)
		}
	}
}

//decodeRequest runs struc decoding and returns a response with error string if it failed, or request
func decodeRequest(rawReq io.Reader, finalReq interface{}) {
	err := json.NewDecoder(rawReq).Decode(&finalReq)
	if err != nil {
		log.Println(err)
	}
}

//validateRequest runs struc validation and returns a response with error string if it failed
func validateRequest(req interface{}) bool {
	validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			log.Println(err)
			return false
		}

		for _, err := range err.(validator.ValidationErrors) {
			log.Println(err)
			return false
		}
	}
	return true
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
