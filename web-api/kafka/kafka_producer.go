package kafka

import (
	"time"

	"go-unit-test/web-api/logger"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

type KafkaProducer interface {
	SendMessage(topic string, msg string) error
}

type KafkaProducerImpl struct {
	producer sarama.SyncProducer
	logger   logger.Logger
}

func NewKafkaProducer(url string, logger logger.Logger) KafkaProducer {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Net.WriteTimeout = 5 * time.Second
	kafkaConfig.Producer.Retry.Max = 0

	producers, err := sarama.NewSyncProducer([]string{url}, kafkaConfig)
	if err != nil {
		panic(err)
	}

	return &KafkaProducerImpl{
		producer: producers,
		logger:   logger,
	}
}

func (kafkaProducer KafkaProducerImpl) SendMessage(topic string, msg string) error {
	context := "kafkaProducer-SendMessage"
	kafkaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}

	_, _, err := kafkaProducer.producer.SendMessage(kafkaMsg)
	if err != nil {
		return err
	}

	kafkaProducer.logger.GetLogger().Info(
		"Send message success",
		zap.String("context", context),
		zap.String("topic", topic),
		zap.String("data", msg),
	)

	return nil
}
