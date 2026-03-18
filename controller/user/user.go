package user

import (
	"context"
	"encoding/json"
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
		Credential    string `json:"credential" binding:"required"`
		IsAccessToken bool   `json:"is_access_token"`
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

	var email, name string

	if req.IsAccessToken {
		// Use access token to get user info
		reqURL := "https://www.googleapis.com/oauth2/v3/userinfo"
		httpReq, _ := http.NewRequest("GET", reqURL, nil)
		httpReq.Header.Set("Authorization", "Bearer "+req.Credential)

		client := &http.Client{}
		resp, err := client.Do(httpReq)
		if err != nil || resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
			return
		}
		defer resp.Body.Close()

		var userInfo struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		}
		importJson := json.NewDecoder(resp.Body)
		if err := importJson.Decode(&userInfo); err != nil {
			c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
			return
		}
		email = userInfo.Email
		name = userInfo.Name
	} else {
		// Validate Google ID token
		clientID := config.GetConfig().GoogleClientID
		payload, err := idtoken.Validate(context.Background(), req.Credential, clientID)
		if err != nil {
			c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
			return
		}
		email = payload.Claims["email"].(string)
		name, _ = payload.Claims["name"].(string)
	}

	token, code_ := user.GoogleLogin(email, name)
	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}

	res.Success()
	res.Token = token
	c.JSON(http.StatusOK, res)
}
