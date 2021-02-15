package message_queue

import (
	"fmt"
	"github.com/golang/glog"
	"log"
	"mlcp/pkg/common"
	"mlcp/pkg/config"
	"mlcp/pkg/worker"
	"time"
	"github.com/streadway/amqp"
)

type MessageQueue interface{
	Produce(data interface{}) error
	Consume() (interface{}, error)
	ListenToQueue(wq *worker.WorkeQueue)
}

type RabbitMQ struct {
	Host string
	QueueName string
}

func InitMQ() (MessageQueue, error) {
	switch config.MQDriver {
	case "rabbitmq":
		return initRabbitMQ(), nil
	}
	return nil, fmt.Errorf("Messagequeue not supported: %s", config.MQDriver)
}

func (rbt *RabbitMQ) ListenToQueue(wq *worker.WorkeQueue) {
	for {
		go func(rbt *RabbitMQ) {
			conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
			if err != nil {
				glog.Errorf("Failed to connect to RabbitMQ: %v", err)
				return
			}
			defer conn.Close()

			ch, err := conn.Channel()
			if err != nil {
				glog.Errorf("Failed to open a channel: %v", err)
				return
			}
			defer ch.Close()

			q, err := ch.QueueDeclare(
				rbt.QueueName, // name
				true,          // durable
				false,         // delete when unused
				false,         // exclusive
				false,         // no-wait
				nil,           // arguments
			)
			if err != nil {
				glog.Errorf("Failed to declare a queue: %v", err)
				return
			}

			err = ch.Qos(
				1,     // prefetch count
				0,     // prefetch size
				false, // global
			)
			if err != nil {
				glog.Errorf("Failed to set QoS", err)
				return
			}
			msgs, err := ch.Consume(
				q.Name, // queue
				"",     // consumer
				false,  // auto-ack
				false,  // exclusive
				false,  // no-local
				false,  // no-wait
				nil,    // args
			)
			if err != nil {
				glog.Errorf("Failed to register a consumer: %v", err)
				return
			}

			forever := make(chan bool)
			go func() {
				for d := range msgs {
					r, err := common.ParseRequest(d.Body)
					if err != nil {
						glog.Errorf("Error parsing request body: %v", err)
					}
					wq.Add(interface{}(common.NewSlot(common.NewCar(r.RegnNo), r.SlotId)))
				}
			}()

			log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
			<-forever
		}(rbt)
		time.Sleep(1*time.Second)
	}
}

func initRabbitMQ() *RabbitMQ {
	rbt := &RabbitMQ {
		Host:      config.MQHost,
		QueueName: config.QueueName,
	}
	return rbt
}

func (rmq *RabbitMQ) Produce(data interface{}) error {
	return nil
}

func (rmq *RabbitMQ) Consume() (interface{}, error) {

	return nil, nil
}