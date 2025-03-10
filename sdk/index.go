package sdk

import (
	"github.com/gin-gonic/gin"
)

type Index struct {
	Notices                       []*Notice                      `json:"Notices"`
	Events                        []*Event                       `json:"Events"`
	Issues                        []interface{}                  `json:"Issues"`
	Maintenance                   *Maintenance                   `json:"Maintenance"`
	Banners                       []*Banner                      `json:"Banners"`
	ServerStatus                  int                            `json:"ServerStatus"`
	LatestClientVersion           string                         `json:"LatestClientVersion"`
	GachaProbabilityDisplay       []*GachaProbabilityDisplay     `json:"GachaProbabilityDisplay"`
	MarketAppleId                 interface{}                    `json:"MarketAppleId"`
	NotificationBeforeMaintenance *NotificationBeforeMaintenance `json:"NotificationBeforeMaintenance"`
	ContentLock                   []interface{}                  `json:"ContentLock"`
	GachaPeriodDisplay            []*GachaPeriodDisplay          `json:"GachaPeriodDisplay"`
	RaidPeriodDisplay             interface{}                    `json:"RaidPeriodDisplay"`
	Survey                        *Survey                        `json:"Survey"`
	GuidePopup                    []*GuidePopup                  `json:"GuidePopup"`
	DownloadLimitTime             float64                        `json:"DownloadLimitTime"`
	EnableSQLite                  bool                           `json:"EnableSQLite"`
}

type Banner struct {
	BannerId            int      `json:"BannerId"`
	StartDate           string   `json:"StartDate"`
	EndDate             string   `json:"EndDate"`
	Url                 string   `json:"Url"`
	FileName            []string `json:"FileName"`
	LinkedLobbyBannerId int      `json:"LinkedLobbyBannerId"`
	BannerType          int      `json:"BannerType"`
	BannerDisplayType   int      `json:"BannerDisplayType"`
}

type Event struct {
	NoticeId     int    `json:"NoticeId"`
	StartDate    string `json:"StartDate"`
	EndDate      string `json:"EndDate"`
	Url          string `json:"Url"`
	Title        string `json:"Title"`
	DisplayOrder int    `json:"DisplayOrder"`
}

type GachaPeriodDisplay struct {
	GachaPeriodDisplayId int    `json:"GachaPeriodDisplayId"`
	Text                 string `json:"Text"`
}

type GachaProbabilityDisplay struct {
	GachaProbabilityDisplayId int    `json:"GachaProbabilityDisplayId"`
	Url                       string `json:"Url"`
	LinkedLobbyBannerId       int    `json:"LinkedLobbyBannerId"`
	GachaDisplayTag           int    `json:"GachaDisplayTag"`
}

type GuidePopup struct {
	GuidePopupId   int         `json:"GuidePopupId"`
	GuidePopupType int         `json:"GuidePopupType"`
	PopupType      int         `json:"PopupType"`
	StartDate      string      `json:"StartDate"`
	EndDate        string      `json:"EndDate"`
	FileName       string      `json:"FileName"`
	Url            string      `json:"Url"`
	Message        string      `json:"Message"`
	SurveyId       int         `json:"SurveyId"`
	NotifyUrl      interface{} `json:"NotifyUrl"`
	GotoUrl        *string     `json:"GotoUrl"`
	DisplayOrder   int         `json:"DisplayOrder"`
	PopupOKText    string      `json:"PopupOKText"`
}

type Maintenance struct {
	StartDate string `json:"StartDate"`
	EndDate   string `json:"EndDate"`
	Text      string `json:"Text"`
}

type Notice struct {
	NoticeId     int    `json:"NoticeId"`
	StartDate    string `json:"StartDate"`
	EndDate      string `json:"EndDate"`
	Url          string `json:"Url"`
	Title        string `json:"Title"`
	DisplayOrder int    `json:"DisplayOrder"`
}

type NotificationBeforeMaintenance struct {
	PopupType int    `json:"PopupType"`
	StartDate string `json:"StartDate"`
	EndDate   string `json:"EndDate"`
	Text      string `json:"Text"`
}

type Survey struct {
	SurveyId  int         `json:"SurveyId"`
	PopupType int         `json:"PopupType"`
	StartDate string      `json:"StartDate"`
	EndDate   string      `json:"EndDate"`
	FileName  string      `json:"FileName"`
	Url       string      `json:"Url"`
	Text      string      `json:"Text"`
	NotifyUrl interface{} `json:"NotifyUrl"`
}

func index(c *gin.Context) {
	// url := c.Request.URL.String()
	// resp, err := http.Get("https://prod-noticeindex.bluearchiveyostar.com" + url)
	// if err != nil {
	// 	return
	// }
	// defer resp.Body.Close()
	// bin, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return
	// }
	// data := new(Index)
	//
	// if json.Unmarshal(bin, &data) != nil {
	// 	return
	// }
	// data.ServerStatus = 1
	// data.Maintenance = &Maintenance{
	// 	StartDate: "2025-02-02T11:00:00",
	// 	EndDate:   "2025-02-02T17:00:00",
	// 	Text:      "定期维护中",
	// }
	// data.NotificationBeforeMaintenance = &NotificationBeforeMaintenance{
	// 	PopupType: 1,
	// 	StartDate: "2023-12-19T00:00:00+09:00",
	// 	EndDate:   "2023-12-27T11:00:00",
	// 	Text:      "更新",
	// }
	// data.Survey = &Survey{
	// 	SurveyId:  199990730,
	// 	PopupType: 1,
	// 	StartDate: "2023-01-11T14:32:37.234028+09:00",
	// 	EndDate:   "2023-01-23T23:59:59",
	// 	FileName:  "Lobby_Banner.png",
	// 	Url:       "https://prod-notice.bluearchiveyostar.com/prod/SurveyGuidePopup/199990730/",
	// 	Text:      "ブルアカ2周年　特別アンケート",
	// 	NotifyUrl: nil,
	// }
	// c.JSON(http.StatusOK, data)
	// return

	i := &Index{
		Banners:                 make([]*Banner, 0),
		ContentLock:             make([]interface{}, 0),
		DownloadLimitTime:       0,
		EnableSQLite:            true,
		Events:                  make([]*Event, 0),
		GachaPeriodDisplay:      make([]*GachaPeriodDisplay, 0),
		GachaProbabilityDisplay: make([]*GachaProbabilityDisplay, 0),
		GuidePopup:              make([]*GuidePopup, 0),
		Issues:                  make([]interface{}, 0),
		LatestClientVersion:     "1.54.327262",
		Maintenance: &Maintenance{
			StartDate: "2025-02-02T11:00:00",
			EndDate:   "2025-02-02T17:00:00",
			Text:      "定期维护中",
		},
		ServerStatus: 1,
		NotificationBeforeMaintenance: &NotificationBeforeMaintenance{
			PopupType: 1,
			StartDate: "2023-12-19T00:00:00+09:00",
			EndDate:   "2023-12-27T11:00:00",
			Text:      "更新",
		},
		MarketAppleId:     nil,
		Notices:           make([]*Notice, 0),
		RaidPeriodDisplay: nil,
		Survey: &Survey{
			SurveyId:  199990730,
			PopupType: 1,
			StartDate: "2023-01-11T14:32:37.234028+09:00",
			EndDate:   "2023-01-23T23:59:59",
			FileName:  "Lobby_Banner.png",
			Url:       "https://prod-notice.bluearchiveyostar.com/prod/SurveyGuidePopup/199990730/",
			Text:      "ブルアカ2周年　特別アンケート",
			NotifyUrl: nil,
		},
	}
	c.JSON(200, i)
}
