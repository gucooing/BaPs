package sdk

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/common/code"
	"github.com/gucooing/BaPs/config"
)

type SDK struct {
	router *gin.Engine
	code   *code.Code
}

func New(router *gin.Engine) *SDK {
	s := &SDK{
		router: router,
		code:   code.NewCode(),
	}

	s.initRouter()
	go s.code.CheckCodeTime()
	return s
}

func (s *SDK) initRouter() {
	s.router.Any("/", handleIndex)

	account := s.router.Group("/account")
	{
		account.POST("/yostar_auth_request", s.YostarAuthRequest)
		account.POST("/yostar_auth_submit", s.YostarAuthSubmit)
	}
	user := s.router.Group("/user")
	{
		user.POST("/yostar_createlogin", s.YostarCreatelogin)
		user.POST("/login", s.YostarLogin)
	}

	gucooingApi := s.router.Group("/gucooing/api", s.autoGucooingApi())
	{
		gucooingApi.GET("/getEmailCode", s.getEmailCode)
	}
}

func handleIndex(c *gin.Context) {
	c.String(http.StatusOK, "Ba Ps!")
}

func (s *SDK) autoGucooingApi() gin.HandlerFunc {
	if config.GetGucooingApiKey() == "" {
		return func(c *gin.Context) {}
	} else {
		return func(c *gin.Context) {
			if c.GetHeader("Authorization-Gucooing") != config.GetGucooingApiKey() {
				c.String(401, "Unauthorized")
			}
		}
	}
}
