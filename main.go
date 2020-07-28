package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
)

type sigfoxRequest struct {
	Data         string `json:"data"`
	Rssi         string `json:"rssi"`
	SeqNumber    string `json:"seqNumber"`
	DeviceTypeId string `json:"deviceTypeId"`
	Id           string `json:"id"`
	Time         string `json:"time"`
	Snr          string `json:"snr"`
	Station      string `json:"station"`
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No handler"))
	})

	http.HandleFunc("/api/track/sigfox", func(w http.ResponseWriter, r *http.Request) {

		req := sigfoxRequest{}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("%+v\n", req)
		//todo do something like mqtt or notify

	})

	log.Fatal(http.ListenAndServe(":80", nil))
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
