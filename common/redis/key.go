package redis

import (
	"GopherAI/config"
	"fmt"
)

// key:特定邮箱-> 验证�?
func GenerateCaptcha(email string) string {
	return fmt.Sprintf(config.DefaultRedisKeyConfig.CaptchaPrefix, email)
}

func GenerateIndexName(filename string) string {
	indexName := fmt.Sprintf(config.DefaultRedisKeyConfig.IndexName, filename)
	return indexName
}

func GenerateIndexNamePrefix(filename string) string {
	prefix := fmt.Sprintf(config.DefaultRedisKeyConfig.IndexNamePrefix, filename)
	return prefix
}
