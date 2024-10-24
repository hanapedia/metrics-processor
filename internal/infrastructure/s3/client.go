package s3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
		keyParentDir: config.S3BucketDir,
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

func (sa *S3Adapter) ParseEndTime() (float64, error) {
	// List the first file in the bucket with the prefix
	resp, err := sa.client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(sa.bucketName),
		Prefix: aws.String(sa.keyParentDir),
		MaxKeys: aws.Int64(1), // To get the first file
	})
	if err != nil {
		return 0, fmt.Errorf("Unable to list items in bucket %q, %w", sa.bucketName, err)
	}

	// Ensure at least one file is returned
	if len(resp.Contents) == 0 {
		return 0, fmt.Errorf("No files found with prefix %s", sa.keyParentDir)
	}

	// Get the first file's key
	fileKey := *resp.Contents[0].Key
	slog.Info("Found file", "file", fileKey)

	// Fetch the file content
	getResp, err := sa.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(sa.bucketName),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		return 0, fmt.Errorf("Unable to download item %q, %w", fileKey, err)
	}
	defer getResp.Body.Close()

	// Read the file content
	body, err := io.ReadAll(getResp.Body)
	if err != nil {
		return 0, fmt.Errorf("Failed to read file content, %w", err)
	}

	// Unmarshal the content into your struct
	var data domain.MetricsMatrix
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, fmt.Errorf("Failed to unmarshal JSON, %w", err)
	}
	slog.Info("Successfuly parsed endtime","bucket", sa.bucketName, "parentDir", sa.keyParentDir, "file", fileKey, "end", data.End)
	return data.End, nil
}

func getS3Key(prefix, name string) string {
	return fmt.Sprintf("%s/%s.json", prefix, name)
}
