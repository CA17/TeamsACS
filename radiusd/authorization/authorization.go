package authorization

import (
	"math"
	"net"
	"time"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2869"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/radiusd/vendors"
)

type Profile interface {
	GetExpireTime() time.Time
	GetInterimInterval() int
	GetAddrPool() string
	GetIpaddr() string
	GetUpRateKbps() int
	GetDownRateKbps() int
	GetDomain() string
	GetLimitPolicy() string
	GetUpLimitPolicy() string
	GetDownLimitPolicy() string
}

func UpdateAuthorization(profile Profile, vendorCode string, accept *radius.Packet) {
	DefaultAuthorization(profile, accept)
	switch vendorCode {
	case vendors.VendorHuawei:
		HuaweiAuthorization(profile, accept)
	case vendors.VendorH3c:
		H3CAuthorization(profile, accept)
	case vendors.VendorRadback:
		RadbackAuthorization(profile, accept)
	case vendors.VendorZte:
		ZteAuthorization(profile, accept)
	case vendors.VendorCisco:
		CiscoAuthorization(profile, accept)
	case vendors.VendorMikrotik:
		MikrotikAuthorization(profile, accept)
	case vendors.VendorIkuai:
		IkuaiAuthorization(profile, accept)
	}
}

func DefaultAuthorization(prof Profile, accept *radius.Packet) {
	var timeout = int64(prof.GetExpireTime().Sub(time.Now()).Seconds())
	if timeout > math.MaxInt32 {
		timeout = math.MaxInt32
	}
	_ = rfc2865.SessionTimeout_Set(accept, rfc2865.SessionTimeout(timeout))

	var interimTimes = prof.GetInterimInterval()
	_ = rfc2869.AcctInterimInterval_Set(accept, rfc2869.AcctInterimInterval(interimTimes))

	addrpool := prof.GetAddrPool()
	if common.IsNotEmptyAndNA(addrpool) {
		rfc2869.FramedPool_SetString(accept, addrpool)
	}
	ipaddr := prof.GetIpaddr()
	if common.IsNotEmptyAndNA(ipaddr) {
		rfc2865.FramedIPAddress_Set(accept, net.ParseIP(ipaddr))
	}
}
