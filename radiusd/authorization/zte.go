package authorization

import (
	"math"

	"layeh.com/radius"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/radiusd/vendors/zte"
)

func ZteAuthorization(prof Profile, accept *radius.Packet) {
	var up = prof.GetUpRateKbps() * 1000
	var down = prof.GetDownRateKbps() * 1000
	if up > math.MaxInt32 {
		up = math.MaxInt32
	}
	if down > math.MaxInt32 {
		down = math.MaxInt32
	}

	zte.ZTERateCtrlSCRUp_Set(accept, zte.ZTERateCtrlSCRUp(up))
	zte.ZTERateCtrlSCRDown_Set(accept, zte.ZTERateCtrlSCRDown(down))

	domain := prof.GetDomain()
	if common.IsNotEmptyAndNA(domain) {
		zte.ZTEContextName_SetString(accept, domain)
	}
}
