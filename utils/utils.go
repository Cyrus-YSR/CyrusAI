package utils

import (
	"GopherAI/model"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
)

func GetRandomNumbers(num int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := ""
	for i := 0; i < num; i++ {
		// 0~9随机�?
		digit := r.Intn(10)
		code += strconv.Itoa(digit)
	}
	return code
}

// MD5 MD5加密
func MD5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

func GenerateUUID() string {
	return uuid.New().String()
}

// �?schema 消息转换为数据库可存储的格式
func ConvertToModelMessage(sessionID string, userName string, msg *schema.Message) *model.Message {
	return &model.Message{
		SessionID: sessionID,
		UserName:  userName,
		Content:   msg.Content,
	}
}

// 将数据库消息转换�?schema 消息（供 AI 使用�?
func ConvertToSchemaMessages(msgs []*model.Message) []*schema.Message {
	schemaMsgs := make([]*schema.Message, 0, len(msgs))
	for _, m := range msgs {
		role := schema.Assistant
		if m.IsUser {
			role = schema.User
		}
		schemaMsgs = append(schemaMsgs, &schema.Message{
			Role:    role,
			Content: m.Content,
		})
	}
	return schemaMsgs
}

// RemoveAllFilesInDir 删除目录中的所有文件（不删除子目录�?
func RemoveAllFilesInDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 目录不存在就算了
		}
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := filepath.Join(dir, entry.Name())
			if err := os.Remove(filePath); err != nil {
				return err
			}
		}
	}
	return nil
}

// ValidateFile 校验文件是否为允许的文本文件�?md �?.txt�?
func ValidateFile(file *multipart.FileHeader) error {
	// 校验文件扩展�?
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".md" && ext != ".txt" {
		return fmt.Errorf("文件类型不正确，只允�?.md �?.txt 文件，当前扩展名: %s", ext)
	}

	return nil
}
