package cwmpserver

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ca17/teamsacs/app"
	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/cwmp"
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/models"
	"github.com/labstack/echo/v4"
)

func (s *CwmpServer) initRouter() {
	s.root.Add(http.MethodPost, "", s.Tr069Index)
	s.root.Add(http.MethodGet, "/cwmpfiles/:session/:token/:filename", s.Tr069Download)
}

// Tr069Download Cwmp 文件下载
func (s *CwmpServer) Tr069Download(c echo.Context) error {
	var session = c.Param("session")
	var token = c.Param("token")
	var filename = c.Param("filename")
	if session == "" || token == "" {
		return c.String(400, "bad request")
	}
	log.Info("cpe fetch cwmp file session = " + session)
	// 文件 token 当日有效
	if token != common.Md5Hash(session+app.Config.Web.Secret+time.Now().Format("20060102")) {
		return c.String(400, "bad token")
	}
	var downloadTask models.CwmpDownloadTask
	common.Must(app.GormDB.Where("session_id = ?", session).First(&downloadTask).Error)
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("Keep-Alive", "timeout=5")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", filename))
	return c.Blob(200, echo.MIMEOctetStream, []byte(downloadTask.Content))
}

func (s *CwmpServer) Tr069Index(c echo.Context) error {
	if app.Config.Cwmp.Debug {
		logRequestHeader(c)
	}
	cookie, _ := c.Cookie(CwmpCookieName)
	if cookie != nil {
		log.Infof("cwmp cooike session sn = %s", cookie.Value)
	}

	requestBody, err := ioutil.ReadAll(c.Request().Body)
	var bodyLen = len(requestBody)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("cwmp read error %s", err.Error()))
	}

	if app.Config.Cwmp.Debug {
		if bodyLen == 0 {
			log.Info("recv cpe empty message")
		} else {
			log.Info(string(requestBody))
		}
	}

	var msg cwmp.Message
	var lastInform *cwmp.Inform

	if bodyLen > 0 {
		msg, err = cwmp.ParseXML(requestBody)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("cwmp read xml error %s", err.Error()))
		}

		if msg.GetName() != "Inform" {
			lastestSn := s.GetLatestCookieSn(c)
			if lastestSn == "" {
				return c.String(http.StatusUnauthorized, "no cookie sn")
			}
			cpe := app.CwmpTable.GetCwmpCpe(lastestSn)
			if cpe.LastInform == nil {
				return c.String(http.StatusUnauthorized, "no cookie cpe data")
			}
			lastInform = cpe.LastInform
		}

		if app.Config.Cwmp.Debug {
			log.Info(common.ToJson(msg))
		}

		switch msg.GetName() {
		case "Inform":
			return s.handleInform(c, lastInform, msg)
		case "TransferComplete":
			return s.handleTransferComplete(c, msg)
		case "GetRPCMethods":
			return s.handleGetRPCMethods(c, msg)
		case "GetParameterValuesResponse":
			s.handleGetParameterValuesResponse(c, msg)
		case "GetParameterNamesResponse":
			s.handleGetParameterNamesResponse(c, msg)
		default:
			s.handleDefaultMesage(c, msg)
		}
	} else {
		lastestSn := s.GetLatestCookieSn(c)
		if lastestSn == "" {
			return noContentResp(c)
		}
		cpe := app.CwmpTable.GetCwmpCpe(lastestSn)

		msg, err := cpe.RecvCwmpEventData(1000, true)
		if err != nil {
			msg, _ = cpe.RecvCwmpEventData(1000, false)
		}

		if msg != nil {
			if msg.Session != "" {
				app.PubDeviceEventLog(lastestSn, "cwmp", msg.Session, msg.Message.GetName(), "info",
					fmt.Sprintf("Send Cwmp %s Message %s", msg.Message.GetName(), common.ToJson(msg.Message)))
			}
			return xmlCwmpMessage(c, msg.Message.CreateXML())
		}
	}

	return noContentResp(c)
}

// 未处理的消息
func (s *CwmpServer) handleDefaultMesage(c echo.Context, msg cwmp.Message) {
	log.Infof("未处理的消息类型 %s", msg.GetName())
	lastestSn := s.GetLatestCookieSn(c)
	if lastestSn != "" {
		app.PubDeviceEventLog(lastestSn, "cwmp", msg.GetID(), msg.GetName(), "info",
			fmt.Sprintf("Recv Cwmp %s Message %s", msg.GetName(), common.ToJson(msg)))
	}
}

