package authorization

import (
	"fmt"

	"layeh.com/radius"

	"github.com/ca17/teamsacs/radiusd/vendors/mikrotik"
)

func MikrotikAuthorization(prof Profile, accept *radius.Packet) {
	_ = mikrotik.MikrotikRateLimit_SetString(accept, fmt.Sprintf("%dk/%dk", prof.GetUpRateKbps(), prof.GetDownRateKbps()))
}
