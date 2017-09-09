package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/WiseGrowth/wisebot-operator/iot"
)

const (
	iotHost     = ""
	certificate = ""
	privateKey  = ""
)

// Message message struct
type Message struct {
	Payload map[string]interface{} `json:"payload"`
	Topic   string                 `json:"topic"`
}

func main() {
	cert, err := tls.X509KeyPair([]byte(certificate), []byte(privateKey))
	if err != nil {
		panic(err)
	}

	client, err := iot.NewClient(
		iot.SetHost(iotHost),
		iot.SetCertificate(cert),
		iot.SetClientID("pigeon-"+strconv.Itoa(rand.Int())),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe(":5151", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var msg Message

		fmt.Println("llego")

		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, "bad request2", http.StatusBadRequest)
			fmt.Println(err)
			return
		}

		fmt.Printf("%#v\n", msg)

		b, err := json.Marshal(msg.Payload)
		if err != nil {
			log.Fatal(err)
		}

		token := client.Publish(msg.Topic, byte(1), false, b)
		if token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		}
	}))
}
