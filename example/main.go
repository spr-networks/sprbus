package main

import (
	"github.com/spr-networks/sprbus"
	//"github.com/spr-networks/sprbus/log"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

var client *sprbus.Client

// var socket = "/tmp/test.sock"
var socket = os.Getenv("TEST_PREFIX") + "/state/api/eventbus.sock"

func custom_subscribe(client *sprbus.Client) {
	time.Sleep(time.Second / 4)

	var wg sync.WaitGroup

	stream, err := client.SubscribeTopic("spr")
	if nil != err {
		log.Fatal(err)
	}

	go func() {
		wg.Add(1)
		for {
			reply, err := stream.Recv()
			if io.EOF == err {
				break
			}

			if nil != err {
				return
			}

			topic := reply.GetTopic()
			value := reply.GetValue()

			fmt.Printf("topic=%v value=%v\n", topic, value)
		}
	}()

}

func custom_publish(client *sprbus.Client) {
	for i := 0; i < 5; i++ {
		_, err := client.Publish("spr:test", "{\"data\": \"test\"}")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func custom_server() {
	//log.Println("server listening...")

	server, err := sprbus.NewServer(socket)
	if err != nil {
		log.Fatal(err)
	}

	// does not return
	log.Println("server:", server)
}

func spr_publish() {
	for i := 0; i < 5; i++ {
		type testS struct {
			Title string
			Body  string
		}

		sprbus.Publish("spr:test:struct", testS{Title: "tttt1111", Body: "datahere"})
		sprbus.Publish("spr:test:string", "s1ACID")
		sprbus.Publish("spr:test:array", []int{11, 23})
	}
}

func spr_event() {
	sprbus.HandleEvent("", func(topic string, json string) {
		fmt.Printf("[sprbus] %v %v\n", topic, json)
	})
}

func spr_log() {
	var log = sprbus.NewLog("spr:log")
	// can modify log
	log.SetLevel(logrus.DebugLevel)
	/*log.SetOutput(os.Stdout)
	log.SetReportCaller(false)
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})*/

	log.Warnf("this is a warning: %v", 1234)
	log.Println("connected to", socket)
	log.Debugf("debug: %v", 1234)
}

func main() {

	go custom_server()
	time.Sleep(time.Second / 2)
	go spr_event()
	time.Sleep(time.Second / 2)

	client, err := sprbus.NewClient(socket)
	defer client.Close()

	if err != nil {
		log.Fatal("err", err)
	}

	//custom_subscribe(client)
	//custom_publish(client)
	spr_publish()

	spr_log()

	time.Sleep(time.Second * 5)
}