// 处理 Inform
func (s *CwmpServer) handleInform(c echo.Context, lastInform *cwmp.Inform, msg cwmp.Message) error {
	lastInform = msg.(*cwmp.Inform)
	log.Infof("Inform Events: %v", lastInform.Events)
	s.SetLatestInformByCookie(c, lastInform.Sn)
	cpe := app.CwmpTable.GetCwmpCpe(lastInform.Sn)
	cpe.CheckRegister(c.RealIP(), lastInform)
	cpe.UpdateStatus(lastInform)
	// 通知系统更新数据
	cpe.NotifyDataUpdate(false)
	// response
	resp := new(cwmp.InformResponse)
	resp.ID = lastInform.ID
	resp.MaxEnvelopes = lastInform.MaxEnvelopes
	response := resp.CreateXML()

	if lastInform.IsEvent(cwmp.EventBootStrap) || lastInform.IsEvent(cwmp.EventBoot) {
		log.ErrorIf(cpe.UpdateManagementAuthInfo("bootstrap-session-"+common.UUID(), 1000, false))
		log.ErrorIf(cpe.UpdateParameterNames("bootstrap-session-"+common.UUID(), 1000, false))
	}

	// log inform
	if lastInform.IsEvent(cwmp.EventValueChange) {
		cpe.NotifyDataUpdate(true)
	}
	return xmlCwmpMessage(c, response)
}

// 处理 TransferComplete
func (s *CwmpServer) handleTransferComplete(c echo.Context, msg cwmp.Message) error {
	tc := msg.(*cwmp.TransferComplete)
	// do something
	resp := new(cwmp.TransferCompleteResponse)
	resp.ID = tc.ID
	response := resp.CreateXML()
	go func() {
		if tc.CommandKey != "" {
			app.PubDeviceEventLog("", "cwmp", tc.CommandKey, msg.GetName(), "info",
				fmt.Sprintf("Recv Cwmp %s Message %s", msg.GetName(), common.ToJson(msg)))
			err := app.GormDB.Model(&models.CwmpDownloadTask{}).Where("session_id = ?", tc.CommandKey).Updates(map[string]interface{}{
				"status":        "complate",
				"fault_code":    tc.FaultCode,
				"fault_string":  tc.FaultString,
				"start_time":    tc.StartTime,
				"complete_time": tc.CompleteTime,
			}).Error
			if err != nil {
				log.Errorf("Update CwmpSession Status error %s", err.Error())
			}
		}
	}()
	return xmlCwmpMessage(c, response)
}

// 处理 GetRPCMethods
func (s *CwmpServer) handleGetRPCMethods(c echo.Context, msg cwmp.Message) error {
	gm := msg.(*cwmp.GetRPCMethods)
	resp := new(cwmp.GetRPCMethodsResponse)
	resp.ID = gm.ID
	response := resp.CreateXML()
	return xmlCwmpMessage(c, response)
}

// 处理 GetParameterValuesResponse
func (s *CwmpServer) handleGetParameterValuesResponse(c echo.Context, msg cwmp.Message) {
	gm := msg.(*cwmp.GetParameterValuesResponse)
	lastestSn := s.GetLatestCookieSn(c)
	if lastestSn != "" {
		app.CwmpTable.GetCwmpCpe(lastestSn).UpdateParamValues(gm.Values)
		app.PubDeviceEventLog(lastestSn, "cwmp", msg.GetID(), msg.GetName(), "info",
			fmt.Sprintf("Recv Cwmp %s Message %s", msg.GetName(), common.ToJson(msg)))
	}
}

// 处理 GetParameterNamesResponse
func (s *CwmpServer) handleGetParameterNamesResponse(c echo.Context, msg cwmp.Message) {
	gm := msg.(*cwmp.GetParameterNamesResponse)
	lastestSn := s.GetLatestCookieSn(c)
	if lastestSn != "" && msg.GetID() != "" {
		if strings.HasPrefix(msg.GetID(), "bootstrap-session") {
			go app.CwmpTable.GetCwmpCpe(lastestSn).ProcessParameterNamesResponse(gm)
		}
		if len(gm.Params) < 100 {
			app.PubDeviceEventLog(lastestSn, "cwmp", msg.GetID(), msg.GetName(), "info",
				fmt.Sprintf("Recv Cwmp %s Message %s", msg.GetName(), common.ToJson(msg)))
		} else {
			app.PubDeviceEventLog(lastestSn, "cwmp", msg.GetID(), msg.GetName(), "info",
				fmt.Sprintf("Recv Cwmp %s Message, names total %d", msg.GetName(), len(gm.Params)))
		}
	}
}

func xmlCwmpMessage(c echo.Context, response []byte) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationXMLCharsetUTF8)
	c.Response().Header().Set("Connection", "keep-alive")
	if app.Config.Cwmp.Debug {
		logResponseHeader(c)
	}
	log.Info(string(response))
	return c.XMLBlob(200, response)
}

func noContentResp(c echo.Context) error {
	c.Response().Header().Set("Connection", "close")
	if app.Config.Cwmp.Debug {
		logResponseHeader(c)
	}
	return c.NoContent(http.StatusNoContent)
}

func logRequestHeader(c echo.Context) {
	for k, v := range c.Request().Header {
		fmt.Printf("%s: %s\n", k, v)
	}
}

func logResponseHeader(c echo.Context) {
	for k, v := range c.Response().Header() {
		fmt.Printf("%s: %s\n", k, v)
	}
}
