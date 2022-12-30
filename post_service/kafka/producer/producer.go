package producer

import (
	"context"
	"encoding/json"
	"exam/post_service/config"
	pb "exam/post_service/genproto/post"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	Conn      *kafka.Conn
	ConnClose func()
}

func NewKafkaProducer(cfg config.Config) (*KafkaProducer, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "kafka:9092", "my_topic", 0)
	if err != nil {
		fmt.Println("error while connection kafka")
		return &KafkaProducer{}, err

	}
	return &KafkaProducer{
		Conn: conn,
		ConnClose: func() {
			conn.Close()
		},
	}, err
}

func (k *KafkaProducer) SendMessage(message *pb.PostReq) error {
	value, err := json.Marshal(message)
	if err != nil {
		fmt.Println("error while marshiling json")
		return err
	}
	_, err = k.Conn.WriteMessages(kafka.Message{
		Value: value,
	})
	if err != nil {
		return err
	}
	fmt.Println("Successfully produce review!!!\n", string(value))
	return err
}
