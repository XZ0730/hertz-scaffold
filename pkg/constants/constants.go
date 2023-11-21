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
)
