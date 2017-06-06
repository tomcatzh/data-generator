package awsstorage

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// A StorageS3 is handler of aws S3
type StorageS3 struct {
	config *aws.Config
	bucket string
}

// NewStorageS3 returns StorageS3 handler
func NewStorageS3(region string, bucket string) *StorageS3 {
	return &StorageS3{
		config: &aws.Config{Region: aws.String(region)},
		bucket: bucket,
	}
}

// Save function uploads the reader to S3
func (s *StorageS3) Save(key string, reader io.ReadSeeker) (int64, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("cn-north-1"),
	}))

	uploader := s3manager.NewUploader(sess)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: &s.bucket,
		Key:    &key,
		Body:   reader,
	})
	if err != nil {
		return 0, err
	}

	svc := s3.New(sess)

	resp, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err != nil {
		return 0, err
	}

	return *resp.ContentLength, nil
}
