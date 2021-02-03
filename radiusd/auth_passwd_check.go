package radiusd

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"strings"

	"layeh.com/radius"
	"layeh.com/radius/rfc2759"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc3079"

	"github.com/ca17/teamsacs/common/aes"
	"github.com/ca17/teamsacs/constant"
	"github.com/ca17/teamsacs/models"
	radlog "github.com/ca17/teamsacs/radiusd/radlog"
	"github.com/ca17/teamsacs/radiusd/vendors/microsoft"
)

func (s *AuthService) GetLocalPassword(user *models.Subscribe, isMacAuth bool) (string, error) {
	if isMacAuth {
		return user.GetMacAddr(), nil
	}
	localpwd, err := aes.DecryptFromB64(user.GetPassword(), s.GetAppConfig().System.Aeskey)
	if err != nil {
		return "", fmt.Errorf("user:%s local password is invalid", user.GetUsername())
	}
	return localpwd, nil
}

// check password
// passward is not empty for PAP authentication.
// chapPassword is not empty for chap authentication.
func (s *AuthService) CheckPassword(r *radius.Request, username, localpassword string, radAccept *radius.Packet, isMacAuth bool) error {
	ignoreChk := s.GetStringConfig(constant.RadiusIgnorePwd, constant.DISABLED) == constant.ENABLED
	password := rfc2865.UserPassword_GetString(r.Packet)
	challenge := microsoft.MSCHAPChallenge_Get(r.Packet)
	response := microsoft.MSCHAP2Response_Get(r.Packet)

	if ignoreChk && challenge != nil {
		return nil
	}

	// mschap 认证
	if challenge != nil && response != nil {
		return s.CheckMsChapPassword(username, localpassword, challenge, response, radAccept)
	}

	chapPassword := rfc2865.CHAPPassword_Get(r.Packet)
	if chapPassword != nil && !ignoreChk && !isMacAuth {
		chapChallenge := rfc2865.CHAPChallenge_Get(r.Packet)
		if len(chapPassword) != 17 {
			return fmt.Errorf("user:%s chap password must be 17 bytes", username)
		}

		if len(chapChallenge) != 16 {
			return fmt.Errorf("user:%s chap challenge must be 16 bytes", username)
		}

		w := md5.New()
		w.Write([]byte{chapPassword[0]})
		w.Write([]byte(localpassword))
		w.Write(chapChallenge)
		md5r := w.Sum(nil)
		/*
					for (int i = 0; i < 16; i++)
			            if (chapHash[i] != chapPassword[i + 1])
			                return false;
		*/
		for i := 0; i < 16; i++ {
			if md5r[i] != chapPassword[i+1] {
				return fmt.Errorf("user:%s chap password error", username)
			}
		}

		return nil
	}

	if password != "" && !ignoreChk && !isMacAuth {
		if strings.TrimSpace(password) != localpassword {
			return fmt.Errorf("user:%s pap password is not match", username)
		}
	}

	return nil
}

func (s *AuthService) CheckMsChapPassword(username, password string, challenge, response []byte, radAccept *radius.Packet) error {
	if len(challenge) == 16 && len(response) == 50 {
		ident := response[0]
		peerChallenge := response[2:18]
		peerResponse := response[26:50]
		byteUser := []byte(username)
		bytePwd := []byte(password)
		ntResponse, err := rfc2759.GenerateNTResponse(challenge, peerChallenge, byteUser, bytePwd)
		if err != nil {
			return fmt.Errorf("user:%s mschap access mschap access cannot generate ntResponse", username)
		}

		if bytes.Equal(ntResponse, peerResponse) {
			recvKey, err := rfc3079.MakeKey(ntResponse, bytePwd, false)
			if err != nil {
				return fmt.Errorf("user:%s mschap access cannot make recvKey", username)
			}

			sendKey, err := rfc3079.MakeKey(ntResponse, bytePwd, true)
			if err != nil {
				return fmt.Errorf("user:%s mschap access cannot make sendKey", username)
			}

			authenticatorResponse, err := rfc2759.GenerateAuthenticatorResponse(challenge, peerChallenge, ntResponse, byteUser, bytePwd)
			if err != nil {
				return fmt.Errorf("user:%s mschap access  cannot generate authenticator response", username)

			}

			success := make([]byte, 43)
			success[0] = ident
			copy(success[1:], authenticatorResponse)

			microsoft.MSCHAP2Success_Add(radAccept, []byte(success))
			microsoft.MSMPPERecvKey_Add(radAccept, recvKey)
			microsoft.MSMPPESendKey_Add(radAccept, sendKey)
			microsoft.MSMPPEEncryptionPolicy_Add(radAccept, microsoft.MSMPPEEncryptionPolicy_Value_EncryptionAllowed)
			microsoft.MSMPPEEncryptionTypes_Add(radAccept, microsoft.MSMPPEEncryptionTypes_Value_RC440or128BitAllowed)
			radlog.Infof("user:%s mschap access accept", username)
		}
	}
	return fmt.Errorf("user:%s mschap access reject challenge len or response len error", username)

}
