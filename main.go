package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	_ "time/tzdata"

	"github.com/ca17/teamsacs/app"
	"github.com/ca17/teamsacs/assets"
	"github.com/ca17/teamsacs/common/zaplog/log"
	"github.com/ca17/teamsacs/config"
	"github.com/ca17/teamsacs/controllers"
	"github.com/ca17/teamsacs/installer"
	"github.com/ca17/teamsacs/tr069"
	"github.com/ca17/teamsacs/webserver"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

// 命令行定义
var (
	h         = flag.Bool("h", false, "help usage")
	showVer   = flag.Bool("v", false, "show version")
	conffile  = flag.String("c", "", "config yaml file")
	install   = flag.Bool("install", false, "run install")
	uninstall = flag.Bool("uninstall", false, "run uninstall")
)

// PrintVersion Print version information
func PrintVersion() {
	println(assets.BuildInfo)
}

func printHelp() {
	if *h {
		ustr := fmt.Sprintf("teamsacs version: %s, Usage: teamsacs -h\nOptions:", assets.BuildVersion())
		_, _ = fmt.Fprintf(os.Stderr, ustr)
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	if *showVer {
		PrintVersion()
		os.Exit(0)
	}

	printHelp()

	_config := config.LoadConfig(*conffile)

	// Install as a system service
	if *install {
		err := installer.Install()
		if err != nil {
			log.Error(err)
		}
		return
	}

	// 卸载
	if *uninstall {
		installer.Uninstall()
		return
	}

	app.InitGlobalApplication(_config)

	app.GApp().MigrateDB(false)

	defer app.Release()

	// 管理服务启动
	g.Go(func() error {
		webserver.Init()
		controllers.Init()
		return webserver.Listen()
	})

	g.Go(func() error {
		log.Info("Start tr069 server...")
		return tr069.Listen()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
