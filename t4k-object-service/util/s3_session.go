package util

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"os"
)

func NewS3Session() *session.Session {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "default",
		Config: aws.Config{
			Credentials:      credentials.NewStaticCredentials(os.Getenv("S3_ACCESS_KEY"), os.Getenv("S3_SECRET_KEY"), ""),
			Endpoint:         aws.String(os.Getenv("S3_ENDPOINT")),
			Region:           aws.String("S3_REGION"),
			S3ForcePathStyle: aws.Bool(true),
		},
	})
	if err != nil {
		log.Fatalf("failed to initialize new s3 session: %v", err)
	}

	return sess
}
