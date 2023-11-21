package s3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hanapedia/metrics-processor/internal/domain"
)

type S3Adapter struct {
	client       *s3.S3
	bucketName   string
	keyParentDir string
}

func NewS3Adapter(config *domain.Config) (*S3Adapter, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.AWSRegion)},
	)
	if err != nil {
		return nil, err
	}
	return &S3Adapter{
		client:       s3.New(sess),
		bucketName:   config.S3Bucket,
		keyParentDir: config.TestName,
	}, nil
}

func (sa *S3Adapter) Save(metricsChan <-chan *domain.MetricsMatrix) {
	for metricsMatrix := range metricsChan {
		// Serialize the struct to JSON
		jsonData, err := json.Marshal(metricsMatrix)
		if err != nil {
			slog.Error("Failed to encode to json", "err", err, "name", metricsMatrix.Name)
		}

		key := getS3Key(sa.keyParentDir, metricsMatrix.Name)
		// Upload JSON data to S3
		_, err = sa.client.PutObject(&s3.PutObjectInput{
			Bucket:        aws.String(sa.bucketName),
			Key:           aws.String(key),
			Body:          bytes.NewReader(jsonData),
			ContentLength: aws.Int64(int64(len(jsonData))),
			ContentType:   aws.String("application/json"),
		})
		if err != nil {
			slog.Error("Failed to upload to s3", "err", err, "bucketName", sa.bucketName, "key", key)
		}
	}
}

func getS3Key(prefix, name string) string {
	return fmt.Sprintf("%s/%s.json", prefix, name)
}
