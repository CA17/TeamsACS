/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/common/web"
	"github.com/ca17/teamsacs/constant"
)

type Authlog struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string    `bson:"username,omitempty" json:"username,omitempty"`
	NasAddr   string    `bson:"nas_addr,omitempty" json:"nas_addr,omitempty"`
	Cast      int       `bson:"cast,omitempty" json:"cast,omitempty"`
	Result    string    `bson:"result,omitempty" json:"result,omitempty"`
	Reason    string    `bson:"reason,omitempty" json:"reason,omitempty"`
	Timestamp time.Time `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
}

// Accounting
// Radius Accounting Recode
type Accounting struct {
	ID                string    `bson:"_id,omitempty" json:"id,omitempty"`
	Username          string    `bson:"username" json:"username"`
	NasId             string    `bson:"nas_id" json:"nas_id"`
	NasAddr           string    `bson:"nas_addr" json:"nas_addr"`
	NasPaddr          string    `bson:"nas_paddr" json:"nas_paddr"`
	SessionTimeout    int       `bson:"session_timeout" json:"session_timeout"`
	FramedIpaddr      string    `bson:"framed_ipaddr" json:"framed_ipaddr"`
	FramedNetmask     string    `bson:"framed_netmask" json:"framed_netmask"`
	MacAddr           string    `bson:"mac_addr" json:"mac_addr"`
	NasPort           int64     `bson:"nas_port" json:"nas_port,string"`
	NasClass          string    `bson:"nas_class" json:"nas_class"`
	NasPortId         string    `bson:"nas_port_id" json:"nas_port_id"`
	NasPortType       int       `bson:"nas_port_type" json:"nas_port_type"`
	ServiceType       int       `bson:"service_type" json:"service_type"`
	AcctSessionId     string    `bson:"acct_session_id" json:"acct_session_id"`
	AcctSessionTime   int       `bson:"acct_session_time" json:"acct_session_time"`
	AcctInputTotal    int64     `bson:"acct_input_total" json:"acct_input_total,string"`
	AcctOutputTotal   int64     `bson:"acct_output_total" json:"acct_output_total,string"`
	AcctInputPackets  int       `bson:"acct_input_packets" json:"acct_input_packets"`
	AcctOutputPackets int       `bson:"acct_output_packets" json:"acct_output_packets"`
	AcctStartTime     time.Time `bson:"acct_start_time" json:"acct_start_time"`
	LastUpdate        time.Time `bson:"last_update" json:"last_update"`
	AcctStopTime      time.Time `bson:"acct_stop_time" json:"acct_stop_time"`
}

type RadiusManager struct{ *ModelManager }

func (m *ModelManager) GetRadiusManager() *RadiusManager {
	store, _ := m.ManagerMap.Get("RadiusManager")
	return store.(*RadiusManager)
}

func (m *RadiusManager) GetOnlineCount(username string) (int64, error) {
	coll := m.GetTeamsAcsCollection(TeamsacsOnline)
	return coll.CountDocuments(context.TODO(), bson.M{"username": username})
}

func (m *RadiusManager) GetOnlineCountBySessionid(acct_session_id string) (int64, error) {
	coll := m.GetTeamsAcsCollection(TeamsacsOnline)
	return coll.CountDocuments(context.TODO(), bson.M{"acct_session_id": acct_session_id})
}

func (m *RadiusManager) QueryAccountings(params web.RequestParams) (*web.PageResult, error) {
	return m.QueryPagerItems(params, TeamsacsAccounting)
}

func (m *RadiusManager) QueryAuthlogs(params web.RequestParams) (*web.PageResult, error) {
	return m.QueryPagerItems(params, TeamsacsAuthlog)
}

func (m *RadiusManager) QueryOnlines(params web.RequestParams) (*web.PageResult, error) {
	return m.QueryPagerItems(params, TeamsacsOnline)
}

func (m *RadiusManager) AddRadiusAuthLog(username string, nasip string, result string, reason string, cast int64) error {
	authlog := Authlog{
		ID:        common.UUID(),
		Username:  username,
		NasAddr:   nasip,
		Result:    result,
		Reason:    reason,
		Cast:      int(cast),
		Timestamp: time.Now(),
	}
	coll := m.GetTeamsAcsCollection(TeamsacsAuthlog)
	_, err := coll.InsertOne(context.TODO(), authlog)
	return err
}

func (m *RadiusManager) BatchClearRadiusOnlineDataByNas(nasip, nasid string) error {
	coll := m.GetTeamsAcsCollection(TeamsacsOnline)
	filter := bson.D{
		{"$or",
			bson.A{
				bson.D{{"nas_addr", nasip}},
				bson.D{{"nas_id", nasid}},
			}},
	}
	_, err := coll.DeleteMany(context.TODO(), filter)
	return err
}

func (m *RadiusManager) AddRadiusOnline(ol Accounting) error {
	_, err := m.GetTeamsAcsCollection(TeamsacsOnline).InsertOne(context.TODO(), ol)
	return err
}

func (m *RadiusManager) AddRadiusAccounting(acct Accounting) error {
	acct.AcctStopTime = time.Now()
	_, err := m.GetTeamsAcsCollection(TeamsacsAccounting).InsertOne(context.TODO(), acct)
	return err
}

func (m *RadiusManager) DeleteRadiusOnline(sessionid string) error {
	_, err := m.GetTeamsAcsCollection(TeamsacsOnline).DeleteOne(context.TODO(), bson.M{"acct_session_id": sessionid})
	return err
}

// func (m *RadiusManager) BatchDeleteRadiusOnline(sessionids string) error {
// 	_, err := m.GetTeamsAcsCollection(TeamsacsOnline).DeleteMany(context.TODO(),
// 		bson.M{"acct_session_id": bson.M{"$in": strings.Split(sessionids, ",")}})
// 	return err
// }

func (m *RadiusManager) BatchDeleteRadiusOnline(sessionids string) error {
	for _, sid := range strings.Split(sessionids, ",") {
		m.GetTeamsAcsCollection(TeamsacsOnline).DeleteOne(context.TODO(), bson.M{"acct_session_id": sid})
	}
	return nil
}

func (m *RadiusManager) TruncateRadiusOnline() error {
	_, err := m.GetTeamsAcsCollection(TeamsacsOnline).DeleteMany(context.TODO(), bson.M{})
	return err
}

func (m *RadiusManager) UpdateRadiusOnlineData(acct Accounting) error {
	query := bson.M{"_id": acct.AcctSessionId}
	acct.ID = ""
	_, err := m.GetTeamsAcsCollection(TeamsacsOnline).UpdateOne(context.TODO(), query, acct)
	return err
}

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

// 更新记账信息
func (m *RadiusManager) UpdateRadiusOnline(form *web.WebForm) error {
	user, err := m.GetSubscribeManager().GetSubscribeByUser(form.GetVal("username"))
	if err != nil {
		return fmt.Errorf("user %s not exists", form.GetVal("username"))
	}
	var sessionId = form.GetVal2("acctSessionId", "")
	var statusType = form.GetVal2("acctStatusType", "")
	radOnline := Accounting{
		ID: sessionId,
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
		if user.GetUserType() == "cpe" {
			m.Bus.Publish(constant.EventMosdnsCpeUpdate, radOnline.FramedIpaddr, radOnline.MacAddr)
		}

		ocount, _ := m.GetOnlineCountBySessionid(sessionId)
		if ocount == 0 {
			log.Infof("Add radius online %+v", radOnline)
			return m.AddRadiusOnline(radOnline)
		} else {
			log.Infof("Update radius online %+v", radOnline)
			return m.UpdateRadiusOnlineData(radOnline)
		}
	case "Stop":
		log.Infof("Update radius cdr %+v", radOnline)
		// public cpe update event
		if user.GetUserType() == "cpe" {
			m.Bus.Publish(constant.EventMosdnsCpeRemove, radOnline.FramedIpaddr)
		}
		_ = m.AddRadiusAccounting(radOnline)
		return m.DeleteRadiusOnline(sessionId)
	case "Accounting-On", "Accounting-Off":
		if user.GetUserType() == "cpe" {
			m.Bus.Publish(constant.EventMosdnsCpeClean, radOnline.FramedIpaddr)
		}
	}

	return nil
}

func (m *RadiusManager) ClearExpireOnline() (int, error) {
	filter := bson.M{"last_update": bson.M{"$lte": time.Now().Add(time.Second * 300 * -1)}}
	drr, err := m.GetTeamsAcsCollection(TeamsacsOnline).DeleteMany(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	return int(drr.DeletedCount), nil
}
