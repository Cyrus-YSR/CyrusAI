package user

import (
	"context"
	"net/http"

	"GopherAI/common/code"
	"GopherAI/config"
	"GopherAI/controller"
	"GopherAI/service/user"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
)

type (
	LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	LoginResponse struct {
		controller.Response
		Token string `json:"token,omitempty"`
	}
	RegisterRequest struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Captcha  string `json:"captcha"`
		Password string `json:"password" binding:"required"`
	}
	RegisterResponse struct {
		controller.Response
		Token string `json:"token,omitempty"`
	}

	CaptchaRequest struct {
		Email string `json:"email" binding:"required"`
	}

	CaptchaResponse struct {
		controller.Response
	}

	GoogleLoginRequest struct {
		Credential string `json:"credential" binding:"required"`
	}

	GoogleLoginResponse struct {
		controller.Response
		Token string `json:"token,omitempty"`
	}
)

func Login(c *gin.Context) {
	req := new(LoginRequest)
	res := new(LoginResponse)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	token, code_ := user.Login(req.Username, req.Password)
	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}

	res.Success()
	res.Token = token
	c.JSON(http.StatusOK, res)

}

func Register(c *gin.Context) {
	req := new(RegisterRequest)
	res := new(RegisterResponse)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	token, code_ := user.Register(req.Username, req.Email, req.Password, req.Captcha)
	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}

	res.Success()
	res.Token = token
	c.JSON(http.StatusOK, res)
}

func HandleCaptcha(c *gin.Context) {
	req := new(CaptchaRequest)
	res := new(CaptchaResponse)
	//解析参数
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	//给service层进行处理
	code_ := user.SendCaptcha(req.Email)
	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}
	//匿名字段，其实本身res.Success()调用就是res.Response.Success()
	//res.Response.Success()
	res.Success()
	c.JSON(http.StatusOK, res)
}

func GoogleLogin(c *gin.Context) {
	req := new(GoogleLoginRequest)
	res := new(GoogleLoginResponse)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	// 验证 Google token
	// 从配置文件中读取 Google Client ID
	clientID := config.GetConfig().GoogleClientID
	payload, err := idtoken.Validate(context.Background(), req.Credential, clientID)
	if err != nil {
		// 这里验证失败通常是凭证无效
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	email := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)

	token, code_ := user.GoogleLogin(email, name)
	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}

	res.Success()
	res.Token = token
	c.JSON(http.StatusOK, res)
}
