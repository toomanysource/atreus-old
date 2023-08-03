package minio

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestUploadLocalFileAndGetURL(t *testing.T) {
	fmt.Println("Start Test TestUploadLocalFileAndGetURL")
	// 1.upload video
	if err := UploadLocalFile("test1.mp4", "./test1.mp4", "mybucket", "video/mp4"); err != nil {
		t.Fatal("failed to upload file test1.mp4")
	}
	// or use default method
	if err := UploadLocalFileWithDefaultConfig("test2.mp4", "./test2.mp4"); err != nil {
		t.Fatal("failed to upload file test2.mp4")
	}

	var expireTime = time.Second * 10
	// get target file's url
	if url1, err := GetFileURL("mybucket", "test1.mp4", expireTime); err != nil {
		t.Fatal("failed to get url from file test1.mp4")
	} else {
		_url1 := strings.Split(url1.String(), "?")[0]
		// if _url1 != "http://{baseURL}/{bucketName}/test1.mp4" {
		// 	t.Fatalf("url error,expected:%s,get:%s", "test1.mp4", _url1)
		// }
		fmt.Println(_url1)
	}

	if url2, err := GetFileURL("mybucket", "test2.mp4", expireTime); err != nil {
		t.Fatal("failed to get url from file test2.mp4")
	} else {
		_url2 := strings.Split(url2.String(), "?")[0]
		// if _url2 != "http://{baseURL}/{bucketName}/test2.mp4" {
		// 	t.Fatalf("url error,expected:%s,get:%s", "test2.mp4", _url2)
		// }
		fmt.Println(_url2)
	}

	// 2.upload image
	if err := UploadLocalFile("test1.png", "./test1.png", "mybucket", "image/png"); err != nil {
		t.Fatal("failed to upload file test1.png")
	}

	// if file type is image,can't use default type

	// get target file's url
	if url1, err := GetFileURL("mybucket", "test1.png", expireTime); err != nil {
		t.Fatal("failed to get url from file test1.png")
	} else {
		_url1 := strings.Split(url1.String(), "?")[0]
		// if _url1 != "http://{baseURL}/{bucketName}/test1.png" {
		// 	t.Fatalf("url error,expected:%s,get:%s", "test1.png", _url1)
		// }
		fmt.Println(_url1)
	}

}

func TestUploadFileAndGetURL(t *testing.T) {
	fmt.Println("Start Test TestUploadFileAndGetURL")
	// upload video using default method
	file1, _ := os.Open("test1.mp4")
	defer file1.Close()
	if err := UploadFileWithDefaultConfig("test2.mp4", file1, -1); err != nil {
		t.Fatal("failed to upload file test1.mp4")
	}

	// upload image
	file2, _ := os.Open("test1.png")
	defer file1.Close()
	if err := UploadFile("mybucket", "test1.png", file2, -1, "image/png"); err != nil {
		t.Fatal("failed to upload file test1.png")
	}

	var expireTime = time.Second * 10
	// get target file's url
	if url1, err := GetFileURL("mybucket", "test1.mp4", expireTime); err != nil {
		t.Fatal("failed to get url from file test1.mp4")
	} else {
		_url1 := strings.Split(url1.String(), "?")[0]
		// if _url1 != "http://{baseURL}/{bucketName}/test1.mp4" {
		// 	t.Fatalf("url error,expected:%s,get:%s", "test1.mp4", _url1)
		// }
		fmt.Println(_url1)
	}

	if url2, err := GetFileURL("mybucket", "test1.png", expireTime); err != nil {
		t.Fatal("failed to get url from file test2.mp4")
	} else {
		_url2 := strings.Split(url2.String(), "?")[0]
		// if _url2 != "http://{baseURL}/{bucketName}/test2.mp4" {
		// 	t.Fatalf("url error,expected:%s,get:%s", "test2.mp4", _url2)
		// }
		fmt.Println(_url2)
	}
}

func BenchmarkUploadFile(b *testing.B) {
	fmt.Println("BenchmarkUploadFile")
	// upload video using default method
	file1, _ := os.Open("test1.mp4")
	defer file1.Close()

	for i := 0; i < 100; i++ {
		if err := UploadFileWithDefaultConfig("test1.mp4", file1, -1); err != nil {
			b.Fatal("failed to upload file test1.mp4")
		}
	}
}

func BenchmarkGetFileURL(b *testing.B) {
	var expireTime = time.Second * 10
	// get target file's url
	for i := 0; i < 1000; i++ {
		if url, err := GetFileURL("mybucket", "test1.mp4", expireTime); err != nil {
			b.Fatal("failed to get url from file test1.mp4")
		} else {
			_url := strings.Split(url.String(), "?")[0]
			// if _url != "http://{baseURL}/{bucketName}/test1.mp4" {
			// 	t.Fatalf("url error,expected:%s,get:%s", "test1.mp4", _url)
			// }
			fmt.Println("get", _url, "count:", i)
		}
	}
}
