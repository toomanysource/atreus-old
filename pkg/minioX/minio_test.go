package minioX

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	endpoint    = "localhost:9000"
	acessKeyId  = "acessKeyId"
	secret      = "secret"
	useSSL      = false
	bucketName  = "mybucket"
	location    = "beijing"
	contentType = "video/mp4"
)

func NewMinioClient() *minio.Client {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(acessKeyId, secret, ""),
		Secure: useSSL,
		//Region: location,
	})

	if err != nil {
		log.Fatalf("minio client init failed,err: %v", err)
	}
	log.Info("minioExtra enabled successfully")
	return client
}

func TestExistBucket(t *testing.T) {
	minioClient := NewMinioClient()

	client := NewClient(NewExtraConn(minioClient), NewIntraConn(minioClient))

	if exist, err := client.ExistBucket(context.Background(), bucketName); !exist || err != nil {
		t.Fatal("bucket miss,err:", err)
	}
}

func TestUploadFileAndGetURL(t *testing.T) {
	minioClient := NewMinioClient()

	client := NewClient(NewExtraConn(minioClient), NewIntraConn(minioClient))

	f, err := os.Open("test1.mp4")
	if err != nil {
		t.Fatalf("open file test1.mp4 error:%v", err)
	}
	defer f.Close()
	// 获取文件统计信息
	fileInfo, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	sizeInBytes := fileInfo.Size()
	// 上传整个文件
	if err := client.UploadSizeFile(context.Background(), bucketName, "test1.mp4", f, sizeInBytes, minio.PutObjectOptions{
		ContentType: contentType,
	}); err != nil {
		t.Fatal(err)
	}
	// 获取文件URL
	url, err := client.GetFileURL(context.Background(), bucketName, "test1.mp4", time.Second)
	if err != nil {
		t.Fatal(err)
	}

	_url := strings.Split(url.String(), "?")[0]

	fmt.Println("get file url:", _url)

}

func BenchmarkUploadFile(b *testing.B) {
	minioClient := NewMinioClient()

	client := NewClient(NewExtraConn(minioClient), NewIntraConn(minioClient))

	f, err := os.Open("test1.mp4")
	if err != nil {
		b.Fatalf("open file test1.mp4 error:%v", err)
	}
	defer f.Close()
	// 获取文件统计信息
	fileInfo, err := f.Stat()
	if err != nil {
		b.Fatal(err)
	}
	sizeInBytes := fileInfo.Size()

	for i := 0; i < 10; i++ {
		if err := client.UploadSizeFile(context.Background(), bucketName, "test1.mp4", f, sizeInBytes, minio.PutObjectOptions{
			ContentType: contentType,
		}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetFileURL(b *testing.B) {
	minioClient := NewMinioClient()

	client := NewClient(NewExtraConn(minioClient), NewIntraConn(minioClient))

	for i := 0; i < 1000; i++ {
		_, err := client.GetFileURL(context.Background(), bucketName, "test1.mp4", time.Second)
		if err != nil {
			b.Fatal(err)
		}
		fmt.Println("Get time", i)
	}
}
