package session

import (
	"GopherAI/common/aihelper"
	"GopherAI/common/code"
	"GopherAI/dao/message"
	"GopherAI/dao/session"
	"GopherAI/model"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

var ctx = context.Background()

func GetUserSessionsByUserName(userName string) ([]model.SessionInfo, error) {
	//获取用户的所有会话ID

	manager := aihelper.GetGlobalManager()
	sessionIDs := manager.GetUserSessions(userName)

	if len(sessionIDs) == 0 {
		return []model.SessionInfo{}, nil
	}

	dbSessions, err := session.GetSessionsByIDs(sessionIDs)
	if err != nil {
		log.Println("GetUserSessionsByUserName GetSessionsByIDs error:", err)

		fallback := make([]model.SessionInfo, 0, len(sessionIDs))
		for _, sid := range sessionIDs {
			fallback = append(fallback, model.SessionInfo{
				SessionID: sid,
				Title:     sid,
			})
		}
		return fallback, nil
	}

	titleByID := make(map[string]string, len(dbSessions))
	for _, s := range dbSessions {
		titleByID[s.ID] = s.Title
	}

	result := make([]model.SessionInfo, 0, len(sessionIDs))
	for _, sid := range sessionIDs {
		title := titleByID[sid]
		if title == "" {
			title = sid
		}
		result = append(result, model.SessionInfo{
			SessionID: sid,
			Title:     title,
		})
	}

	return result, nil
}

func CreateSessionAndSendMessage(userName string, userQuestion string, modelType string) (string, string, code.Code) {
	//1：创建一个新的会�?
	newSession := &model.Session{
		ID:       uuid.New().String(),
		UserName: userName,
		Title:    userQuestion, // 可以根据需求设置标题，这边暂时用用户第一次的问题作为标题
	}
	createdSession, err := session.CreateSession(newSession)
	if err != nil {
		log.Println("CreateSessionAndSendMessage CreateSession error:", err)
		return "", "", code.CodeServerBusy
	}

	//2：获取AIHelper并通过其管理消�?
	manager := aihelper.GetGlobalManager()
	config := map[string]interface{}{
		"apiKey":   "your-api-key", // TODO: 从配置中获取
		"username": userName,       // 用于 RAG 模型获取用户文档
	}
	helper, err := manager.GetOrCreateAIHelper(userName, createdSession.ID, modelType, config)
	if err != nil {
		log.Println("CreateSessionAndSendMessage GetOrCreateAIHelper error:", err)
		return "", "", code.AIModelFail
	}

	//3：生成AI回复
	aiResponse, err_ := helper.GenerateResponse(userName, ctx, userQuestion)
	if err_ != nil {
		log.Println("CreateSessionAndSendMessage GenerateResponse error:", err_)
		return "", "", code.AIModelFail
	}

	return createdSession.ID, aiResponse.Content, code.CodeSuccess
}

func CreateStreamSessionOnly(userName string, userQuestion string) (string, code.Code) {
	newSession := &model.Session{
		ID:       uuid.New().String(),
		UserName: userName,
		Title:    userQuestion,
	}
	createdSession, err := session.CreateSession(newSession)
	if err != nil {
		log.Println("CreateStreamSessionOnly CreateSession error:", err)
		return "", code.CodeServerBusy
	}
	return createdSession.ID, code.CodeSuccess
}

func StreamMessageToExistingSession(userName string, sessionID string, userQuestion string, modelType string, writer http.ResponseWriter) code.Code {
	log.Printf("[Service] StreamMessageToExistingSession Start. User=%s, Session=%s, Model=%s", userName, sessionID, modelType)
	// 确保 writer 支持 Flush
	flusher, ok := writer.(http.Flusher)
	if !ok {
		log.Println("StreamMessageToExistingSession: streaming unsupported")
		return code.CodeServerBusy
	}

	manager := aihelper.GetGlobalManager()
	config := map[string]interface{}{
		"apiKey":   "your-api-key", // TODO: 从配置中获取
		"username": userName,       // 用于 RAG 模型获取用户文档
	}
	log.Println("[Service] Getting AIHelper...")
	helper, err := manager.GetOrCreateAIHelper(userName, sessionID, modelType, config)
	if err != nil {
		log.Println("StreamMessageToExistingSession GetOrCreateAIHelper error:", err)
		return code.AIModelFail
	}
	log.Println("[Service] AIHelper Obtained. Starting StreamResponse...")

	cb := func(msg string) {
		log.Printf("[SSE] Sending chunk: %s (len=%d)\n", msg, len(msg))
		payload, err := json.Marshal(map[string]string{
			"type":    "delta",
			"content": msg,
		})
		if err != nil {
			log.Println("[SSE] Marshal error:", err)
			return
		}
		_, err = writer.Write([]byte("data: " + string(payload) + "\n\n"))
		if err != nil {
			log.Println("[SSE] Write error:", err)
			return
		}
		flusher.Flush()
		log.Println("[SSE] Flushed")
	}

	_, err_ := helper.StreamResponse(userName, ctx, cb, userQuestion)
	if err_ != nil {
		log.Println("StreamMessageToExistingSession StreamResponse error:", err_)
		return code.AIModelFail
	}

	_, err = writer.Write([]byte("data: [DONE]\n\n"))
	if err != nil {
		log.Println("StreamMessageToExistingSession write DONE error:", err)
		return code.AIModelFail
	}
	flusher.Flush()

	return code.CodeSuccess
}

func CreateStreamSessionAndSendMessage(userName string, userQuestion string, modelType string, writer http.ResponseWriter) (string, code.Code) {

	sessionID, code_ := CreateStreamSessionOnly(userName, userQuestion)
	if code_ != code.CodeSuccess {
		return "", code_
	}

	code_ = StreamMessageToExistingSession(userName, sessionID, userQuestion, modelType, writer)
	if code_ != code.CodeSuccess {

		return sessionID, code_
	}

	return sessionID, code.CodeSuccess
}

func ChatSend(userName string, sessionID string, userQuestion string, modelType string) (string, code.Code) {
	//1：获取AIHelper
	manager := aihelper.GetGlobalManager()
	config := map[string]interface{}{
		"username": userName, // 用于 RAG 模型获取用户文档（若当前用户选择了RAG模型，该字段将会被用到）
	}
	helper, err := manager.GetOrCreateAIHelper(userName, sessionID, modelType, config)
	if err != nil {
		log.Println("ChatSend GetOrCreateAIHelper error:", err)
		return "", code.AIModelFail
	}

	//2：生成AI回复
	aiResponse, err_ := helper.GenerateResponse(userName, ctx, userQuestion)
	if err_ != nil {
		log.Println("ChatSend GenerateResponse error:", err_)
		return "", code.AIModelFail
	}

	return aiResponse.Content, code.CodeSuccess
}

func GetChatHistory(userName string, sessionID string) ([]model.History, code.Code) {
	// 获取AIHelper中的消息历史
	manager := aihelper.GetGlobalManager()
	helper, exists := manager.GetAIHelper(userName, sessionID)
	if !exists {
		return nil, code.CodeServerBusy
	}

	messages := helper.GetMessages()
	history := make([]model.History, 0, len(messages))

	// 转换消息为历史格式（根据消息顺序或内容判断用�?AI消息�?
	for i, msg := range messages {
		isUser := i%2 == 0
		history = append(history, model.History{
			IsUser:  isUser,
			Content: msg.Content,
		})
	}

	return history, code.CodeSuccess
}

func ChatStreamSend(userName string, sessionID string, userQuestion string, modelType string, writer http.ResponseWriter) code.Code {

	return StreamMessageToExistingSession(userName, sessionID, userQuestion, modelType, writer)
}

func DeleteSession(userName string, sessionID string) code.Code {
	sess, err := session.GetSessionByID(sessionID)
	if err != nil {
		log.Println("DeleteSession GetSessionByID error:", err)
		return code.CodeServerBusy
	}
	if sess.UserName != userName {
		return code.CodeInvalidParams
	}
	if err := message.DeleteMessagesBySessionID(sessionID); err != nil {
		log.Println("DeleteSession DeleteMessagesBySessionID error:", err)
		return code.CodeServerBusy
	}
	if err := session.DeleteSessionByID(sessionID); err != nil {
		log.Println("DeleteSession DeleteSessionByID error:", err)
		return code.CodeServerBusy
	}
	manager := aihelper.GetGlobalManager()
	manager.RemoveAIHelper(userName, sessionID)
	return code.CodeSuccess
}
