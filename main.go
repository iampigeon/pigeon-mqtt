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

type service struct {
	Client *iot.Client
}

func (s *service) Approve(content []byte) (valid bool, err error) {
	if content == nil {
		return false, errors.New("Invalid message content")
	}

	fmt.Println(string(content))
	m := new(pigeon.MQTT)

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

	m := new(pigeon.MQTT)
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
	port := flag.Int("port", 9010, "host of the service")
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
}
