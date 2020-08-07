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
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/go-playground/validator.v9"
)

const (
	dateTimeFormat = "2006-01-02 15:04:05"
	sigfoxTopic    = "sigfox_tracker_datas"
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

	switch msg.Topic() {
	case sigfoxTopic:
		handleSigfoxMessage(msg)
	default:
		log.Printf("no handler for topic %s", msg.Topic())
		return

	}
}

//handleSigfoxMessage decodes and saves a sigfox payload
func handleSigfoxMessage(msg mqtt.Message) {
	sigfoxRequest := &requests.Sigfox{}
	decodeRequest(bytes.NewReader(msg.Payload()), sigfoxRequest)
	if true == validateRequest(sigfoxRequest) {
		lat, lon, alt := parseCoords(sigfoxRequest.Data)
		log.Printf("LAT: %f, LON: %f, ALT: %f", lat, lon, alt)

		i, err := strconv.Atoi(sigfoxRequest.Time)
		if err != nil {
			log.Println(err)
			return
		}

		location, err := time.LoadLocation("Europe/Paris")
		if err != nil {
			log.Println(err)
			return
		}

		locale := time.Unix(int64(i), 0).In(location)

		stmt, err := db.Prepare("INSERT INTO position (lat, lon, alt, at, raw, origin) VALUES(?,?,?,?,?,?)")
		if err != nil {
			log.Println(err)
			return
		}

		_, err = stmt.Exec(lat, lon, alt, locale.Format(dateTimeFormat), msg.Payload(), msg.Topic())
		if err != nil {
			log.Println(err)
			return
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

//parseCoords gets a hexa string, splits it in 3 and parses coords by unpacking it
func parseCoords(hexa string) (float32, float32, float32) {
	b, err := hex.DecodeString(hexa)
	if err != nil {
		panic(err)
	}

	var lat, lon, alt float32
	buf := bytes.NewReader(b[:4])
	err = binary.Read(buf, binary.LittleEndian, &lat)
	if err != nil {
		log.Println("binary.Read failed:", err)
	}

	buf = bytes.NewReader(b[4:8])
	err = binary.Read(buf, binary.LittleEndian, &lon)
	if err != nil {
		log.Println("binary.Read failed:", err)
	}

	buf = bytes.NewReader(b[8:])
	err = binary.Read(buf, binary.LittleEndian, &alt)
	if err != nil {
		log.Println("binary.Read failed:", err)
	}

	return lat, lon, alt
}
