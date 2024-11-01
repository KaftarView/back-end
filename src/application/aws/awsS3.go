package application_aws

import (
	"first-project/src/bootstrap"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
