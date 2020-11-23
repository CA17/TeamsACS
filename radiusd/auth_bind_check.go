package radiusd

import (
	"errors"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/models"
	"github.com/ca17/teamsacs/radiusd/radparser"
)

// CheckVlanBind
// vlanid binding detection
// Only if both user vlanid and request vlanid are valid.
// If user vlanid is empty, update user vlanid directly.
func (s *AuthService) CheckVlanBind(user *models.Subscribe, vendorReq *radparser.VendorRequest) error {
	if user.GetBindVlan() == 0 {
		return nil
	}
	reqvid1 := int(vendorReq.Vlanid1)
	reqvid2 := int(vendorReq.Vlanid2)
	var vlanid1 = user.GetVlan1()
	var vlanid2 = user.GetVlan2()
	if vlanid1 != 0 && vendorReq.Vlanid1 != 0 && vlanid1 != reqvid1 {
		return errors.New("user vlanid1 bind not match")
	}

	if vlanid2 != 0 && reqvid2 != 0 && vlanid2 != reqvid2 {
		return errors.New("user vlanid2 bind not match")
	}

	return nil
}

// CheckMacBind
// mac binding detection
// Detected only if both user mac and request mac are valid.
// If user mac is empty, update user mac directly.
func (s *AuthService) CheckMacBind(user *models.Subscribe, vendorReq *radparser.VendorRequest) error {
	if user.GetBindVlan() == 0 {
		return nil
	}
	var mac = user.GetMacAddr()
	if common.IsNotEmptyAndNA(mac) && vendorReq.Macaddr != "" && mac != vendorReq.Macaddr {
		return errors.New("user mac bind not match")
	}
	return nil
}


// UpdateBind
// update mac or vlan
func (s *AuthService) UpdateBind(user *models.Subscribe, vendorReq *radparser.VendorRequest) {
	var mac = user.GetMacAddr()
	var username = user.GetUsername()
	var vlanid1 = user.GetVlan1()
	var vlanid2 = user.GetVlan2()
	if mac != vendorReq.Macaddr {
		s.UpdateUserMac(username, vendorReq.Macaddr)
	}
	reqvid1 := int(vendorReq.Vlanid1)
	reqvid2 := int(vendorReq.Vlanid2)
	if vlanid1 != reqvid1 {
		s.UpdateUserVlanid2(username, reqvid1)
	}
	if vlanid2 != reqvid2 {
		s.UpdateUserVlanid2(username, reqvid2)
	}
}
