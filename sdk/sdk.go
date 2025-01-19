package sdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/common/code"
	"github.com/gucooing/BaPs/config"
	"github.com/gucooing/BaPs/pkg/alg"
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
	s.router.GET("/:url.json", s.connectionGroups)

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

	gucooingApi := s.router.Group("/gucooing/api", alg.AutoGucooingApi())
	{
		gucooingApi.GET("/ba/getEmailCode", s.getEmailCode)
	}
}

func handleIndex(c *gin.Context) {
	c.String(http.StatusOK, "Ba Ps!")
}

func (s *SDK) GetOuterAddr() string {
	cfg := config.GetHttpNet()
	if cfg.Tls {
		return fmt.Sprintf("https://%s:%s", cfg.OuterAddr, cfg.OuterPort)
	} else {
		return fmt.Sprintf("http://%s:%s", cfg.OuterAddr, cfg.OuterPort)
	}
}

type ConnectionGroupS struct {
	ConnectionGroups []*ConnectionGroup `json:"ConnectionGroups"`
}
type ConnectionGroup struct {
	Name                       string                     `json:"Name"`
	ManagementDataUrl          string                     `json:"ManagementDataUrl"`
	IsProductionAddressables   bool                       `json:"IsProductionAddressables"`
	ApiUrl                     string                     `json:"ApiUrl"`
	GatewayUrl                 string                     `json:"GatewayUrl"`
	KibanaLogUrl               string                     `json:"KibanaLogUrl"`
	ProhibitedWordBlackListUri string                     `json:"ProhibitedWordBlackListUri"`
	ProhibitedWordWhiteListUri string                     `json:"ProhibitedWordWhiteListUri"`
	CustomerServiceUrl         string                     `json:"CustomerServiceUrl"`
	OverrideConnectionGroups   []*OverrideConnectionGroup `json:"OverrideConnectionGroups"`
	BundleVersion              string                     `json:"BundleVersion"`
}

type OverrideConnectionGroup struct {
	Name                       string `json:"Name"`
	AddressablesCatalogUrlRoot string `json:"AddressablesCatalogUrlRoot"`
}

func (s *SDK) connectionGroups(c *gin.Context) {
	url := c.Request.URL.String()
	resp, err := http.Get("https://yostar-serverinfo.bluearchiveyostar.com" + url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bin, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	data := new(ConnectionGroupS)

	if json.Unmarshal(bin, &data) != nil {
		return
	}
	for _, group := range data.ConnectionGroups {
		group.ApiUrl = fmt.Sprintf("%s/api/", s.GetOuterAddr())
		group.GatewayUrl = fmt.Sprintf("%s/getEnterTicket/", s.GetOuterAddr())
	}

	c.JSON(200, data)
}
