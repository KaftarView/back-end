package application_aws

import (
	"first-project/src/bootstrap"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AWSS3 struct {
	constants *bootstrap.Constants
	bucket    *bootstrap.Bucket
	session   *session.Session
	s3Client  *s3.S3
}

func NewAWSS3(constants *bootstrap.Constants, bucket *bootstrap.Bucket) *AWSS3 {
	return &AWSS3{
		constants: constants,
		bucket:    bucket,
	}
}

func (a *AWSS3) getS3Client() {
	if a.session != nil && a.s3Client != nil {
		return
	}
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(a.bucket.AccessKey, a.bucket.SecretKey, ""),
		Region:      aws.String(a.bucket.Region),
		Endpoint:    aws.String(a.bucket.Endpoint),
	})

	if err != nil {
		panic(fmt.Errorf("unable to create AWS session, %v", err))
	}

	a.session = sess
	a.s3Client = s3.New(sess)
}

func (a *AWSS3) UploadObject(file *multipart.FileHeader, objectTittle string, objectID int) {
	a.getS3Client()
	fileReader, err := file.Open()
	if err != nil {
		panic(fmt.Errorf("unable to open file %q, %v", file.Filename, err))
	}
	defer fileReader.Close()

	_, err = a.s3Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(a.bucket.Name),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && (aerr.Code() == s3.ErrCodeNoSuchBucket || aerr.Code() == "NotFound") {
			_, err = a.s3Client.CreateBucket(&s3.CreateBucketInput{
				Bucket: aws.String(a.bucket.Name),
			})
			if err != nil {
				panic(fmt.Errorf("unable to create bucket %q, %v", a.bucket.Name, err))
			}

			err = a.s3Client.WaitUntilBucketExists(&s3.HeadBucketInput{
				Bucket: aws.String(a.bucket.Name),
			})
			if err != nil {
				panic(fmt.Errorf("unable to confirm bucket %q exists, %v", a.bucket.Name, err))
			}
		} else {
			panic(fmt.Errorf("unable to check bucket %q, %v", a.bucket.Name, err))
		}
	}

	uploader := s3manager.NewUploader(a.session)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(a.bucket.Name),
		Key:    aws.String(a.constants.ObjectStorage.GetObjectKey(int(objectID), objectTittle, file.Filename)),
		Body:   fileReader,
	})
	if err != nil {
		panic(fmt.Errorf("unable to upload %q to %q, %v", file.Filename, a.bucket.Name, err))
	}
}

func (a *AWSS3) DeleteObject(objectName string) {
	a.getS3Client()
	_, err := a.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(a.bucket.Name),
		Key:    aws.String(objectName),
	})
	if err != nil {
		panic(fmt.Errorf("unable to upload %q to %q, %v", objectName, a.bucket.Name, err))
	}

	err = a.s3Client.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(a.bucket.Name),
		Key:    aws.String(objectName),
	})
	if err != nil {
		panic(fmt.Errorf("unable to open file %q, %v", objectName, err))
	}
}

func (a *AWSS3) ListObjects() []map[string]interface{} {
	a.getS3Client()
	resp, err := a.s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(a.bucket.Name),
	})
	if err != nil {
		panic(fmt.Errorf("unable to list items in bucket %q, %v", a.bucket.Name, err))
	}

	itemMap := make([]map[string]interface{}, 0)
	for _, item := range resp.Contents {
		itemData := map[string]interface{}{
			"Name":         *item.Key,
			"LastModified": *item.LastModified,
			"Size":         *item.Size,
			"StorageClass": *item.StorageClass,
		}
		itemMap = append(itemMap, itemData)
	}

	return itemMap
}

func (a *AWSS3) listSessionVideos(sessionID uint) []string {
	var videoKeys []string
	prefix := fmt.Sprintf("session/%d/", int(sessionID))
	err := a.s3Client.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket: aws.String(a.bucket.Name),
		Prefix: aws.String(prefix),
	}, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			videoKeys = append(videoKeys, *obj.Key)
		}
		return true
	})

	if err != nil {
		panic(fmt.Errorf("unable to list videos for session %s: %v", prefix, err))
	}
	return videoKeys
}

func (a *AWSS3) GetSessionVideoURLs(sessionID uint) map[string]string {
	a.getS3Client()
	videoKeys := a.listSessionVideos(sessionID)
	videoURLs := make(map[string]string)
	for _, key := range videoKeys {
		req, _ := a.s3Client.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(a.bucket.Name),
			Key:    aws.String(key),
		})

		urlStr, err := req.Presign(24 * time.Hour)
		if err != nil {
			panic(fmt.Errorf("failed to sign request for %s: %v", key, err))
		}
		videoURLs[key] = urlStr
	}
	return videoURLs
}
