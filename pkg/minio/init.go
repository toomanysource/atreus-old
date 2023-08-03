/*
 * @Author: DyamidSteve 2857199455@qq.com
 * @Date: 2023-08-02 19:10:59
 * @LastEditors: DyamidSteve 2857199455@qq.com
 * @LastEditTime: 2023-08-02 19:10:59
 * @FilePath: /Atreus/pkg/minio/init.go
 * @Description: Minio Object Storage Initialization
 */

package minio

import (
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v6"
	"gopkg.in/yaml.v3"
)

const (
	//configFilePath = "../config"
	configPath = "../config/minio.yaml"
)

var (
	minioClient *minio.Client
	cfg         MinioConfig
)

type MinioConfig struct {
	// just modify in minio.yaml
	Endpoint        string `yaml:"endpoint" json:"endpoint"`
	AccessKeyId     string `yaml:"acesskeyid" json:"acesskeyid"`
	SecretAccessKey string `yaml:"secret" json:"secret"`
	UseSSL          bool   `yaml:"usessl" json:"usessl"`

	// you can use the default config below to UploadFile.
	BucketName  string `yaml:"bucketname" json:"bucketname"`
	Location    string `yaml:"location" json:"location"`
	ContentType string `yaml:"contentType" json:"contentType"`
}

func readConfigYaml(cfg *MinioConfig) {
	// load yaml config
	yamlFile, err := os.ReadFile(configPath) //ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println(err.Error())
	}
	// load the yaml into a structure
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// Minio Object Storage Initial
func init() {
	// read congif from ../config/minio.yaml
	readConfigYaml(&cfg)
	//fmt.Println(cfg)
	client, err := minio.New(cfg.Endpoint, cfg.AccessKeyId, cfg.SecretAccessKey, cfg.UseSSL)
	if err != nil {
		log.Fatalln("minio client init failed,err:", err)
	}
	log.Println("minio client init successfully")
	minioClient = client

	// check whether the bucket exists
	if exists, err := minioClient.BucketExists(cfg.BucketName); !(err == nil && exists) {
		log.Fatalf("minio buckect %s miss,err: %v\n", cfg.BucketName, err)
	}
}
