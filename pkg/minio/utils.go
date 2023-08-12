package minio

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

// UploadLocalFile 将本地文件上传至minio
func UploadLocalFile(filePath string, bucketName string, fileName string, opt minio.PutObjectOptions) error {
	// check whether the bucket exists
	if exists, err := Client.BucketExists(context.Background(), bucketName); !(err == nil && exists) {
		return fmt.Errorf("minio buckect %s miss,err: %v\n", bucketName, err)
	}
	uploadInfo, err := Client.FPutObject(
		context.Background(), bucketName, fileName, filePath, opt)
	if err != nil {
		return fmt.Errorf("failed uploaded object, err : %w", err)
	}
	fmt.Println("Successfully uploaded object: ", uploadInfo)
	return nil
}

// UploadSizeFile 读取固定大小的文件并上传至minio(主要使用)
func UploadSizeFile(bucketName string, fileName string, reader io.Reader, size int64, opt minio.PutObjectOptions) error {
	if exists, err := Client.BucketExists(context.Background(), bucketName); !(err == nil && exists) {
		return fmt.Errorf("minio buckect %s miss,err: %v\n", bucketName, err)
	}
	uploadInfo, err := Client.PutObject(context.Background(), bucketName, fileName, reader, size, opt)
	if err != nil {
		return fmt.Errorf("failed uploaded bytes, err : %w", err)
	}
	fmt.Println("Successfully uploaded bytes: ", uploadInfo)
	return nil
}

// GetFileURL 根据文件名从minio获取文件URL
func GetFileURL(bucketName string, fileName string, timeLimit time.Duration) (*url.URL, error) {
	if exists, err := Client.BucketExists(context.Background(), bucketName); !(err == nil && exists) {
		return nil, fmt.Errorf("minio buckect %s miss,err: %v\n", bucketName, err)
	}
	reqParams := make(url.Values)
	reqParams.Set("response-content", "attachment; filename=\""+fileName+"\"")

	presignedURL, err := Client.PresignedGetObject(context.Background(), bucketName, fileName, timeLimit, reqParams)
	if err != nil {
		return nil, fmt.Errorf("failed generated presigned URL, err : %w", err)
	}
	fmt.Println("Successfully generated presigned URL", presignedURL)
	return presignedURL, nil
}
