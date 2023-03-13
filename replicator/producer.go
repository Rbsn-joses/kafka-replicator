package replicator

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var kafkaProducerBoostrapServers = os.Getenv("KAFKA_PRODUCER_BOOSTRAP_SERVERS")
var kafkaProducerTopic = os.Getenv("KAFKA_PRODUCER_TOPIC")
var kafkaProducerTlsEnabled = os.Getenv("KAFKA_PRODUCER_TLS_ENABLED")
var kafkaProducerSslCertificateLocation = os.Getenv("KAFKA_PRODUCER_SSL_CERTIFICATE_LOCATION")
var kafkaProducerSslPrivateKeyLocation = os.Getenv("KAFKA_PRODUCER_SSL_PRIVATE_KEY_LOCATION")
var kafkaProducerSslPrivateKeyPassword = os.Getenv("KAFKA_PRODUCER_SSL_PRIVATE_KEY_LOCATION")

func init() {

	if kafkaProducerBoostrapServers == "" {
		log.Panic("Variável KAFKA_PRODUCER_BOOSTRAP_SERVERS não declarada")
	}
	if kafkaProducerTopic == "" {
		log.Panic("Variável KAFKA_PRODUCER_TOPIC não declarada")

	}

	if kafkaProducerTlsEnabled == "true" {
		if kafkaProducerSslCertificateLocation == "" {
			log.Panic("Variável KAFKA_PRODUCER_SSL_CERTIFICATE_LOCATION não declarada")

		}
		if kafkaProducerSslPrivateKeyLocation == "" {
			log.Panic("Variável KAFKA_PRODUCER_SSL_PRIVATE_KEY_LOCATION não declarada")

		}
		if kafkaProducerSslPrivateKeyPassword == "" {
			log.Panic("Variável KAFKA_PRODUCER_SSL_PRIVATE_KEY_PASSWORD não declarada")

		}

	}

}
func SendTopicMessages(value []byte) {
	var err error
	var p *kafka.Producer
	if kafkaProducerTlsEnabled == "true" {
		p, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaProducerBoostrapServers})

		if err != nil {
			fmt.Printf("Failed to create producer: %s\n", err)
			os.Exit(1)
		}
	} else {
		p, err = kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers":        kafkaProducerBoostrapServers,
			"ssl.key.location":         kafkaProducerSslPrivateKeyLocation,
			"ssl.certificate.location": kafkaProducerSslCertificateLocation,
			"ssl.key.password":         kafkaProducerSslPrivateKeyPassword,
			"security.protocol":        "SSL",
		})

		if err != nil {
			fmt.Printf("Failed to create producer: %s\n", err)
			os.Exit(1)
		}
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case kafka.Error:
				fmt.Printf("Error: %v\n", ev)
			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}
	}()
	deliveryChan := make(chan kafka.Event)
	go func() {
		for e := range deliveryChan {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					fmt.Printf("Replicando mensagem %s no tópico %s\n", string(value), kafkaProducerTopic)
				}

			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
			close(deliveryChan)
		}
	}()
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kafkaProducerTopic, Partition: kafka.PartitionAny},
		Value:          []byte(value),
		Headers:        []kafka.Header{{Key: "replicator", Value: []byte("kafka to kafka")}},
	}, deliveryChan)
	if err != nil {
		close(deliveryChan)
		if err.(kafka.Error).Code() == kafka.ErrQueueFull {
			time.Sleep(time.Second)
		}
		fmt.Printf("Failed to produce message: %v\n", err)
	}

}
