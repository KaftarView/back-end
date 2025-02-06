package application

import (
	"first-project/src/bootstrap"
	"first-project/src/enums"
	"fmt"
	"mime/multipart"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Service struct {
	constants *bootstrap.Constants
	storage   *bootstrap.S3
	clients   *s3.S3
	uploader  *s3manager.Uploader
	buckets   map[enums.BucketType]string
}

func NewS3Service(
	constants *bootstrap.Constants,
	storage *bootstrap.S3,
) *S3Service {
	buckets := make(map[enums.BucketType]string)
	buckets[enums.EventsBucket] = storage.Buckets.EventsBucket
	buckets[enums.PodcastsBucket] = storage.Buckets.PodcastsBucket
	buckets[enums.NewsBucket] = storage.Buckets.NewsBucket
	buckets[enums.JournalsBucket] = storage.Buckets.JournalsBucket
	buckets[enums.ProfilesBucket] = storage.Buckets.ProfilesBucket
	return &S3Service{
		constants: constants,
		storage:   storage,
		buckets:   buckets,
	}
}

func (s3Service *S3Service) setS3Client(bucketType enums.BucketType) {
	bucketTypes := enums.GetAllBucketTypes()
	if !slices.Contains(bucketTypes, bucketType) {
		panic(fmt.Errorf("bucket not exist"))
	}
	if s3Service.uploader != nil && s3Service.clients != nil {
		return
	}
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(s3Service.storage.AccessKey, s3Service.storage.SecretKey, ""),
		Region:      aws.String(s3Service.storage.Region),
		Endpoint:    aws.String(s3Service.storage.Endpoint),
	})

	if err != nil {
		panic(fmt.Errorf("unable to create AWS session, %w", err))
	}

	s3Service.uploader = s3manager.NewUploader(sess)
	s3Service.clients = s3.New(sess)
}

func (s3Service *S3Service) UploadObject(bucketType enums.BucketType, key string, file *multipart.FileHeader) {
	s3Service.setS3Client(bucketType)
	bucket := s3Service.buckets[bucketType]

	fileReader, err := file.Open()
	if err != nil {
		panic(fmt.Errorf("unable to open file %q, %w", file.Filename, err))
	}
	defer fileReader.Close()

	_, err = s3Service.clients.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && (aerr.Code() == s3.ErrCodeNoSuchBucket || aerr.Code() == "NotFound") {
			_, err = s3Service.clients.CreateBucket(&s3.CreateBucketInput{
				Bucket: aws.String(bucket),
			})
			if err != nil {
				panic(fmt.Errorf("unable to create bucket %q, %w", bucket, err))
			}

			err = s3Service.clients.WaitUntilBucketExists(&s3.HeadBucketInput{
				Bucket: aws.String(bucket),
			})
			if err != nil {
				panic(fmt.Errorf("unable to confirm bucket %q exists, %w", bucket, err))
			}
		} else {
			panic(fmt.Errorf("unable to check bucket %q, %w", bucket, err))
		}
	}

	_, err = s3Service.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   fileReader,
	})
	if err != nil {
		panic(fmt.Errorf("unable to upload %q to %q, %w", file.Filename, bucket, err))
	}
}

func (s3Service *S3Service) DeleteObject(bucketType enums.BucketType, key string) error {
	s3Service.setS3Client(bucketType)
	bucket := s3Service.buckets[bucketType]

	_, err := s3Service.clients.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("unable to upload %q to %q, %w", key, bucket, err)
	}

	err = s3Service.clients.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("unable to open file %q, %w", key, err)
	}
	return nil
}

func (s3Service *S3Service) GetPresignedURL(bucketType enums.BucketType, objectKey string, expiration time.Duration) string {
	s3Service.setS3Client(bucketType)
	bucket := s3Service.buckets[bucketType]

	req, _ := s3Service.clients.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	})

	url, err := req.Presign(expiration)
	if err != nil {
		panic(fmt.Errorf("failed to generate presigned URL: %w", err))
	}

	return url
}
