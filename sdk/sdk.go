package sdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/config"
)

type SDK struct {
	router *gin.Engine
}

func New(router *gin.Engine) *SDK {
	s := &SDK{
		router: router,
	}

	s.initRouter()
	return s
}

func (s *SDK) initRouter() {
	s.router.Any("/", handleIndex)
	s.router.GET("/:url.json", s.connectionGroups)

	s.router.GET("/prod/index.json", index)

	account := s.router.Group("/account")
	{
		account.POST("/yostar_auth_request", s.YostarAuthRequest)
		account.POST("/yostar_auth_submit", s.YostarAuthSubmit)
	}
	user := s.router.Group("/user")
	{
		user.POST("/yostar_createlogin", s.YostarCreatelogin)
		user.POST("/login", s.YostarLogin)
		user.POST("/agreement", agreement)
	}
	hash := s.router.Group("/r76_d32k9xg20divs4806ybp_2")
	{
		hash.GET("/TableBundles/TableCatalog.hash", func(c *gin.Context) {
			c.String(200, "4136508985")
		})
		hash.GET("/MediaResources/Catalog/MediaCatalog.hash", func(c *gin.Context) {
			c.String(200, "2581543713")
		})
		hash.GET("/iOS/bundleDownloadInfo.hash", func(c *gin.Context) {
			c.String(200, "970608301")
		})
		hash.GET("/iOS/catalog_iOS.hash", func(c *gin.Context) {
			c.String(200, "a0d3861cb71f215d5d3033d3eee04172")
		})
	}
	app := s.router.Group("/app")
	{
		app.Any("/getSettings", getSettings)
		app.Any("/getCode", getCode)
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
		group.ManagementDataUrl = fmt.Sprintf("%s/prod/index.json", s.GetOuterAddr())
		group.ProhibitedWordBlackListUri = "https://ba-oss.alsl.xyz/prod/ProhibitedWord/blacklist.csv"
		group.ProhibitedWordWhiteListUri = "https://ba-oss.alsl.xyz/prod/ProhibitedWord/whitelist.csv"
	}

	c.JSON(200, data)
}
