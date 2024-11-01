package application_aws

import (
	"first-project/src/bootstrap"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AWSS3 struct {
	bucket *bootstrap.Bucket
}

func NewAWSS3(bucket *bootstrap.Bucket) *AWSS3 {
	return &AWSS3{
		bucket: bucket,
	}
}
func (a *AWSS3) CreateBucket() {
	sess, _ := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(a.bucket.AccessKey, a.bucket.SecretKey, ""),
	})
	svc := s3.New(sess, &aws.Config{
		Region:   aws.String(a.bucket.Region),
		Endpoint: aws.String(a.bucket.Endpoint),
	})

	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(a.bucket.Name),
	})
	if err != nil {
		panic(fmt.Errorf("unable to create bucket %q, %v", a.bucket.Name, err))
	}

	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(a.bucket.Name),
	})
	if err != nil {
		panic(fmt.Errorf("error occurred while waiting for bucket to be created, %q, %v", a.bucket.Name, err))
	}
}

func (a *AWSS3) UploadObject(objectPath string) {
	file, err := os.Open(objectPath)
	if err != nil {
		panic(fmt.Errorf("unable to open file %q, %v", objectPath, err))
	}
	defer file.Close()

	sess, _ := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(a.bucket.AccessKey, a.bucket.SecretKey, ""),
		Region:      aws.String(a.bucket.Region),
		Endpoint:    aws.String(a.bucket.Endpoint),
	})

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(a.bucket.Name),
		Key:    aws.String(filepath.Base(objectPath)),
		Body:   file,
	})
	if err != nil {
		panic(fmt.Errorf("unable to upload %q to %q, %v", objectPath, a.bucket.Name, err))
	}
}
