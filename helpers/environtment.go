package helpers

import "os"

func Env(key string, defaultValue string) (value string) {
	value, found := os.LookupEnv(key)
	if !found {
		value = defaultValue
	}

	return value
}
