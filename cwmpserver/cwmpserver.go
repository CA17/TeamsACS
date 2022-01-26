package cwmpserver

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/ca17/teamsacs/app"
	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	elog "github.com/labstack/gommon/log"
)

var server *CwmpServer

const CwmpCookieName = "tr069_cookie"

type CwmpServer struct {
	root     *echo.Echo
	sesslock sync.Mutex
}

func Listen() error {
	server = NewTr069Server()
	server.initRouter()
	return server.Start()
}

func NewTr069Server() *CwmpServer {
	s := new(CwmpServer)
	s.root = echo.New()
	s.sesslock = sync.Mutex{}
	s.root.Pre(middleware.RemoveTrailingSlash())
	s.root.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					if app.Config.Cwmp.Debug {
						log.Errorf("%+v", r)
					}
					c.Error(echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
				}
			}()
			return next(c)
		}
	})
	s.root.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "cwmp-acs ${time_rfc3339} ${remote_ip} ${method} ${uri} ${protocol} ${status} ${id} ${user_agent} ${latency} ${bytes_in} ${bytes_out} ${error}\n",
		Output: os.Stdout,
	}))
	// s.root.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	// 	log.Info(string(resBody))
	// 	log.Info(string(resBody))
	// }))
	// s.root.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
	// 	Skipper: func(c echo.Context) bool {
	// 		rpath := c.Request().RequestURI
	// 		if strings.HasPrefix(rpath, "/cwmp/files") {
	// 			return true
	// 		} else {
	// 			return false
	// 		}
	// 	},
	// 	Validator: func(username, password string, c echo.Context) (bool, error) {
	// 		if username == "" {
	// 			return false, nil
	// 		}
	// 		return true, nil
	// 	},
	// 	Realm: "Restricted",
	// }))
	s.root.IPExtractor = echo.ExtractIPFromRealIPHeader()
	// s.root.Use(session.Middleware(sessions.NewCookieStore([]byte(app.Config.Web.Secret))))
	s.root.HideBanner = true
	s.root.Logger.SetLevel(common.If(app.Config.Cwmp.Debug, elog.DEBUG, elog.INFO).(elog.Lvl))
	s.root.Debug = app.Config.Cwmp.Debug
	return s
}

// Start 启动服务器
func (s *CwmpServer) Start() (err error) {
	log.Infof("启动 Tr069 API 服务器 %s:%d", app.Config.Cwmp.Host, app.Config.Cwmp.Port)
	if app.Config.Cwmp.Tls {
		err = s.root.StartTLS(fmt.Sprintf("%s:%d", app.Config.Cwmp.Host, app.Config.Cwmp.Port),
			path.Join(app.Config.GetPrivateDir(), "metaslink.tls.crt"), path.Join(app.Config.GetPrivateDir(), "metaslink.tls.key"))
	} else {
		err = s.root.Start(fmt.Sprintf("%s:%d", app.Config.Cwmp.Host, app.Config.Cwmp.Port))
	}
	if err != nil {
		log.Errorf("启动 Cwmp API 服务器错误 %s", err.Error())
	}
	return err
}

func (s *CwmpServer) GetLatestCookieSn(c echo.Context) string {
	cookie, err := c.Cookie(CwmpCookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func (s *CwmpServer) SetLatestInformByCookie(c echo.Context, sn string) {
	cookie := new(http.Cookie)
	cookie.Name = CwmpCookieName
	cookie.Value = sn
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)
}
