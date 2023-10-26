package config

import "os"

const (
	IsDev          = "DEBUG"
	LogLevel       = "LOG_LEVEL"
	DBHost         = "DB_HOST"
	DBReadPort     = "DB_PORT"
	DBDatabase     = "DB_DATABASE"
	DBReadUsername = "DB_USERNAME"
	DBReadPassword = "DB_PASSWORD"
	RedisAddr      = "REDIS_ADDR"
	RedisPassword  = "REDIS_PASSWORD"
	GatewayAddress = "GATEWAY_ADDRESS"
	GrpcAddress    = "GRPC_ADDRESS"
	SessionLifeDay = "SESSION_LIFE_DAY"
	CookieName     = "COOKIE_NAME"
)

var defaults = map[string]string{
	IsDev:          "true",
	LogLevel:       "0",
	DBHost:         "127.0.0.1",
	DBDatabase:     "quick_web",
	DBReadUsername: "root",
	DBReadPassword: "123456",
	DBReadPort:     "3306",
	RedisAddr:      "127.0.0.1:6379",
	RedisPassword:  "123456",
	GatewayAddress: "0.0.0.0:3000",
	GrpcAddress:    "0.0.0.0:8000",
	SessionLifeDay: "1",
	CookieName:     "quick_web",
}

func Get(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	defValue, ok := defaults[key]
	if !ok {
		return ""
	}

	return defValue
}
