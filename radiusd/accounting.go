package radiusd

import (
	"context"
	"fmt"
	"time"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2866"

	"github.com/ca17/teamsacs/constant"
	"github.com/ca17/teamsacs/models"
	"github.com/ca17/teamsacs/radiusd/debug"
	"github.com/ca17/teamsacs/radiusd/radlog"
	"github.com/ca17/teamsacs/radiusd/radparser"
)

func (s *AcctService) processAcctStart(r *radius.Request, vr *radparser.VendorRequest, username string, vpe *models.Vpe, nasrip string) {
	online := GetRadiusOnlineFromRequest(r, vr, vpe, nasrip)
	err := s.Manager.GetRadiusManager().AddRadiusOnline(online)
	if err != nil {
		radlog.Errorf("AddRadiusOnline user:%s error %s", username, err.Error())
	}
}

func (s *AcctService) processAcctUpdateBefore(r *radius.Request, vr *radparser.VendorRequest, user *models.Subscribe, vpe *models.Vpe, nasrip string) {
	// 用户状态变更为停用后触发下线
	var username = user.GetUsername()
	if user.GetStatus() == constant.DISABLED {
		s.processAcctDisconnect(r, vpe, username, nasrip)
	}

	// 用户过期后触发下线
	if user.GetExpireTime().Before(time.Now()) {
		s.processAcctDisconnect(r, vpe, username, nasrip)
	}

	s.processAcctUpdate(r, vr, username, vpe, nasrip)
}

func (s *AcctService) processAcctUpdate(r *radius.Request, vr *radparser.VendorRequest, username string, vpe *models.Vpe, nasrip string) {
	online := GetRadiusOnlineFromRequest(r, vr, vpe, nasrip)
	// 更新在线信息
	err := s.Manager.GetRadiusManager().UpdateRadiusOnlineData(online)
	if err != nil {
		radlog.Errorf("UpdateRadiusOnlineData user:%s error, %s", username, err.Error())
	}

}

func (s *AcctService) processAcctStop(r *radius.Request, vr *radparser.VendorRequest, username string, vpe *models.Vpe, nasrip string) {
	online := GetRadiusOnlineFromRequest(r, vr, vpe, nasrip)
	if err := s.Manager.GetRadiusManager().AddRadiusAccounting(online); err != nil {
		radlog.Errorf("AddRadiusAccounting user:%s error %s ", username, err.Error())
	}
	if err := s.Manager.GetRadiusManager().DeleteRadiusOnline(online.AcctSessionId); err != nil {
		radlog.Errorf("DeleteRadiusOnline user:%s error %s ", username, err.Error())
	}
}

func (s *AcctService) processAcctNasOn(r *radius.Request) {
	err := s.Manager.GetRadiusManager().BatchClearRadiusOnlineDataByNas(
		rfc2865.NASIPAddress_Get(r.Packet).String(),
		rfc2865.NASIdentifier_GetString(r.Packet),
	)
	if err != nil {
		radlog.Errorf("BatchClearRadiusOnlineDataByNas error, %s", err.Error())
	}
}

func (s *AcctService) processAcctNasOff(r *radius.Request) {
	err := s.Manager.GetRadiusManager().BatchClearRadiusOnlineDataByNas(
		rfc2865.NASIPAddress_Get(r.Packet).String(),
		rfc2865.NASIdentifier_GetString(r.Packet),
	)
	if err != nil {
		radlog.Errorf("BatchClearRadiusOnlineDataByNas error, %s", err.Error())
	}
}

func (s *AcctService) processAcctDisconnect(r *radius.Request, vpe *models.Vpe, username, nasrip string) {
	packet := radius.New(radius.CodeDisconnectRequest, []byte(vpe.GetSecret()))
	sessionid := rfc2866.AcctSessionID_GetString(r.Packet)
	if sessionid == "" {
		radlog.Errorf("radius disconnect user:%s, but sessionid is empty", username)
		return
	}
	_ = rfc2865.UserName_SetString(packet, username)
	_ = rfc2866.AcctSessionID_Set(packet, []byte(sessionid))
	var coaPort = vpe.GetCoaPort()
	radlog.Infof("disconnect user:%s => (%s:%d): %s", username, nasrip, coaPort, debug.FormatPacket(packet))
	response, err := radius.Exchange(context.Background(), packet, fmt.Sprintf("%s:%d", nasrip, coaPort))
	if err != nil {
		radlog.Errorf("radius disconnect user:%s failure", username)
	}
	radlog.Info("radius disconnect resp from (%s:%s): %s ", nasrip, coaPort, debug.FormatPacket(response))
}
