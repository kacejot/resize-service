package utils

import (
	"fmt"
	"os"
)

// EnvOrDie returns environment variable or panics in case if
// environment variable is empty
func EnvOrDie(env string) string {
	result, err := Env(env)
	UnwrapWith(err, fmt.Sprintf("please set %s environment variable to proceed", env))

	return result
}

// Env returns environment variable or error in case if
// environment variable is empty
func Env(env string) (string, error) {
	result := os.Getenv(env)
	if "" == result {
		return result, fmt.Errorf("failed to load %s environment variable", env)
	}

	return result, nil
}

// EnvOr returns environment variable or default value in case if
// environment variable is empty
func EnvOr(env string, defaultEnv string) string {
	result := os.Getenv(env)
	if "" == result {
		fmt.Printf("environment variable %s is not found. Using %s as default value", env, defaultEnv)
		return defaultEnv
	}

	return result
}
