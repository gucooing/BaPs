package gateway

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/db"
	"github.com/gucooing/BaPs/mx/proto"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *Gateway) getEnterTicket(c *gin.Context) {
	if !alg.CheckGateWay(c) {
		c.JSON(404, gin.H{})
		return
	}
	bin, err := alg.GetFormMx(c)
	if err != nil {
		return
	}
	rsp := &proto.QueuingGetTicketResponse{}
	defer g.send(c, rsp)
	req := new(proto.QueuingGetTicketRequest)
	err = json.Unmarshal(bin, req)
	if err != nil {
		logger.Debug("request err:%s c--->s:%s", err.Error(), string(bin))
		return
	}
	yoStarUserLogin := db.GetYoStarUserLoginByYostarUid(req.YostarUID)
	if yoStarUserLogin == nil {
		return
	}
	if yoStarUserLogin.YostarLoginToken != req.YostarToken ||
		yoStarUserLogin.YostarLoginToken == "" {
		return
	}
	yoStarUserLogin.YostarLoginToken = ""
	if err = db.UpdateYoStarUserLogin(yoStarUserLogin); err != nil {
		return
	}
	enterTicket := fmt.Sprintf("%v%s", g.snow.GenId(), alg.RandStr(10))
	if !enter.AddEnterTicket(yoStarUserLogin.AccountServerId, enterTicket) {
		return
	}
	rsp.EnterTicket = enterTicket
	rsp.Protocol = req.Protocol
	logger.Debug("EnterTicket交换成功:%s", rsp.EnterTicket)
}
