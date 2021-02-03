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

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"
	_ "time/tzdata"

	"github.com/op/go-logging"
	"golang.org/x/sync/errgroup"

	"github.com/ca17/teamsacs/common/installer"
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/config"
	"github.com/ca17/teamsacs/freeradius"
	"github.com/ca17/teamsacs/grpcservice"
	"github.com/ca17/teamsacs/message"
	"github.com/ca17/teamsacs/models"
	"github.com/ca17/teamsacs/nbi"
	"github.com/ca17/teamsacs/radiusd"
	"github.com/ca17/teamsacs/radiusd/radlog"
	"github.com/ca17/teamsacs/scheduler"
	"github.com/ca17/teamsacs/syslogd"
)

var (
	g errgroup.Group

	BuildVersion   string
	ReleaseVersion string
	BuildTime      string
	BuildName      string
	CommitID       string
	CommitDate     string
	CommitUser     string
	CommitSubject  string
)

//go:generate esc -o common/resources/resources.go -pkg resources -ignore=".DS_Store" resources
//go:generate protoc -I ./grpcservice --go_out=plugins=grpc:./grpcservice  ./grpcservice/service.proto

// Command line definition
var (
	h               = flag.Bool("h", false, "help usage")
	showVer         = flag.Bool("v", false, "show version")
	debug           = flag.Bool("X", false, "run debug level")
	syslogaddr      = flag.String("syslog", "", "syslog addr x.x.x.x:x")
	conffile        = flag.String("c", "/etc/teamsacs.yaml", "config yaml/json file")
	dev             = flag.Bool("dev", false, "run develop mode")
	port            = flag.Int("p", 0, "web port")
	install         = flag.Bool("install", false, "run install")
	startRadius     = flag.Bool("radiusd", false, "run radius server")
	startFreeradius = flag.Bool("freeradius-api", true, "run freeradius api")
	startNbi        = flag.Bool("nbi-service", true, "run northbound interface api")
	startRfc3164    = flag.Bool("syslog-rfc3164", false, "run rfc3164 syslog server")
	startRfc5424    = flag.Bool("syslog-rfc5424", false, "run rfc5424 syslog server")
	uninstall       = flag.Bool("uninstall", false, "run uninstall")
	initcfg         = flag.Bool("initcfg", false, "write default config > /etc/teamsacs.yaml")
	initAdmin       = flag.Bool("init-admin", false, "init admin api secret")
	adminName       = flag.String("admin", "admin", "init admin api secret with username ")
	messaged        = flag.Bool("messaged", true, "listen message pub server")
)

// Print version information
func PrintVersion() {
	fmt.Fprintf(os.Stdout, "build name:\t%s\n", BuildName)
	fmt.Fprintf(os.Stdout, "build version:\t%s\n", BuildVersion)
	fmt.Fprintf(os.Stdout, "build time:\t%s\n", BuildTime)
	fmt.Fprintf(os.Stdout, "release version:\t%s\n", ReleaseVersion)
	fmt.Fprintf(os.Stdout, "Commit ID:\t%s\n", CommitID)
	fmt.Fprintf(os.Stdout, "Commit Date:\t%s\n", CommitDate)
	fmt.Fprintf(os.Stdout, "Commit Username:\t%s\n", CommitUser)
	fmt.Fprintf(os.Stdout, "Commit Subject:\t%s\n", CommitSubject)
}

