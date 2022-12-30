package kafka

import (
	"exam/post_service/config"
	"exam/post_service/kafka/producer"
	"fmt"
)

type Kafka struct {
	KafkaFuncs *producer.KafkaProducer
}

type KafkaI interface {
	Produce() *producer.KafkaProducer
}

func NewKafka(cfg config.Config) (KafkaI, func(), error) {
	kafka, err := producer.NewKafkaProducer(cfg)
	if err != nil {
		fmt.Println("error conn kafka producer")
		return &Kafka{}, func() {}, err
	}
	return &Kafka{
			KafkaFuncs: kafka,
		}, func() {
			kafka.ConnClose()
		}, err
}

func (k *Kafka) Produce() *producer.KafkaProducer {
	return k.KafkaFuncs
}
