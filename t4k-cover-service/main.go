package main

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/huaouo/t4k/common"
	"github.com/huaouo/t4k/t4k-mq-service/rpc"
	"log"
	"os"
)

func main() {
	mqClient := rpc.NewMqClient(rpc.NewMqServiceClient())
	pollMq(mqClient)
}

func pollMq(mqClient rpc.MqClient) {
	endpoint := os.Getenv("OBJECT_SERVICE_ADDR") + ":" + os.Getenv("OBJECT_SERVICE_LISTEN_PORT")
	coverGen := CoverGenerator{
		VideoUrlPrefix: "http://" + endpoint + common.ObjectServiceVideoPathPrefix,
		CoverUrlPrefix: "http://" + endpoint + common.ObjectServiceCoverPathPrefix,
	}
	subClient, err := mqClient.Subscribe(context.TODO())
	if err != nil {
		log.Printf("failed to subscribe to MQ: %v", err)
		return
	}
	err = subClient.Send(&rpc.SubRequestOrAck{
		QueueName: aws.String(common.MqCoverQueueName),
	})
	if err != nil {
		log.Printf("failed to initialize subscription: %v", err)
		return
	}

	for {
		resp, err := subClient.Recv()
		if err != nil {
			log.Printf("failed to receive cover task: %v", err)
			break
		}
		objectId := string(resp.Content)
		err = coverGen.GenerateCover(objectId)
		if err != nil {
			log.Printf("failed to generate cover: %v", err)
			break
		}
		err = subClient.Send(&rpc.SubRequestOrAck{})
		if err != nil {
			log.Printf("failed to ack MQ: %v", err)
			break
		}
	}
}
