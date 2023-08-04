package minio

import (
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v6"
)

// UploadLocalFile will upload the local file to minio based on the (fileName,filePath,bucketName,contentType) you provided.
// bucketName in minio is similar to a folder in op system.
// contentType means the file type(such as application/json,video/mp4).
func UploadLocalFile(fileName string, filePath string, bucketName string, contentType string) error {
	// use the minio.FPutObject to upload File
	_, err := minioClient.FPutObject(bucketName, fileName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return fmt.Errorf("upload file failed,err: %v", err)
	}

	//log.Printf("success")
	return nil
}

// UploadLocalFileWithDefaultConfig will upload the local file to minio based on the (fileName,filePath) and on the default (bucketName,contentType) in Config.
func UploadLocalFileWithDefaultConfig(fileName string, filePath string) error {
	// use the minio.FPutObject to upload File
	_, err := minioClient.FPutObject(cfg.BucketName, fileName, filePath, minio.PutObjectOptions{ContentType: cfg.ContentType})
	if err != nil {
		return fmt.Errorf("upload file failed,err: %v", err)
	}

	//log.Printf("success")
	return nil
}

// UploadFile will upload file to minio based on the reader you provided,further more,you should specify (fileName,bucketName,contentType,objectsize).
// bucketName in minio is similar to a folder in op system.
// contentType means the file type(such as application/json,video/mp4)
// objectsize is the byte size of the object file(if less than 0,then the method will upload file until read EOF).
func UploadFile(bucketName string, fileName string, reader io.Reader, objectsize int64, contentType string) error {
	// use the minio.PutObject to upload File
	if objectsize < int64(-1) {
		objectsize = -1
	}
	_, err := minioClient.PutObject(bucketName, fileName, reader, objectsize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return fmt.Errorf("upload %s of size %d failed, %s", bucketName, objectsize, err)
	}
	return nil
}

// UploadFile will upload file to minio based on the (fileName,reader,objectsize) you provided and on the default (bucketName,contentType) in Config.
// if objectsize less than 0,the method will upload file until read EOF.
func UploadFileWithDefaultConfig(fileName string, reader io.Reader, objectsize int64) error {
	// use the minio.PutObject to upload File
	if objectsize < int64(-1) {
		objectsize = -1
	}
	_, err := minioClient.PutObject(cfg.BucketName, fileName, reader, objectsize, minio.PutObjectOptions{ContentType: cfg.ContentType})
	if err != nil {
		return fmt.Errorf("upload %s of size %d failed, %s", cfg.BucketName, objectsize, err)
	}
	return nil
}

// GetFileUrl will get URL from minio base on the fileName in the given bucket.
// bucketName in minio is similar to a folder in op system.
// when it expire the duration of timeLimit,it will return an TimeExpire error.
// return the file's net URL and error
func GetFileURL(bucketName string, fileName string, timeLimit time.Duration) (*url.URL, error) {
	reqParams := make(url.Values)
	if timeLimit <= 0 {
		timeLimit = time.Minute
	}
	presignedUrl, err := minioClient.PresignedGetObject(bucketName, fileName, timeLimit, reqParams)
	if err != nil {
		return nil, fmt.Errorf("get url of file %s from bucket %s failed, %s", fileName, bucketName, err)
	}
	// TODO: url need to be split from url.String(),example: _url := strings.Split(url.String(), "?")[0]
	return presignedUrl, nil
}

// GetFileUrl will get URL from minio base on the fileName in default bucket from Config.
// when it expire the duration of timeLimit,it will return an TimeExpire error.
// return the file's net URL and error
func GetFileURLWithDefaultBucket(fileName string, timeLimit time.Duration) (*url.URL, error) {
	reqParams := make(url.Values)
	if timeLimit <= 0 {
		timeLimit = time.Minute
	}
	presignedUrl, err := minioClient.PresignedGetObject(cfg.BucketName, fileName, timeLimit, reqParams)
	if err != nil {
		return nil, fmt.Errorf("get url of file %s from bucket %s failed, %s", fileName, cfg.BucketName, err)
	}
	// TODO: url need to be split from url.String(),example: _url := strings.Split(url.String(), "?")[0]
	return presignedUrl, nil
}
