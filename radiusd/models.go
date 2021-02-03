package radiusd

import "time"

type AuthorizationProfile struct {
	ExpireTime      time.Time
	InterimInterval int
	AddrPool        string
	Ipaddr          string
	UpRateKbps      int
	DownRateKbps    int
	Domain          string
	LimitPolicy     string
	UpLimitPolicy   string
	DownLimitPolicy string
}

func (a AuthorizationProfile) GetExpireTime() time.Time {
	return a.ExpireTime
}

func (a AuthorizationProfile) GetInterimInterval() int {
	return a.InterimInterval
}

func (a AuthorizationProfile) GetAddrPool() string {
	return a.AddrPool
}

func (a AuthorizationProfile) GetIpaddr() string {
	return a.Ipaddr
}

func (a AuthorizationProfile) GetUpRateKbps() int {
	return a.UpRateKbps
}

func (a AuthorizationProfile) GetDownRateKbps() int {
	return a.DownRateKbps
}

func (a AuthorizationProfile) GetDomain() string {
	return a.Domain
}

func (a AuthorizationProfile) GetLimitPolicy() string {
	return a.LimitPolicy
}

func (a AuthorizationProfile) GetUpLimitPolicy() string {
	return a.UpLimitPolicy
}

func (a AuthorizationProfile) GetDownLimitPolicy() string {
	return a.DownLimitPolicy
}
