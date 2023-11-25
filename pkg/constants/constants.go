package constants

import "time"

const (
	JWTValue = "MTAxNTkwMTg1Mw=="

	// snowflake
	SnowflakeWorkerID     = 0
	SnowflakeDatacenterID = 0

	MaxConnections  = 1000
	MaxQPS          = 100
	MaxVideoSize    = 300000
	MaxListLength   = 100
	MaxIdleConns    = 10
	MaxGoroutines   = 10
	MaxOpenConns    = 100
	ConnMaxLifetime = 10 * time.Second

	PhoneRegexp  = "^1[345789]{1}\\d{9}$"
	CardIdRegexp = "(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)"
)
