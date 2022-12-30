package consumer

import (
	"context"
	"encoding/json"
	"exam/review_service/config"
	ps "exam/review_service/genproto/post"
	rs "exam/review_service/genproto/review"
	"exam/review_service/pkg/logger"
	"exam/review_service/storage"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/segmentio/kafka-go"
)

type KafkaConnConsumer struct {
	Reader    *kafka.Reader
	ConnClose func()
	Cfg       config.Config
	Logger    logger.Logger
	Storage   storage.IStorage
}

func NewKafkaConsumer(cfg config.Config, db *sqlx.DB) (*KafkaConnConsumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"kafka:9092"},
		Topic:     "my_topic",
		Partition: 0,
		MinBytes:  1e3,
		MaxBytes:  10e6,
	})
	return &KafkaConnConsumer{
		Reader:  reader,
		Storage: storage.NewStoragePg(db),
		ConnClose: func() {
			reader.Close()
		},
		Cfg: cfg,
	}, nil
}

func (k *KafkaConnConsumer) Start() error {
	fmt.Println("Start waiting for messages >>")
	for {
		m, err := k.Reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("error while reading message")
			return err
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		req := &ps.PostReq{}
		err = json.Unmarshal(m.Value, req)
		if err != nil {
			fmt.Println("error while unmarshiling kafka message")
			return err
		}
		for _, review := range req.Reviews {
			reviews := &rs.ReviewRequest{
				Id:          review.Id,
				OwnerId:     review.OwnerId,
				PostId:      review.PostId,
				Name:        review.Name,
				Rating:      review.Rating,
				Description: review.Description,
			}
			k.Storage.Review().CreateReview(reviews)
		}
		fmt.Println("Start waiting message!!!")
	}
}
