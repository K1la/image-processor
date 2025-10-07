package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/K1la/image-processor/internal/config"
	"github.com/K1la/image-processor/internal/model"
	"github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
	"time"
)

type Kafka struct {
	consumer *kafka.Consumer
	producer *kafka.Producer
}

func New(cfg config.Kafka) *Kafka {
	address := cfg.Host + cfg.Port

	consumer := kafka.NewConsumer([]string{address}, "images", "imgs")
	producer := kafka.NewProducer([]string{address}, "images")

	return &Kafka{consumer: consumer, producer: producer}
}

func (k *Kafka) ProduceMessage(message model.Message) error {
	strategy := retry.Strategy{
		Attempts: 3,
		Delay:    time.Second,
		Backoff:  1,
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not marshal kafka message payload to produce: %w", err)
	}

	if err = k.producer.SendWithRetry(context.Background(), strategy, nil, payload); err != nil {
		return fmt.Errorf("could not send message to kafka: %w", err)
	}

	zlog.Logger.Info().Str("message id", message.ID.String()).Msg("successfully produced message to kafka with id <-")
	return nil
}

func (k *Kafka) ConsumeMessage() (*model.Message, error) {
	strategy := retry.Strategy{
		Attempts: 3,
		Delay:    time.Second,
		Backoff:  1,
	}

	kafkaMessage, err := k.consumer.FetchWithRetry(context.Background(), strategy)
	if err != nil {
		return nil, fmt.Errorf("could not consume message from kafka: %w", err)
	}

	var message model.Message
	if err = json.Unmarshal(kafkaMessage.Value, &message); err != nil {
		return nil, fmt.Errorf("could not unmarshal kafka message payload: %w", err)
	}

	if err = k.consumer.Commit(context.Background(), kafkaMessage); err != nil {
		return nil, fmt.Errorf("could not commit kafka message: %w", err)
	}

	zlog.Logger.Info().Msg("successfully consumed message from kafka with id: " + message.ID.String())
	return &message, nil
}
