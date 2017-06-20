package data

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type storageS3 struct {
	config   *aws.Config
	bucket   string
	partSize int64
}

// NewStorageS3 returns storageS3 handler
func newStorageS3(region string, bucket string, partSize int64) *storageS3 {
	return &storageS3{
		config:   &aws.Config{Region: aws.String(region)},
		bucket:   bucket,
		partSize: partSize,
	}
}

// Save function uploads the reader to S3
func (s *storageS3) Save(key string, reader io.Reader) (int64, error) {
	sess := session.Must(session.NewSession(s.config))

	svc := s3.New(sess)

	uploader := s3manager.NewUploaderWithClient(svc, func(u *s3manager.Uploader) {
		if s.partSize > 5*1024*1024 {
			u.PartSize = s.partSize
			u.Concurrency = 1
		}
		u.MaxUploadParts = int((5 * 1024 * 1024 * 1024 * 1024) / u.PartSize)
	})
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: &s.bucket,
		Key:    &key,
		Body:   reader,
	})
	if err != nil {
		return 0, err
	}

	resp, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err != nil {
		return 0, err
	}

	return *resp.ContentLength, nil
}