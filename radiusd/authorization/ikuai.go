package authorization

import (
	"math"

	"layeh.com/radius"

	"github.com/ca17/teamsacs/radiusd/vendors/ikuai"
)

func IkuaiAuthorization(prof Profile, accept *radius.Packet) {
	var up = prof.GetUpRateKbps() * 1024 * 8
	var down = prof.GetDownRateKbps() * 1024 * 8
	if up > math.MaxInt32 {
		up = math.MaxInt32
	}
	if down > math.MaxInt32 {
		down = math.MaxInt32
	}

	ikuai.RPUpstreamSpeedLimit_Set(accept, ikuai.RPUpstreamSpeedLimit(up))
	ikuai.RPDownstreamSpeedLimit_Set(accept, ikuai.RPDownstreamSpeedLimit(down))
}
