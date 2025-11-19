package pkg

import (
	"os"
)

func GetEnvDefault(key, defVal string) string {
	val, err := os.LookupEnv(key)
	if !err {
		return defVal
	}
	return val
}
