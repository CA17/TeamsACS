package radiusd

import (
	"fmt"

	"layeh.com/radius"

	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/models"
)

func ListenRadiusAuthServer(manager *models.ModelManager) error {
	service := NewAuthService(NewRadiusService(manager))
	server := radius.PacketServer{
		Addr:               fmt.Sprintf("%s:%d", manager.Config.Radiusd.Host, manager.Config.Radiusd.AuthPort),
		Handler:            service,
		SecretSource:       service,
		InsecureSkipVerify: true,
	}

	log.Infof("Starting Radius Auth server on %s", server.Addr)
	return server.ListenAndServe()
}

func ListenRadiusAcctServer(manager *models.ModelManager) error {
	service := NewAcctService(NewRadiusService(manager))
	server := radius.PacketServer{
		Addr:               fmt.Sprintf("%s:%d", manager.Config.Radiusd.Host, manager.Config.Radiusd.AcctPort),
		Handler:            service,
		SecretSource:       service,
		InsecureSkipVerify: true,
	}

	log.Infof("Starting Radius Acct server on %s", server.Addr)
	return server.ListenAndServe()
}
