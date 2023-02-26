package supervise

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ca17/teamsacs/app"
	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/cwmp"
	"github.com/ca17/teamsacs/common/timeutil"
	"github.com/ca17/teamsacs/common/web"
	"github.com/ca17/teamsacs/events"
	"github.com/ca17/teamsacs/models"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
)

func execCwmpUpdateFirmware(c echo.Context, devids []string, firmwareId, session string) error {
	var devs []models.NetCpe
	common.Must(app.GDB().Where("id in ?", devids).First(&devs).Error)

	if len(devs) > 100 {
		return c.JSON(http.StatusOK, web.RestError(fmt.Sprintf(
			"The maximum number of devices sent each time cannot exceed 100")))
	}

	var firmwareCfg models.CwmpFirmwareConfig
	err := app.GDB().Where("id=?", firmwareId).First(&firmwareCfg).Error
	if err != nil {
		return c.JSON(http.StatusOK, web.RestError(fmt.Sprintf("No firmware configuration found")))
	}

	if firmwareCfg.Content == "" {
		return c.JSON(http.StatusOK, web.RestError(fmt.Sprintf("The firmware configuration content is empty")))
	}

	for _, dev := range devs {
		if common.IsEmptyOrNA(dev.Sn) {
			return c.JSON(http.StatusOK, web.RestError(fmt.Sprintf("Device SN %s invalid", dev.Sn)))
		}

		cpe := app.GApp().CwmpTable().GetCwmpCpe(dev.Sn)
		if !app.GApp().MatchDevice(dev, firmwareCfg.Oui, firmwareCfg.ProductClass, firmwareCfg.SoftwareVersion) {
			events.PubSuperviseLog(dev.ID, session, "error",
				fmt.Sprintf("cpe %s not match CwmpFirmwareConfig", dev.Sn))
			continue
		}

		go func(devitem models.NetCpe) {
			scontent := app.GApp().InjectCwmpConfigVars(dev.Sn, firmwareCfg.Content, nil)
			// 创建脚本下发记录
			// session, _ := common.UUIDBase32()
			scriptSession := &models.CwmpConfigSession{
				ID:              common.UUIDint64(),
				ConfigId:        cast.ToString(firmwareCfg.ID),
				CpeId:           devitem.ID,
				Session:         session,
				Name:            "CwmpUpdateFirmware",
				Level:           "manage",
				SoftwareVersion: firmwareCfg.SoftwareVersion,
				ProductClass:    firmwareCfg.ProductClass,
				Oui:             firmwareCfg.Oui,
				Content:         scontent,
				ExecStatus:      "initialize",
				LastError:       "",
				Timeout:         120,
				ExecTime:        time.Now(),
				RespTime:        timeutil.EmptyTime,
				CreatedAt:       time.Time{},
				UpdatedAt:       time.Time{},
			}
			common.Must(app.GDB().Create(scriptSession).Error)

			// 文件下载 token 当日有效
			var token = common.Md5Hash(session + app.GConfig().Tr069.Secret + time.Now().Format("20060102"))
			err := cpe.SendCwmpEventData(models.CwmpEventData{
				Session: session,
				Sn:      devitem.Sn,
				Message: &cwmp.Download{
					ID:         session,
					Name:       "Cwmp FirmwareConfig Task",
					NoMore:     0,
					CommandKey: session,
					FileType:   "1 Firmware Upgrade Image",
					URL: fmt.Sprintf("%s/cwmpfiles/%s/%s/latest.xml",
						app.GApp().GetTr069SettingsStringValue(app.ConfigTR069AccessAddress), session, token),
					Username:       "",
					Password:       "",
					FileSize:       len([]byte(scontent)),
					TargetFileName: session + ".xml",
					DelaySeconds:   5,
					SuccessURL:     "",
					FailureURL:     "",
				},
			}, 5000, true)
			if err != nil {
				events.PubSuperviseLog(devitem.ID, session, "error",
					fmt.Sprintf("TR069 Push firmware configuration timed out %s", err.Error()))
				return
			}

			go connectDeviceAuth(session, dev)

		}(dev)
	}

	return c.JSON(200, web.RestSucc("The command to push the firmware has been sent, please check the management log later, please do not repeat it in a short time"))

}
