package apiserver

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ca17/teamsacs/app"
	"github.com/ca17/teamsacs/assets"
	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/common/tpl"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	elog "github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

type Skips struct {
	Paths   []string `json:"paths"`
	Prefixs []string `json:"prefixs"`
	Suffixs []string `json:"suffixs"`
}

var server *ApiServer

type ApiServer struct {
	root      *echo.Echo
	dbapi     *echo.Group
	pgapi     *echo.Group
	jwtConfig middleware.JWTConfig
}

func Init() {
	server = NewApiServer()
}

func Listen() error {
	return server.Start()
}

// NewApiServer 创建管理系统服务器
func NewApiServer() *ApiServer {
	s := &ApiServer{}
	s.root = echo.New()
	s.root.Pre(middleware.RemoveTrailingSlash())
	s.root.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 9,
	}))
	// 失败恢复处理中间件
	s.root.Use(ServerRecover(app.Config.System.Debug))
	// 日志处理中间件
	s.root.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: app.Config.System.Appid + " ${time_rfc3339} ${remote_ip} ${method} ${uri} ${protocol} ${status} ${id} ${user_agent} ${latency} ${bytes_in} ${bytes_out} ${error}\n",
		Output: os.Stdout,
	}))

	// 模板加载
	s.root.Renderer = tpl.NewCommonTemplate(assets.TemplatesFs, []string{"templates"}, map[string]interface{}{
		"pagever": func() int64 {
			return time.Now().Unix()
		},
		"buildver": func() string {
			return strings.TrimSpace(assets.BuildVer)
		},
	})

	s.root.HideBanner = true
	// 设置日志级别
	s.root.Logger.SetLevel(common.If(app.Config.System.Debug, elog.DEBUG, elog.INFO).(elog.Lvl))
	s.root.Debug = app.Config.System.Debug

	// JWT 中间件
	s.jwtConfig = middleware.JWTConfig{
		SigningKey:    []byte(app.Config.Web.Secret),
		SigningMethod: middleware.AlgorithmHS256,
		Skipper: func(c echo.Context) bool {
			if os.Getenv("TEAMSACS_DEVMODE") == "true" {
				return true
			}
			if common.InSlice(c.Request().RequestURI, []string{"/status"}) {
				return true
			}
			return false
		},
		ErrorHandler: func(err error) error {
			return echo.NewHTTPError(http.StatusBadRequest, "资源访问受限", err.Error())
		},
	}
	s.root.Use(middleware.JWTWithConfig(s.jwtConfig))
	s.initRouter()
	for _, r := range s.root.Routes() {
		log.Infof(fmt.Sprintf("%s %s", r.Method, r.Path))
	}
	return s
}

// Start 启动服务器
func (s *ApiServer) Start() (err error) {
	log.Infof("启动管理服务器 %s:%d", app.Config.Web.Host, app.Config.Web.Port)
	if app.Config.Web.Tls {
		err = s.root.StartTLS(fmt.Sprintf("%s:%d", app.Config.Web.Host, app.Config.Web.Port),
			path.Join(app.Config.GetPrivateDir(), "metaslink.tls.crt"), path.Join(app.Config.GetPrivateDir(), "metaslink.tls.key"))
	} else {
		err = s.root.Start(fmt.Sprintf("%s:%d", app.Config.Web.Host, app.Config.Web.Port))
	}
	if err != nil {
		log.Errorf("启启动管理服务器错误 %s", err.Error())
	}
	return err
}

// ParseJwtToken 解析 Jwt Token
func (s *ApiServer) ParseJwtToken(tokenstr string) (jwt.MapClaims, error) {
	config := s.jwtConfig
	token, err := jwt.Parse(tokenstr, func(t *jwt.Token) (interface{}, error) {
		// Check the signing method
		if t.Method.Alg() != config.SigningMethod {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		if len(config.SigningKeys) > 0 {
			if kid, ok := t.Header["kid"].(string); ok {
				if key, ok := config.SigningKeys[kid]; ok {
					return key, nil
				}
			}
			return nil, fmt.Errorf("unexpected jwt key id=%v", t.Header["kid"])
		}
		return config.SigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims := token.Claims.(jwt.MapClaims)
	return claims, err
}

// ServerRecover Web 服务恢复处理中间件
func ServerRecover(debug bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					if debug {
						log.Errorf("%+v", errors.WithStack(err))
					}
					c.Error(echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
				}
			}()
			return next(c)
		}
	}
}
