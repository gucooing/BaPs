package sdk

import (
	"fmt"
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

	gucooingApi := s.router.Group("/gucooing/api", s.autoGucooingApi())
	{
		gucooingApi.GET("/ba/getEmailCode", s.getEmailCode)
		gucooingApi.GET("/ba/getPlayerBin", s.getPlayerBin)
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
				c.Abort()
			}
		}
	}
}

func (s *SDK) GetOuterAddr() string {
	cfg := config.GetHttpNet()
	if cfg.Tls {
		return fmt.Sprintf("https://%s:%s", cfg.OuterAddr, cfg.OuterPort)
	} else {
		return fmt.Sprintf("http://%s:%s", cfg.OuterAddr, cfg.OuterPort)
	}
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
	type ConnectionGroupS struct {
		ConnectionGroups []*ConnectionGroup `json:"ConnectionGroups"`
	}
	rsp := &ConnectionGroupS{
		ConnectionGroups: []*ConnectionGroup{},
	}
	connectionGroup := &ConnectionGroup{
		Name:                       "Prod-Audit",
		ManagementDataUrl:          "https://prod-noticeindex.bluearchiveyostar.com/prod/index.json",
		IsProductionAddressables:   true,
		ApiUrl:                     fmt.Sprintf("%s/api/", s.GetOuterAddr()),
		GatewayUrl:                 fmt.Sprintf("%s/getEnterTicket/", s.GetOuterAddr()),
		KibanaLogUrl:               "https://prod-logcollector.bluearchiveyostar.com:5300",
		ProhibitedWordBlackListUri: "https://prod-notice.bluearchiveyostar.com/prod/ProhibitedWord/blacklist.csv",
		ProhibitedWordWhiteListUri: "https://prod-notice.bluearchiveyostar.com/prod/ProhibitedWord/whitelist.csv",
		CustomerServiceUrl:         "https://bluearchive.jp/contact-1-hint",
		OverrideConnectionGroups: []*OverrideConnectionGroup{
			{
				Name:                       "1.0",
				AddressablesCatalogUrlRoot: "https://prod-clientpatch.bluearchiveyostar.com/m28_1_0_1_mashiro3",
			},
			{
				Name:                       "1.52",
				AddressablesCatalogUrlRoot: "https://prod-clientpatch.bluearchiveyostar.com/r75_49ajrpwcziy395uuk0jq_2",
			},
		},
		BundleVersion: "li3pmyogha",
	}
	rsp.ConnectionGroups = append(rsp.ConnectionGroups, connectionGroup)
	c.JSON(200, rsp)
}
