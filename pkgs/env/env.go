package pkgs_env

import (
	"os"
	"strconv"
)

func GetEnvString(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
func GetEnvInt64(key string, defaultVal int64) int64 {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	res, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defaultVal
	}
	return res
}
