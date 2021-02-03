package authorization

import (
	"fmt"

	"layeh.com/radius"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/radiusd/vendors/cisco"
)

func CiscoAuthorization(prof Profile, accept *radius.Packet) {
	upLimitPolicy := prof.GetUpLimitPolicy()
	if common.IsNotEmptyAndNA(upLimitPolicy) {
		cisco.CiscoAVPair_Add(accept, []byte(fmt.Sprintf("sub-qos-policy-in=%s", upLimitPolicy)))
	}
	downLimitPolicy := prof.GetDownLimitPolicy()
	if common.IsNotEmptyAndNA(downLimitPolicy) {
		cisco.CiscoAVPair_Add(accept, []byte(fmt.Sprintf("sub-qos-policy-out=%s", downLimitPolicy)))
	}
}
