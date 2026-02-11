package image

import (
	"GopherAI/common/aihelper"
	"GopherAI/common/code"
	"GopherAI/controller"
	"GopherAI/service/image"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
)

type (
	RecognizeImageResponse struct {
		ClassName string `json:"class_name,omitempty"` // AI回答
		controller.Response
	}
	AnalyzeImageResponse struct {
		ClassName    string                   `json:"class_name,omitempty"`
		TopK         []map[string]interface{} `json:"top_k,omitempty"`
		AnalysisText string                   `json:"analysis_text,omitempty"`
		controller.Response
	}
)

func RecognizeImage(c *gin.Context) {
	res := new(RecognizeImageResponse)
	file, err := c.FormFile("image")
	if err != nil {
		log.Println("FormFile fail ", err)
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	className, err := image.RecognizeImage(file)
	if err != nil {
		log.Println("RecognizeImage fail ", err)
		c.JSON(http.StatusOK, res.CodeOf(code.CodeServerBusy))
		return
	}

	res.Success()
	res.ClassName = className
	c.JSON(http.StatusOK, res)
}

func AnalyzeImage(c *gin.Context) {
	res := new(AnalyzeImageResponse)
	file, err := c.FormFile("image")
	if err != nil {
		log.Println("FormFile fail ", err)
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}
	className, err := image.RecognizeImage(file)
	if err != nil {
		log.Println("RecognizeImage fail ", err)
		c.JSON(http.StatusOK, res.CodeOf(code.CodeServerBusy))
		return
	}
	res.Success()
	res.ClassName = className
	res.TopK = []map[string]interface{}{
		{"label": className, "score": 1.0},
	}
	username := c.GetString("userName")
	ctx := context.Background()
	factory := aihelper.GetGlobalFactory()
	// 使用基础 OpenAI 模型 (Type 1) 而非 MCP 模型 (Type 3)，避免工具调用提示词干扰生成结果
	model, err := factory.CreateAIModel(ctx, "1", map[string]interface{}{"username": username})
	if err != nil {
		log.Println("CreateAIModel fail ", err)
		c.JSON(http.StatusOK, res)
		return
	}
	var b strings.Builder
	b.WriteString("你是一位专业的图像讲解助手，请基于识别结果做详细中文解析。\n")
	b.WriteString("识别结果（label/score）如下：\n")
	b.WriteString("- ")
	b.WriteString(className)
	b.WriteString(" / 1.0\n")
	b.WriteString("\n请按照以下结构输出：\n")
	b.WriteString("1) 整体描述（一句话概括）；\n")
	b.WriteString("2) 主要物体及可见特征；\n")
	b.WriteString("3) 可能的场景或用途；\n")
	b.WriteString("4) 风险或注意事项；\n")
	b.WriteString("5) 延伸建议（可选）。\n")
	msgs := []*schema.Message{
		{Role: "system", Content: "你是严谨的中文图像讲解助手"},
		{Role: "user", Content: b.String()},
	}
	resp, err := model.GenerateResponse(ctx, msgs)
	if err == nil && resp != nil {
		res.AnalysisText = resp.Content
	}
	c.JSON(http.StatusOK, res)
}

func fmtFloat(f float32) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.4f", f), "0"), ".")
}
