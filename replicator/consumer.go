package replicator

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func init() {
	if kafkaConsumerBoostrapServers == "" {
		log.Panic("Variável KAFKA_CONSUMER_BOOSTRAP_SERVERS não declarada")
	}
	if kafkaConsumerTopic == "" {
		log.Panic("Variável KAFKA_CONSUMER_TOPIC não declarada")

	}
	if kafkaConsumerGroupID == "" {
		log.Panic("Variável KAFKA_CONSUMER_GROUP_ID não declarada")

	}
	if kafkaConsumerTlsEnabled == "true" {
		if kafkaConsumerSslCertificateLocation == "" {
			log.Panic("Variável KAFKA_CONSUMER_SSL_CERTIFICATE_LOCATION não declarada")

		}
		if kafkaConsumerSslPrivateKeyLocation == "" {
			log.Panic("Variável KAFKA_CONSUMER_SSL_PRIVATE_KEY_LOCATION não declarada")

		}
		if kafkaConsumerSslPrivateKeyPassword == "" {
			log.Panic("Variável KAFKA_CONSUMER_SSL_PRIVATE_KEY_PASSWORD não declarada")

		}

	}

}

var kafkaConsumerBoostrapServers = os.Getenv("KAFKA_CONSUMER_BOOSTRAP_SERVERS")
var kafkaConsumerTopic = os.Getenv("KAFKA_CONSUMER_TOPIC")
var kafkaConsumerTlsEnabled = os.Getenv("KAFKA_CONSUMER_TLS_ENABLED")
var kafkaConsumerSslCertificateLocation = os.Getenv("KAFKA_CONSUMER_SSL_CERTIFICATE_LOCATION")
var kafkaConsumerSslPrivateKeyLocation = os.Getenv("KAFKA_CONSUMER_SSL_PRIVATE_KEY_LOCATION")
var kafkaConsumerSslPrivateKeyPassword = os.Getenv("KAFKA_CONSUMER_SSL_PRIVATE_KEY_PASSWORD")
var replicationEnabled = os.Getenv("KAFKA_REPLICATION_ENABLED")
var kafkaConsumerGroupID = os.Getenv("KAFKA_CONSUMER_GROUP_ID")

func init() {

}
func GetTopicMessages() {
	var err error
	var c *kafka.Consumer
	if kafkaConsumerTlsEnabled == "true" {
		c, err = kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers":        kafkaConsumerBoostrapServers,
			"group.id":                 kafkaConsumerGroupID,
			"security.protocol":        "SSL",
			"ssl.key.location":         kafkaConsumerSslPrivateKeyLocation,
			"ssl.certificate.location": kafkaConsumerSslCertificateLocation,
			"ssl.key.password":         kafkaConsumerSslPrivateKeyPassword,
			"auto.offset.reset":        "earliest",
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
			os.Exit(1)
		}
	} else {
		c, err = kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers": kafkaConsumerBoostrapServers,
			"group.id":          kafkaConsumerGroupID,
			"auto.offset.reset": "earliest",
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("Created Consumer %v\n", c)
	err = c.SubscribeTopics(strings.Split(kafkaConsumerTopic, ","), nil)

	//run := true

	if err != nil {
		panic(err)
	}
	run := true
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}
				if replicationEnabled == "true" {
					SendTopicMessages(e.Value)
				}

			case kafka.Error:

				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	c.Close()
}
