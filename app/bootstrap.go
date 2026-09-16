package app

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/zaplog/log"
	"github.com/ca17/teamsacs/models"
	"gorm.io/gorm"
)

const (
	defaultBootstrapAdminUsername = "admin"
	defaultBootstrapAdminPassword = "teamsacs"
)

// Bootstrap runs idempotent post-migration seeding.
// Safe for existing production databases: never overwrites an existing operator
// (including password), and only creates the default admin when that username
// is missing.
func (a *Application) Bootstrap() error {
	if err := a.ensureBootstrapAdmin(); err != nil {
		return err
	}
	if err := a.ensureSettings(); err != nil {
		return err
	}
	if err := a.ensureDefaultPNode(); err != nil {
		return err
	}
	a.checkAppVersion()
	return nil
}

func bootstrapAdminUsername() string {
	if v := strings.TrimSpace(os.Getenv("TEAMSACS_BOOTSTRAP_ADMIN_USERNAME")); v != "" {
		return v
	}
	return defaultBootstrapAdminUsername
}

func bootstrapAdminPassword() (password string, fromEnv bool) {
	if v := os.Getenv("TEAMSACS_BOOTSTRAP_ADMIN_PASSWORD"); v != "" {
		return v, true
	}
	return defaultBootstrapAdminPassword, false
}

func (a *Application) ensureBootstrapAdmin() error {
	username := bootstrapAdminUsername()
	password, fromEnv := bootstrapAdminPassword()

	var existing models.SysOpr
	err := a.gormDB.Where("username = ?", username).First(&existing).Error
	switch {
	case err == nil:
		// Existing production / already-bootstrapped install: do not touch.
		log.Infof("bootstrap: operator %q already exists, skip seeding", username)
		return nil
	case !errors.Is(err, gorm.ErrRecordNotFound):
		return err
	}

	opr := &models.SysOpr{
		ID:        common.UUIDint64(),
		Realname:  "administrator",
		Mobile:    "0000",
		Email:     "N/A",
		Username:  username,
		Password:  common.Sha256HashWithSalt(password, common.SecretSalt),
		Level:     "super",
		Status:    "enabled",
		Remark:    "bootstrap",
		LastLogin: time.Now(),
	}
	if err := a.gormDB.Create(opr).Error; err != nil {
		return err
	}

	if fromEnv {
		log.Infof("bootstrap: created operator %q with TEAMSACS_BOOTSTRAP_ADMIN_PASSWORD", username)
	} else {
		log.Warnf("bootstrap: created operator %q with default password %q; change it immediately after first login",
			username, defaultBootstrapAdminPassword)
	}
	return nil
}

func (a *Application) ensureSettings() error {
	var checkConfig = func(sortid int, stype, cname, value, remark string) error {
		var count int64
		if err := a.gormDB.Model(&models.SysConfig{}).Where("type = ? and name = ?", stype, cname).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return nil
		}
		return a.gormDB.Create(&models.SysConfig{
			ID:     0,
			Sort:   sortid,
			Type:   stype,
			Name:   cname,
			Value:  value,
			Remark: remark,
		}).Error
	}

	for sortid, name := range ConfigConstants {
		var err error
		switch name {
		case ConfigSystemTitle:
			err = checkConfig(sortid, "system", ConfigSystemTitle, "TeamsACS Management System", "System title")
		case ConfigSystemTheme:
			err = checkConfig(sortid, "system", ConfigSystemTheme, "light", "System theme")
		case ConfigSystemLoginRemark:
			err = checkConfig(sortid, "system", ConfigSystemLoginRemark, "Recommended browser: Chrome/Edge", "Login page description")
		case ConfigSystemLoginSubtitle:
			err = checkConfig(sortid, "system", ConfigSystemLoginSubtitle, "TeamsACS Community Edition", "Login form title")
		case ConfigCpeAutoRegister:
			err = checkConfig(sortid, "tr069", ConfigCpeAutoRegister, "enabled", "Auto register CPE device")
		case ConfigTR069AccessAddress:
			err = checkConfig(sortid, "tr069", ConfigTR069AccessAddress, "http://127.0.0.1:2999", "Teamsacs TR069 access address, HTTP | https://domain:port")
		case ConfigTR069AccessPassword:
			err = checkConfig(sortid, "tr069", ConfigTR069AccessPassword, "teamsacstr069password", "Teamsacs TR069 access password, It is provided to CPE to access TeamsACS")
		case ConfigCpeConnectionRequestPassword:
			err = checkConfig(sortid, "tr069", ConfigCpeConnectionRequestPassword, "teamsacscpepassword", "CPE Connection authentication password, It is provided to TeamsACS to access CPE")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Application) ensureDefaultPNode() error {
	var pnode models.NetNode
	err := a.gormDB.Where("id=?", AutoRegisterPopNodeId).First(&pnode).Error
	switch {
	case err == nil:
		return nil
	case !errors.Is(err, gorm.ErrRecordNotFound):
		return err
	}
	return a.gormDB.Create(&models.NetNode{
		ID:     AutoRegisterPopNodeId,
		Name:   "default",
		Remark: "Device auto-registration node",
	}).Error
}
