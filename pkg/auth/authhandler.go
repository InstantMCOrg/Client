package auth

import "os"

const key = "auth"

func HasAuthKey() bool {
	_, ok := os.LookupEnv(key)
	return ok
}

func GetAuthKey() string {
	return os.Getenv(key)
}
