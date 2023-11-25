package config

type server struct {
	Secret  []byte
	Version string
	Name    string
}

type snowflake struct {
	WorkerID      int64 `mapstructure:"worker-id"`
	DatancenterID int64 `mapstructure:"datancenter-id"`
}

type mySQL struct {
	Addr     string
	Database string
	Username string
	Password string
	Charset  string
}

// type rabbitMQ struct {
// 	Addr     string
// 	Username string
// 	Password string
// }

type redis struct {
	Addr     string
	Password string
}

type qss struct {
	AccessKey   string `mapstructure:"access-key"`
	SerectKey   string `mapstructure:"serect-key"`
	QiniuServer string `mapstructure:"qiniu-server"`
	Bucket      string
}

type config struct {
	Server    server
	Snowflake snowflake
	MySQL     mySQL
	Redis     redis
	QSS       qss
	// Jaeger        jaeger
	// Etcd          etcd
	// RabbitMQ      rabbitMQ
	// OSS           oss
	// USS           uss
	// Elasticsearch elasticsearch
}
