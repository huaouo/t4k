package rpc

import (
	"github.com/huaouo/t4k/common"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type MqHandler struct {
	UnimplementedMqServer
	Conn *amqp.Connection
}

func (h *MqHandler) Publish(srv Mq_PublishServer) error {
	ch, err := h.Conn.Channel()
	defer ch.Close()
	if err != nil {
		log.Printf("failed to create new mq channel: %v", err)
		return common.ErrInternal
	}
	for {
		req, err := srv.Recv()
		if err != nil {
			break
		}

		err = ch.Publish(
			"",
			req.GetQueueName(),
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				Body:         req.GetContent(),
			})
		if err != nil {
			log.Printf("failed to publish message to mq: %v", err)
			return common.ErrInternal
		}
	}
	return nil
}

func (h *MqHandler) Subscribe(srv Mq_SubscribeServer) error {
	ch, err := h.Conn.Channel()
	defer ch.Close()
	if err != nil {
		log.Printf("failed to create new mq channel: %v", err)
		return common.ErrInternal
	}
	req, err := srv.Recv()
	if err != nil {
		log.Printf("failed to receive requests: %v", err)
		return common.ErrInternal
	}

	msg, err := ch.Consume(
		req.GetQueueName(), // queue
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("failed to start consume: %v", err)
		return common.ErrInternal
	}

	for d := range msg {
		err := srv.Send(&SubResponse{
			Content: d.Body,
		})
		if err != nil {
			log.Printf("failed to send to client: %v", err)
			return common.ErrInternal
		}
		_, err = srv.Recv()
		if err != nil {
			log.Printf("failed to receive ack from client: %v", err)
			return common.ErrInternal
		}
		err = d.Ack(false)
		if err != nil {
			log.Printf("failed to ack mq: %v", err)
			return common.ErrInternal
		}
	}
	return nil
}
