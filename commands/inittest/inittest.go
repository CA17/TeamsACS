package main

import (
	"github.com/ca17/teamsacs/app"
	"github.com/ca17/teamsacs/config"
)

func main() {
	app.InitGlobalApplication(config.LoadConfig("../teamsacs.yml"))
	app.GApp().InitTest()
}
