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
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"gopkg.in/yaml.v3"
)

const (
	configPath = "../config/minio.yaml"
)

var (
	Client *minio.Client
	conf   Config
)

type Config struct {
	Endpoint     string `yaml:"endpoint"`
	AccessKeyId  string `yaml:"accessKeyId"`
	AccessSecret string `yaml:"accessSecret"`
	UseSSL       bool   `yaml:"useSSL"`
	BucketName   string `yaml:"bucketName" `
}

func readConfig(conf *Config) {
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading configuration file, err: %v", err)
	}
	if err = yaml.Unmarshal(yamlFile, conf); err != nil {
		log.Fatalf("Error reading configuration file, err: %v", err)
	}
}

// Minio Object Storage Initial
func init() {
	readConfig(&conf)
	var err error
	Client, err = minio.New(conf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKeyId, conf.AccessSecret, ""),
		Secure: conf.UseSSL,
	})
	if err != nil {
		log.Fatalf("minio client init failed,err: %v", err)
	}
	log.Println("minio client init successfully")
}
