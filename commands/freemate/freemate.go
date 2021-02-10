package main

import (
	sha256_ "crypto/sha256"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
	h = flag.Bool("h", false, "help usage")
	x = flag.Bool("X", false, "debug")
	t = flag.String("t", "bs2radiuis", "api token")
)

func startFreeradius() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if err := exec.Command("chown freerad:freerad -R /var/log/freeradius").Run(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("start freeradius ")
	cmd := exec.Command("/usr/sbin/freeradius")
	if *x {
		cmd = exec.Command("/usr/sbin/freeradius", "-X")
	}
	cmd.Stdin = os.Stderr
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	// cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
	}
}

func startCheckProc() {
	ticker := time.NewTicker(time.Millisecond * 5000)
	go func() {
		for t := range ticker.C {
			_ = t.String()
			ps, _ := process.Processes()
			count := 0
			for _, p := range ps {
				name, _ := p.Name()
				st, _ := p.Status()
				if strings.Contains(name, "freeradius") {
					// fmt.Println(fmt.Sprintf("%s %s", name, st))
					if st == "Z" {
						fmt.Println(fmt.Sprintf("%s %s", name, st))
						// syscall.Kill(int(p.Pid), syscall.SIGKILL)
						p.Resume()
					}
					time.Sleep(time.Second * 3)
					if st, _ := p.Status(); st == "S" {
						count += 1
					}
				}
			}
			if count == 0 {
				go startFreeradius()
			}
		}
	}()
}


func KillRadiusProc() {
	ps, _ := process.Processes()
	for _, p := range ps {
		name, _ := p.Name()
		if strings.Contains(name, "freeradius") {
			syscall.Kill(int(p.Pid), syscall.SIGKILL)
		}
	}
}

func main() {
	flag.Parse()

	if *h == true {
		ustr := "daemon version: daemon/1.0, Usage:\ndaemon -h\nOptions:"
		fmt.Fprintf(os.Stderr, ustr)
		flag.PrintDefaults()
		return
	}

	startCheckProc()

	g.Go(func() error {
		return startApi()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

}

func Sha256HashWithSalt(src string, salt string) string {
	h := sha256_.New()
	h.Write([]byte(src))
	h.Write([]byte(salt))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

type ApiRequest struct {
	Sign string `json:"sign" query:"sign" form:"sign"`
	Data string `json:"data" query:"data" form:"data"`
}

// Handler
func clientUpdate(c echo.Context) error {
	req := new(ApiRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if Sha256HashWithSalt(req.Data, *t) != req.Sign {
		return c.String(http.StatusForbidden, "Reject")
	}
	fmt.Println(req.Data)
	err := ioutil.WriteFile("/etc/freeradius/clients.conf", []byte(req.Data), 0644)
	if err != nil {
		return c.String(http.StatusOK, "Failure")
	}
	go KillRadiusProc()
	return c.String(http.StatusOK, "Success")
}

func startApi() error {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HideBanner = true
	e.POST("/client/update", clientUpdate)
	return e.Start(":1815")
}