func printHelp() {
	if *h {
		ustr := fmt.Sprintf("%s version: %s, Usage:%s -h\nOptions:", BuildName, BuildVersion, BuildName)
		fmt.Fprintf(os.Stderr, ustr)
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func setupAppconfig() *config.AppConfig {
	appconfig := config.LoadConfig(*conffile)
	if *port > 0 {
		appconfig.NBI.Port = *port
	}

	if *syslogaddr != "" {
		appconfig.System.SyslogAddr = *syslogaddr
	}
	if *debug {
		appconfig.NBI.Debug = *debug
		appconfig.Radiusd.Debug = *debug
		appconfig.Grpc.Debug = *debug
	}
	appconfig.InitDirs()
	return appconfig
}

func setupLogging(appconfig *config.AppConfig) {
	// system logging
	level := logging.INFO
	if appconfig.NBI.Debug {
		level = logging.DEBUG
	}
	log.SetupLog(level, appconfig.System.SyslogAddr, appconfig.GetLogDir(), appconfig.System.Appid)

	// radius logging
	radlevel := logging.INFO
	if appconfig.Radiusd.Debug {
		radlevel = logging.DEBUG
	}
	radlog.SetupLog(radlevel, appconfig.System.SyslogAddr, appconfig.GetLogDir(), "Radiusd")

}

func installService(appconfig *config.AppConfig) bool {
	// 安装为系统服务
	if *install {
		err := installer.Install(appconfig)
		if err != nil {
			log.Error(err)
		}
		return true
	}

	// 卸载
	if *uninstall {
		installer.Uninstall()
		return true
	}
	return false
}

func ionitConfig(appconfig *config.AppConfig) bool {
	if *initcfg {
		err := installer.InitConfig(appconfig)
		if err != nil {
			log.Error(err)
		}
		return true
	}
	return false
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	if *showVer {
		PrintVersion()
		os.Exit(0)
	}

	printHelp()

	appconfig := setupAppconfig()

	// set logging level
	setupLogging(appconfig)

	if installService(appconfig) {
		return
	}

	if ionitConfig(appconfig) {
		return
	}

	manager := models.NewModelManager(appconfig, *dev)

	if *initAdmin && *adminName != "" {
		apisecret, err := manager.GetOpsManager().InitSuper(*adminName)
		if err != nil {
			fmt.Fprintln(os.Stdout, err.Error())
			os.Exit(0)
		}
		fmt.Fprintln(os.Stdout, *adminName+" "+apisecret)
		return
	}

	if *dev {
		log.Debug("Running for Dev Mode")
	}

	if *startRadius {
		g.Go(func() error {
			log.Info("Start Radius auth Server ...")
			return radiusd.ListenRadiusAuthServer(manager)
		})

		g.Go(func() error {
			log.Info("Start Radius acct Server ...")
			return radiusd.ListenRadiusAcctServer(manager)
		})
	}

	time.Sleep(time.Millisecond * 50)

	g.Go(func() error {
		log.Info("Start Grpc Server ...")
		return grpcservice.StartGrpcServer(manager)
	})

	if *startFreeradius {
		g.Go(func() error {
			log.Info("Start FreeRADIUS API Server ...")
			return freeradius.ListenFreeRADIUSServer(manager)
		})
	}

	time.Sleep(time.Millisecond * 50)

	if *startNbi {
		g.Go(func() error {
			log.Info("Start NBI Server ...")
			return nbi.ListenNBIServer(manager)
		})
	}

	syslogserv := syslogd.NewSyslogServer(manager)
	if *startRfc3164 || os.Getenv("TEAMSACS_RFC3164") == "true" {
		g.Go(func() error {
			log.Info("Start rfc3164 Syslog Server ...")
			return syslogserv.StartRfc3164()
		})
	}

	if *startRfc5424 || os.Getenv("TEAMSACS_RFC5424") == "true" {
		g.Go(func() error {
			log.Info("Start rfc5424 Syslog Server ...")
			return syslogserv.StartRfc5424()
		})
	}

	g.Go(func() error {
		log.Info("Start Syslog Server ...")
		return syslogserv.StartTextlog()
	})

	g.Go(func() error {
		return scheduler.Start(manager)
	})

	time.Sleep(time.Millisecond * 50)

	if *messaged {
		g.Go(func() error {
			log.Info("Start Message publish Server ...")
			return message.NewPubSubService(manager).StartPubServer()
		})
	}

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
