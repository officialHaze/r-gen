package util

import "github.com/google/uuid"

func UUIDGen() string {
	return uuid.New().String()
}
