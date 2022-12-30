package kafka

import (
	"exam/review_service/config"
	"exam/review_service/kafka/consumer"
	"exam/review_service/pkg/logger"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type KafkaCon struct {
	KafkaConn *consumer.KafkaConnConsumer
}

type KafkaConnI interface {
	Reads() *consumer.KafkaConnConsumer
}

func NewKafkaReader(cfg config.Config, log logger.Logger, db *sqlx.DB) (KafkaConnI, func(), error) {
	kafka_reader, err := consumer.NewKafkaConsumer(cfg, db)
	if err != nil {
		fmt.Println("error conn kafka")
		return nil, nil, err
	}
	return &KafkaCon{
			KafkaConn: kafka_reader,
		}, func() {
			kafka_reader.Reader.Close()
		}, nil
}

func (k KafkaCon) Reads() *consumer.KafkaConnConsumer {
	return k.KafkaConn
}
