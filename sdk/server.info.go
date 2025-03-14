package sdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type ServerInfo struct {
	loadSync                 sync.RWMutex
	upTime                   time.Time
	OverrideConnectionGroups []OverrideConnectionGroup `json:"OverrideConnectionGroups"`
	BundleVersion            string                    `json:"BundleVersion"`
}

func (s *SDK) GetServerInfo(c *gin.Context) *ServerInfo {
	getServerInfo := func() {
		s.serverinfo.loadSync.Lock()
		defer s.serverinfo.loadSync.Unlock()
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
		data := &ConnectionGroupS{
			ConnectionGroups: make([]*ConnectionGroup, 0),
		}
		err = json.Unmarshal(bin, &data)
		for _, v := range data.ConnectionGroups {
			s.serverinfo.BundleVersion = v.BundleVersion
			s.serverinfo.OverrideConnectionGroups = v.OverrideConnectionGroups
			s.serverinfo.upTime = time.Now()
			break
		}
	}

	s.serverinfo.loadSync.RLock()
	defer s.serverinfo.loadSync.RUnlock()
	if s.serverinfo.upTime.Add(15 * time.Minute).Before(time.Now()) {
		s.serverinfo.loadSync.RUnlock()
		getServerInfo()
		s.serverinfo.loadSync.RLock()
	}

	return s.serverinfo
}

type ConnectionGroupS struct {
	ConnectionGroups []*ConnectionGroup `json:"ConnectionGroups"`
}
type ConnectionGroup struct {
	Name                       string                    `json:"Name"`
	ManagementDataUrl          string                    `json:"ManagementDataUrl"`
	IsProductionAddressables   bool                      `json:"IsProductionAddressables"`
	ApiUrl                     string                    `json:"ApiUrl"`
	GatewayUrl                 string                    `json:"GatewayUrl"`
	KibanaLogUrl               string                    `json:"KibanaLogUrl"`
	ProhibitedWordBlackListUri string                    `json:"ProhibitedWordBlackListUri"`
	ProhibitedWordWhiteListUri string                    `json:"ProhibitedWordWhiteListUri"`
	CustomerServiceUrl         string                    `json:"CustomerServiceUrl"`
	OverrideConnectionGroups   []OverrideConnectionGroup `json:"OverrideConnectionGroups"`
	BundleVersion              string                    `json:"BundleVersion"`
}

type OverrideConnectionGroup struct {
	Name                       string `json:"Name"`
	AddressablesCatalogUrlRoot string `json:"AddressablesCatalogUrlRoot"`
}

func (s *SDK) connectionGroups(c *gin.Context) {
	data := &ConnectionGroupS{
		ConnectionGroups: make([]*ConnectionGroup, 0),
	}
	sinfo := s.GetServerInfo(c)
	connectionGroup := &ConnectionGroup{
		Name:                       "Prod-Audit",
		ManagementDataUrl:          "https://prod-noticeindex.bluearchiveyostar.com/prod/index.json",
		IsProductionAddressables:   false,
		ApiUrl:                     fmt.Sprintf("%s/api/", s.GetOuterAddr()),
		GatewayUrl:                 fmt.Sprintf("%s/getEnterTicket/", s.GetOuterAddr()),
		KibanaLogUrl:               "https://prod-logcollector.bluearchiveyostar.com:5300", //    fmt.Sprintf("%s/client/log/", s.GetOuterAddr()),
		ProhibitedWordBlackListUri: "https://prod-notice.bluearchiveyostar.com/prod/ProhibitedWord/blacklist.csv",
		ProhibitedWordWhiteListUri: "https://prod-notice.bluearchiveyostar.com/prod/ProhibitedWord/whitelist.csv",
		CustomerServiceUrl:         "https://bluearchive.jp/contact-1-hint",
		BundleVersion:              sinfo.BundleVersion,
		OverrideConnectionGroups:   sinfo.OverrideConnectionGroups,
	}
	data.ConnectionGroups = append(data.ConnectionGroups, connectionGroup)

	c.JSON(200, data)
}
