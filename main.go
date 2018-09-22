package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/WiseGrowth/wisebot-operator/iot"
)

// Message message struct
type Message struct {
	Payload map[string]interface{} `json:"mqtt_payload"`
	Topic   string                 `json:"mqtt_topic"`
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

	//get os env
	http.ListenAndServe(":5151", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var msg Message

		fmt.Println("llego")

		b, _ := ioutil.ReadAll(r.Body)
		fmt.Println(string(b))

		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, "bad request2", http.StatusBadRequest)
			fmt.Println("BOOM HERE")
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
		fmt.Println("camilito skt1")
	}))
}
