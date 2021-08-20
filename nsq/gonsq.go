package main

import (
	"flag"
	"log"

	"github.com/nsqio/go-nsq"
)

var (
	role = flag.String("role", "", "input your role")
)

const (
	nsqdAddr       = "127.0.0.1:4150"
	nsqlookupdAddr = "127.0.0.1:4161"
)

func main() {
	flag.Parse()
	topic := "lab_first"

	switch *role {
	case "producer":
		producer, err := nsq.NewProducer(nsqdAddr, nsq.NewConfig())
		if err != nil {
			panic(err)
		}

		producer.Publish(topic, []byte("this is the first message!"))

		producer.Stop()
	case "consumer":
		consumer, err := nsq.NewConsumer(topic, "channel-01", nsq.NewConfig())
		if err != nil {
			panic(err)
		}
		ch := make(chan *nsq.Message, 1)
		consumer.AddHandler(makeNsqHandler(ch))

		if err := consumer.ConnectToNSQLookupd(nsqlookupdAddr); err != nil {
			panic(err)
		}

		select {
		case msg := <-ch:
			log.Println("Consumer: ", string(msg.Body))
		}
	default:
		flag.Usage()
	}
}

type funcHandler func(message *nsq.Message) error

func (f funcHandler) HandleMessage(message *nsq.Message) error {
	return f(message)
}

func makeNsqHandler(chmsg chan *nsq.Message) nsq.Handler {
	return funcHandler(func(message *nsq.Message) error {
		log.Println("received messge")
		chmsg <- message
		return nil
	})
}
