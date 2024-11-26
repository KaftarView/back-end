package application_aws

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

// type S3service struct {
// constants *bootstrap.Constants
// 	bucket    *bootstrap.Bucket // has to change to 3 different bucket names
// 	session   *session.Session
// 	s3Client  *s3.S3
// }

// func NewAWSS3(constants *bootstrap.Constants, bucket *bootstrap.Bucket) *S3service {
// 	return &S3service{
// 		constants: constants,
// 		bucket:    bucket,
// 	}
// }

// type S3service struct {
// 	constants   *bootstrap.Constants
// 	bucketConfigs map[enums.BucketType]*bootstrap.Bucket
// 	clients     map[enums.BucketType]*s3.S3
// 	uploader    map[enums.BucketType]*s3manager.Uploader
// }

type S3service struct {
	constants *bootstrap.Constants
	buckets   map[enums.BucketType]*bootstrap.Bucket
	clients   map[enums.BucketType]*s3.S3
	uploader  map[enums.BucketType]*s3manager.Uploader
}

func NewS3Service(
	constants *bootstrap.Constants,
	bannerBucket *bootstrap.Bucket,
	sessionsBucket *bootstrap.Bucket,
	podcastsBucket *bootstrap.Bucket,

) *S3service {
	buckets := make(map[enums.BucketType]*bootstrap.Bucket)
	buckets[enums.BannersBucket] = bannerBucket
	buckets[enums.SessionsBucket] = sessionsBucket
	buckets[enums.PodcastsBucket] = podcastsBucket
	return &S3service{
		constants: constants,
		buckets:   buckets,
		clients:   make(map[enums.BucketType]*s3.S3),
		uploader:  make(map[enums.BucketType]*s3manager.Uploader),
	}
}

func (s3Service *S3service) setS3Client(bucketType enums.BucketType) {
	bucketTypes := enums.GetAllBucketTypes()
	if !slices.Contains(bucketTypes, bucketType) {
		panic(fmt.Errorf("bucket not exist"))
	}
	if s3Service.uploader[bucketType] != nil && s3Service.clients[bucketType] != nil {
		return
	}
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(s3Service.buckets[bucketType].AccessKey, s3Service.buckets[bucketType].SecretKey, ""),
		Region:      aws.String(s3Service.buckets[bucketType].Region),
		Endpoint:    aws.String(s3Service.buckets[bucketType].Endpoint),
	})

	if err != nil {
		panic(fmt.Errorf("unable to create AWS session, %v", err))
	}

	s3Service.uploader[bucketType] = s3manager.NewUploader(sess)
	s3Service.clients[bucketType] = s3.New(sess)
}

func (s3Service *S3service) UploadObject(bucketType enums.BucketType, key string, file *multipart.FileHeader) {
	s3Service.setS3Client(bucketType)
	bucket := s3Service.buckets[bucketType]

	fileReader, err := file.Open()
	if err != nil {
		panic(fmt.Errorf("unable to open file %q, %v", file.Filename, err))
	}
	defer fileReader.Close()

	_, err = s3Service.clients[bucketType].HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucket.Name),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && (aerr.Code() == s3.ErrCodeNoSuchBucket || aerr.Code() == "NotFound") {
			_, err = s3Service.clients[bucketType].CreateBucket(&s3.CreateBucketInput{
				Bucket: aws.String(bucket.Name),
			})
			if err != nil {
				panic(fmt.Errorf("unable to create bucket %q, %v", bucket.Name, err))
			}

			err = s3Service.clients[bucketType].WaitUntilBucketExists(&s3.HeadBucketInput{
				Bucket: aws.String(bucket.Name),
			})
			if err != nil {
				panic(fmt.Errorf("unable to confirm bucket %q exists, %v", bucket.Name, err))
			}
		} else {
			panic(fmt.Errorf("unable to check bucket %q, %v", bucket.Name, err))
		}
	}

	_, err = s3Service.uploader[bucketType].Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket.Name),
		Key:    aws.String(key),
		Body:   fileReader,
	})
	if err != nil {
		panic(fmt.Errorf("unable to upload %q to %q, %v", file.Filename, bucket.Name, err))
	}
}

func (s3Service *S3service) DeleteObject(bucketType enums.BucketType, objectName string) {
	s3Service.setS3Client(bucketType)
	bucket := s3Service.buckets[bucketType]

	_, err := s3Service.clients[bucketType].DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket.Name),
		Key:    aws.String(objectName),
	})
	if err != nil {
		panic(fmt.Errorf("unable to upload %q to %q, %v", objectName, bucket.Name, err))
	}

	err = s3Service.clients[bucketType].WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket.Name),
		Key:    aws.String(objectName),
	})
	if err != nil {
		panic(fmt.Errorf("unable to open file %q, %v", objectName, err))
	}
}

func (s3Service *S3service) GetPresignedURL(bucketType enums.BucketType, objectKey string, expiration time.Duration) string {
	s3Service.setS3Client(bucketType)
	bucket := s3Service.buckets[bucketType]

	req, _ := s3Service.clients[bucketType].GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket.Name),
		Key:    aws.String(objectKey),
	})

	url, err := req.Presign(expiration)
	if err != nil {
		panic(fmt.Errorf("failed to generate presigned URL: %w", err))
	}

	return url
}
