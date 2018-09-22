package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/WiseGrowth/wisebot-operator/iot"
	"github.com/iampigeon/pigeon"
	"github.com/iampigeon/pigeon/backend"
)

// Message message struct
type Message struct {
	Payload map[string]interface{} `json:"mqtt_payload"`
	Topic   string                 `json:"mqtt_topic"`
}

type service struct {
	Client *iot.Client
}

func (s *service) Approve(content []byte) (valid bool, err error) {
	if content == nil {
		return false, errors.New("Invalid message content")
	}

	fmt.Println(string(content))
	m := new(Message)

	err = json.Unmarshal(content, m)
	if err != nil {
		return false, err
	}

	// validate topic to avoid  breakout
	fmt.Println(m)

	return true, nil
}

func (s *service) Deliver(content []byte) error {
	if !s.Client.IsConnected() {
		// TODO: sentry ? notification to someone
		log.Println("service unavailable")
	}

	m := new(Message)
	err := json.Unmarshal(content, m)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(m.Payload)
	if err != nil {
		return err
	}

	token := s.Client.Publish(m.Topic, byte(1), false, payload)
	if token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	fmt.Println("camilito skt1")

	log.Printf("message received: %s", content)
	return nil
}

func main() {
	host := flag.String("host", "", "host of the service")
	port := flag.Int("port", 5151, "host of the service")
	flag.Parse()

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

	addr := fmt.Sprintf("%s:%d", *host, *port)

	log.Printf("Serving at %s", addr)

	svc := &service{Client: client}

	if err := backend.ListenAndServe(pigeon.NetAddr(addr), svc); err != nil {
		log.Fatal(err)
	}

	//get os env
	//http.ListenAndServe(":5151", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//var msg Message

	//fmt.Println("llego")

	//b, _ := ioutil.ReadAll(r.Body)
	//fmt.Println(string(b))

	//if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
	//http.Error(w, "bad request2", http.StatusBadRequest)
	//fmt.Println("BOOM HERE")
	//fmt.Println(err)
	//return
	//}

	//fmt.Printf("%#v\n", msg)

	//b, err := json.Marshal(msg.Payload)
	//if err != nil {
	//log.Fatal(err)
	//}
}
