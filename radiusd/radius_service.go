package radiusd

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"time"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2866"
	"layeh.com/radius/rfc2869"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/config"
	"github.com/ca17/teamsacs/models"
	"github.com/ca17/teamsacs/radiusd/radlog"
	"github.com/ca17/teamsacs/radiusd/radparser"
)

const (
	VendorMikrotik = "14988"
	VendorIkuai    = "10055"
	VendorHuawei   = "2011"
	VendorZte      = "3902"
	VendorH3c      = "25506"
	VendorRadback  = "2352"
	VendorCisco    = "9"

	RadiusAuthlogAll  = "all"
	RadiusAuthlogNone = "none"
	RadiusAuthSucces  = "success"
	RadiusAuthFailure = "failure"
)

type RadiusService struct {
	Manager *models.ModelManager
}

func NewRadiusService(manager *models.ModelManager) *RadiusService {
	return &RadiusService{Manager: manager}
}

func (s *RadiusService) GetAppConfig() *config.AppConfig {
	return s.Manager.Config
}

func (s *RadiusService) RADIUSSecret(ctx context.Context, remoteAddr net.Addr) ([]byte, error) {
	return []byte("greensecret"), nil
}

// 查询 NAS 设备, 优先查询IP, 然后ID
func (s *RadiusService) GetNas(ip, identifier string) (*models.Vpe, error) {
	vstore := s.Manager.GetVpeManager()
	vpe, err := vstore.GetVpeByIpaddr(ip)
	if err != nil {
		nvpe, err := vstore.GetVpeByIdentifier(identifier)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("Unauthorized access to device, Ip=%s, Identifier=%s, %s", ip, identifier, err.Error())
			}
			return nil, err
		}
		return nvpe, nil
	}
	return vpe, nil
}

// 获取有效用户, 初步判断用户有效性
func (s *RadiusService) GetUser(username string, macauth bool) (*models.Subscribe, error) {
	m := s.Manager.GetSubscribeManager()
	user := new(models.Subscribe)
	var err error
	if macauth {
		user, err = m.GetSubscribeByUser(username)
		if err != nil {
			return nil, err
		}
	} else {
		user, err = m.GetSubscribeByMac(username)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("user:%s not exists", username)
			}
			return nil, err
		}
	}
	if user.GetStatus() == common.DISABLED {
		return nil, fmt.Errorf("user:%s status is disabled", username)
	}

	if user.GetExpireTime().Before(time.Now()) {
		return nil, fmt.Errorf("user:%s expire", username)
	}
	return user, nil
}

// 获取用户, 不判断用户过期等状态
func (s *RadiusService) GetUserForAcct(username string) (*models.Subscribe, error) {
	m := s.Manager.GetSubscribeManager()
	user, err := m.GetSubscribeByUser(username)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (s *RadiusService) UpdateUserMac(username string, macaddr string) {
	err := s.Manager.GetSubscribeManager().UpdateSubscribeByUsername(username, models.Attributes{"macaddr": macaddr})
	if err != nil {
		radlog.Warningf("update user:%s mac_addr:%s error", username, macaddr)
	}
}

func (s *RadiusService) UpdateUserVlanid1(username string, vlanid1 int) {
	err := s.Manager.GetSubscribeManager().UpdateSubscribeByUsername(username, models.Attributes{"vlanid1": vlanid1})
	if err != nil {
		radlog.Warningf("update user:%s vlanid1:%s error", username, vlanid1)
	}
}

func (s *RadiusService) UpdateUserVlanid2(username string, vlanid2 int) {
	err := s.Manager.GetSubscribeManager().UpdateSubscribeByUsername(username, models.Attributes{"vlanid2": vlanid2})
	if err != nil {
		radlog.Warningf("update user:%s vlanid2:%s error", username, vlanid2)
	}
}

func (s *RadiusService) GetIntConfig(name string, defval int64) int64 {
	return s.Manager.GetConfigManager().GetRadiusConfigIntValue(name, defval)
}

func (s *RadiusService) GetStringConfig(name string, defval string) string {
	return s.Manager.GetConfigManager().GetRadiusConfigStringValue(name, defval)
}

func GetRadiusOnlineFromRequest(r *radius.Request, vr *radparser.VendorRequest, vpe *models.Vpe, nasrip string) models.Accounting {

	acctInputOctets := int(rfc2866.AcctInputOctets_Get(r.Packet))
	acctInputGigawords := int(rfc2869.AcctInputGigawords_Get(r.Packet))
	acctOutputOctets := int(rfc2866.AcctOutputOctets_Get(r.Packet))
	acctOutputGigawords := int(rfc2869.AcctOutputGigawords_Get(r.Packet))

	getAcctStartTime := func(sessionTime int) time.Time {
		m, _ := time.ParseDuration(fmt.Sprintf("-%ds", sessionTime))
		return time.Now().Add(m)
	}
	return models.Accounting{
		Username:          rfc2865.UserName_GetString(r.Packet),
		NasId:             common.IfEmptyStr(rfc2865.NASIdentifier_GetString(r.Packet), common.NA),
		NasAddr:           vpe.GetIpaddr(),
		NasPaddr:          nasrip,
		SessionTimeout:    int(rfc2865.SessionTimeout_Get(r.Packet)),
		FramedIpaddr:      common.IfEmptyStr(rfc2865.FramedIPAddress_Get(r.Packet).String(), common.NA),
		FramedNetmask:     common.IfEmptyStr(rfc2865.FramedIPNetmask_Get(r.Packet).String(), common.NA),
		MacAddr:           common.IfEmptyStr(vr.Macaddr, common.NA),
		NasPort:           0,
		NasClass:          common.NA,
		NasPortId:         common.IfEmptyStr(rfc2869.NASPortID_GetString(r.Packet), common.NA),
		NasPortType:       0,
		ServiceType:       0,
		AcctSessionId:     rfc2866.AcctSessionID_GetString(r.Packet),
		AcctSessionTime:   int(rfc2866.AcctSessionTime_Get(r.Packet)),
		AcctInputTotal:    int64(acctInputOctets) + int64(acctInputGigawords)*4*1024*1024*1024,
		AcctOutputTotal:   int64(acctOutputOctets) + int64(acctOutputGigawords)*4*1024*1024*1024,
		AcctInputPackets:  int(rfc2866.AcctInputPackets_Get(r.Packet)),
		AcctOutputPackets: int(rfc2866.AcctInputPackets_Get(r.Packet)),
		AcctStartTime:     getAcctStartTime(int(rfc2866.AcctSessionTime_Get(r.Packet))),
		LastUpdate:        time.Now(),
	}

}
