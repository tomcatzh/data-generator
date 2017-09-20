package storage

import (
	"bytes"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func TestS3Save(t *testing.T) {
	var buf bytes.Buffer
	const body = "abcdefghijklmnopqrstuvwxyz0123456789"
	const times = 10000
	for i := 0; i < times; i++ {
		buf.WriteString(body)
	}

	region := "cn-north-1"
	bucket := "live"
	key := "awsstorage_test"

	template := map[string]interface{}{}
	template["Region"] = region
	template["Bucket"] = bucket
	template["PartSizeM"] = (float64)(64)

	s, err := newStorageS3(template)
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}

	l, err := s.Save(key, bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	} else if l != int64(len(body)*times) {
		t.Errorf("Unexcepted content length: %v", l)
	}

	config := &aws.Config{
		Region: aws.String(region),
	}

	sess := session.Must(session.NewSession(config))
	svc := s3.New(sess)

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
}
