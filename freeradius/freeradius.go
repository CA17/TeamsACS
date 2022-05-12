package freeradius

import (
	"fmt"
	"time"

	"github.com/ca17/teamsacs/app"
	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/common/web"
	"github.com/ca17/teamsacs/models"
)

// freeradius api 接口服务

func getOnlineCount(username string) (int64, error) {
	var count int64
	err := app.GormDB.Model(&models.RadiusOnline{}).
		Where("username = ?", username).Count(&count).Error
	return count, err
}

func getOnlineCountBySessionid(acctSessionId string) (int64, error) {
	var count int64
	err := app.GormDB.Model(&models.RadiusOnline{}).
		Where("acct_session_id = ?", acctSessionId).Count(&count).Error
	return count, err
}

// func QueryAccountings(params web.RequestParams) (*web.PageResult, error) {
//	return m.QueryPagerItems(params, TeamsacsAccounting)
// }

// func QueryAuthlogs(params web.RequestParams) (*web.PageResult, error) {
//	return m.QueryPagerItems(params, TeamsacsAuthlog)
// }

// func QueryOnlines(params web.RequestParams) (*web.PageResult, error) {
//	return m.QueryPagerItems(params, TeamsacsOnline)
// }

func addRadiusAuthLog(username string, nasip string, result string, reason string, cast int64) error {
	return app.GormDB.Create(&models.RadiusAuthlog{
		ID:        common.UUIDint64(),
		Username:  username,
		NasAddr:   nasip,
		Result:    result,
		Reason:    reason,
		Cast:      int(cast),
		Timestamp: time.Now(),
	}).Error
}

func BatchClearRadiusOnlineDataByNas(nasip, nasid string) error {
	return app.GormDB.Where("nas_addr = ? or nas_id = ?", nasip, nasid).Delete(models.RadiusOnline{}).Error
}

// func AddRadiusAccounting(acct Accounting) error {
//	acct.AcctStopTime = time.Now()
//	_, err := m.GetTeamsAcsCollection(TeamsacsAccounting).InsertOne(context.TODO(), acct)
//	return err
// }
//
// func DeleteRadiusOnline(sessionid string) error {
//	_, err := m.GetTeamsAcsCollection(TeamsacsOnline).DeleteOne(context.TODO(), bson.M{"acct_session_id": sessionid})
//	return err
// }

// func BatchDeleteRadiusOnline(sessionids string) error {
// 	_, err := m.GetTeamsAcsCollection(TeamsacsOnline).DeleteMany(context.TODO(),
// 		bson.M{"acct_session_id": bson.M{"$in": strings.Split(sessionids, ",")}})
// 	return err
// }

// func BatchDeleteRadiusOnline(sessionids string) error {
//	for _, sid := range strings.Split(sessionids, ",") {
//		m.GetTeamsAcsCollection(TeamsacsOnline).DeleteOne(context.TODO(), bson.M{"acct_session_id": sid})
//	}
//	return nil
// }
//
// func TruncateRadiusOnline() error {
//	_, err := m.GetTeamsAcsCollection(TeamsacsOnline).DeleteMany(context.TODO(), bson.M{})
//	return err
// }
//
// func UpdateRadiusOnlineData(acct Accounting) error {
//	query := bson.M{"_id": acct.AcctSessionId}
//	acct.ID = ""
//	_, err := m.GetTeamsAcsCollection(TeamsacsOnline).UpdateOne(context.TODO(), query, bson.M{"$set": acct})
//	return err
// }

func getAcctStartTime(sessionTime string) time.Time {
	m, _ := time.ParseDuration("-" + sessionTime + "s")
	return time.Now().Add(m)
}

func getInputTotal(form *web.WebForm) int64 {
	var acctInputOctets = form.GetInt64Val("acctInputOctets", 0)
	var acctInputGigawords = form.GetInt64Val("acctInputGigaword", 0)
	return acctInputOctets + acctInputGigawords*4*1024*1024*1024
}

