package replicator

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Rbsn-joses/kafka-replicator/utils"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var kafkaProducerBoostrapServers = os.Getenv("KAFKA_PRODUCER_BOOSTRAP_SERVERS")
var kafkaProducerTopic = os.Getenv("KAFKA_PRODUCER_TOPIC")
var kafkaProducerTlsEnabled = os.Getenv("KAFKA_PRODUCER_TLS_ENABLED")
var kafkaProducerSslCertificateLocation = os.Getenv("KAFKA_PRODUCER_SSL_CERTIFICATE_LOCATION")
var kafkaProducerSslPrivateKeyLocation = os.Getenv("KAFKA_PRODUCER_SSL_PRIVATE_KEY_LOCATION")
var kafkaProducerSslPrivateKeyPassword = os.Getenv("KAFKA_PRODUCER_SSL_PRIVATE_KEY_LOCATION")
var kafkaProducerReplicationFactor = utils.GetenvOrDefaultValue("KAFKA_PRODUCER_REPLICATOR_FACTOR", "3")
var kafkaProducerPartitionNumber = utils.GetenvOrDefaultValue("KAFKA_PRODUCER_REPLICATOR_PARTITION_NUMBER", "3")
var kafkaProducerTimeout = utils.GetenvOrDefaultValue("KAFKA_PRODUCER_TIMEOUT", "60s")

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
func CreateTopic() {
	var err error
	var a *kafka.AdminClient

	if kafkaProducerTlsEnabled == "true" {
		a, err = kafka.NewAdminClient(&kafka.ConfigMap{
			"bootstrap.servers":        kafkaProducerBoostrapServers,
			"ssl.key.location":         kafkaProducerSslPrivateKeyLocation,
			"ssl.certificate.location": kafkaProducerSslCertificateLocation,
			"ssl.key.password":         kafkaProducerSslPrivateKeyPassword,
			"security.protocol":        "SSL",
		})
		if err != nil {
			fmt.Printf("Failed to create Admin client: %s\n", err)
			os.Exit(1)
		}
	} else {
		a, err = kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": kafkaProducerBoostrapServers})
		if err != nil {
			fmt.Printf("Failed to create Admin client: %s\n", err)
			os.Exit(1)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	maxDur, err := time.ParseDuration(kafkaProducerTimeout)
	if err != nil {
		panic(err)
	}
	numParts, err := strconv.Atoi(kafkaProducerPartitionNumber)
	if err != nil {
		fmt.Printf("Invalid partition count: %s: %v\n", os.Args[3], err)
		os.Exit(1)
	}
	replicationFactor, err := strconv.Atoi(kafkaProducerReplicationFactor)
	if err != nil {
		fmt.Printf("Invalid replication factor: %s: %v\n", os.Args[4], err)
		os.Exit(1)
	}

	results, err := a.CreateTopics(
		ctx,
		// Multiple topics can be created simultaneously
		// by providing more TopicSpecification structs here.
		[]kafka.TopicSpecification{{
			Topic:             kafkaProducerTopic,
			NumPartitions:     numParts,
			ReplicationFactor: replicationFactor,
		}},
		// Admin options
		kafka.SetAdminOperationTimeout(maxDur))
	if err != nil {
		fmt.Printf("Failed to create topic: %v\n", err)
		os.Exit(1)
	}

	// Print results
	for _, result := range results {
		fmt.Printf("%s\n", result)
	}

	a.Close()
}
func SendTopicMessages(value []byte) {
	var err error
	var p *kafka.Producer
	if kafkaProducerTlsEnabled == "true" {
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
	} else {
		p, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaProducerBoostrapServers})

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

	if err != nil {
		fmt.Printf("Failed to create topic: %v\n", err)
		os.Exit(1)
	}
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
