package utils

import uuid "github.com/satori/go.uuid"

func GenerateUUID() string {
	//生成uuid
	uuid := uuid.NewV4()
	return uuid.String()
}