func getOutputTotal(form *web.WebForm) int64 {
	var acctOutputOctets = form.GetInt64Val("acctOutputOctets", 0)
	var acctOutputGigawords = form.GetInt64Val("acctOutputGigawords", 0)
	return acctOutputOctets + acctOutputGigawords*4*1024*1024*1024
}

// updateRadiusOnline 更新记账信息
func updateRadiusOnline(form *web.WebForm) error {
	var user models.RadiusUser
	err := app.GormDB.Where("username=?", form.GetVal("username")).First(&user).Error
	if err != nil {
		return fmt.Errorf("user %s not exists", form.GetVal("username"))
	}
	var sessionId = form.GetVal2("acctSessionId", "")
	var statusType = form.GetVal2("acctStatusType", "")
	radOnline := models.RadiusOnline{
		ID:                common.UUIDint64(),
		Username:          form.GetVal("username"),
		NasId:             form.GetVal("nasid"),
		NasAddr:           form.GetVal("nasip"),
		NasPaddr:          form.GetVal("nasip"),
		SessionTimeout:    form.GetIntVal("sessionTimeout", 0),
		FramedIpaddr:      form.GetVal2("framedIPAddress", "0.0.0.0"),
		FramedNetmask:     form.GetVal2("framedIPNetmask", common.NA),
		MacAddr:           form.GetVal2("macAddr", common.NA),
		NasPort:           0,
		NasClass:          common.NA,
		NasPortId:         form.GetVal2("nasPortId", common.NA),
		NasPortType:       0,
		ServiceType:       0,
		AcctSessionId:     sessionId,
		AcctSessionTime:   form.GetIntVal("acctSessionTime", 0),
		AcctInputTotal:    getInputTotal(form),
		AcctOutputTotal:   getOutputTotal(form),
		AcctInputPackets:  form.GetIntVal("acctInputPackets", 0),
		AcctOutputPackets: form.GetIntVal("acctOutputPackets", 0),
		AcctStartTime:     getAcctStartTime(form.GetVal2("acctSessionTime", "0")),
		LastUpdate:        time.Now(),
	}

	switch statusType {
	case "Start", "Update", "Alive", "Interim-Update":
		// public cpe update event
		ocount, _ := getOnlineCountBySessionid(sessionId)
		if ocount == 0 {
			log.Infof("Add radius online %+v", radOnline)
			return app.GormDB.Create(&radOnline).Error
		} else {
			log.Infof("Update radius online %+v", radOnline)
			return app.GormDB.Model(models.RadiusOnline{}).
				Where("acct_session_id=?", sessionId).
				Updates(radOnline).Error
		}
	case "Stop":
		log.Infof("Update radius cdr %+v", radOnline)
		app.GormDB.Create(&models.RadiusAccounting{
			ID:                radOnline.ID,
			Username:          radOnline.Username,
			NasId:             radOnline.NasId,
			NasAddr:           radOnline.NasAddr,
			NasPaddr:          radOnline.NasPaddr,
			SessionTimeout:    radOnline.SessionTimeout,
			FramedIpaddr:      radOnline.FramedIpaddr,
			FramedNetmask:     radOnline.FramedNetmask,
			MacAddr:           radOnline.MacAddr,
			NasPort:           radOnline.NasPort,
			NasClass:          radOnline.NasClass,
			NasPortId:         radOnline.NasPortId,
			NasPortType:       radOnline.NasPortType,
			ServiceType:       radOnline.ServiceType,
			AcctSessionId:     radOnline.AcctSessionId,
			AcctSessionTime:   radOnline.AcctSessionTime,
			AcctInputTotal:    radOnline.AcctInputTotal,
			AcctOutputTotal:   radOnline.AcctOutputTotal,
			AcctInputPackets:  radOnline.AcctInputPackets,
			AcctOutputPackets: radOnline.AcctOutputPackets,
			AcctStartTime:     radOnline.AcctStartTime,
			LastUpdate:        time.Now(),
			AcctStopTime:      time.Now(),
		})
		return app.GormDB.Where("acct_session_id=?", sessionId).
			Delete(&models.RadiusOnline{}).Error
	case "Accounting-On", "Accounting-Off":
	}

	return nil
}
