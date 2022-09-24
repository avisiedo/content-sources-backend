package event

import (
	"github.com/content-services/content-sources-backend/pkg/config"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"
)

// Adapted from: https://github.com/RedHatInsights/playbook-dispatcher/blob/master/internal/response-consumer/main.go#L21
// func Start(
// 	ctx context.Context,
// 	cfg *config.Configuration,
// 	handler Eventable,
// ) {
// 	var (
// 		// schemas schema.TopicSchemas
// 		err error
// 	)
// 	// schemas, err = schema.LoadSchemas()
// 	// utils.DieOnError(err)

// 	consumer, err := NewConsumer(ctx, cfg)
// 	utils.DieOnError(err)

// 	start := NewConsumerEventLoop(ctx, consumer /* nil, nil, schemas*/, handler.OnMessage)
// 	start()
// }

func Start(config *config.Configuration, handler Eventable) {
	var (
		err error
		// msg      *kafka.Message
		consumer *kafka.Consumer
	)

	if consumer, err = NewConsumer(config); err != nil {
		log.Logger.Panic().Msg("[Start] error creating consumer")
		return
	}
	defer consumer.Close()

	start := NewConsumerEventLoop(consumer, handler)
	start()
}
