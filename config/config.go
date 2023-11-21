package config

import (
	"log"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

var (
	Server    *server
	Mysql     *mySQL
	Snowflake *snowflake
	//Etcd          *etcd
	//RabbitMQ      *rabbitMQ
	//Redis         *redis
	runtime_viper = viper.New()
)

func Init(path string, filename, service string) {

	_, err := ReadYamlConfig(path, filename)
	if err != nil {
		klog.Errorf("[config] error : %v", err.Error())
		panic(err)
	}
	configMapping(service)
}

func ReadYamlConfig(filePath, fileName string) (*viper.Viper, error) {
	return readFile(filePath, fileName, "yaml")
}

func readFile(filePath, fileName, configType string) (*viper.Viper, error) {
	runtime_viper.AddConfigPath(filePath)
	runtime_viper.SetConfigName(fileName)
	runtime_viper.SetConfigType(configType)
	if err := runtime_viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			klog.Error("Not found config file")
			return nil, err
		} else {
			klog.Errorf("read config file error, %v \n", err)
			return nil, err
		}
	}
	return runtime_viper, nil
}

func configMapping(srv string) {
	c := new(config)
	if err := runtime_viper.Unmarshal(&c); err != nil {
		log.Fatal(err)
	}
	Snowflake = &c.Snowflake

	Server = &c.Server
	Server.Secret = []byte(runtime_viper.GetString("server.jwt-secret"))

	Mysql = &c.MySQL
	// RabbitMQ = &c.RabbitMQ
	// Redis = &c.Redis
	// OSS = &c.OSS
	// Jaeger = &c.Jaeger
	// USS = &c.USS
	// Elasticsearch = &c.Elasticsearch
	// Service = GetService(srv)
}
