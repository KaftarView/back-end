package application_interfaces

import (
	"first-project/src/enums"
	"mime/multipart"
	"time"
)

type S3Service interface {
	DeleteObject(bucketType enums.BucketType, key string) error
	GetPresignedURL(bucketType enums.BucketType, objectKey string, expiration time.Duration) string
	UploadObject(bucketType enums.BucketType, key string, file *multipart.FileHeader)
}
