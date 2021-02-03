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

package installer

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/config"
)

var InstallScript = `#!/bin/bash -x
groupadd teamsacs
useradd teamsacs -g teamsacs -M -s /sbin/nologin
mkdir -p {{workdir}}
chown -R teamsacs.teamsacs {{workdir}}
chmod -R 700 {{workdir}}
install -m 777 ./teamsacs /usr/local/bin/teamsacs 
chown teamsacs.teamsacs /etc/teamsacs.yaml 
test -d /usr/lib/systemd/system || mkdir -p /usr/lib/systemd/system
cat>/usr/lib/systemd/system/teamsacs.service<<EOF
[Unit]
Description=teamsacs
After=network.target

[Service]
Environment=GODEBUG=x509ignoreCN=0
LimitNOFILE=65535
LimitNPROC=65535
Username=teamsacs
ExecStart={{command}}

[Install]
WantedBy=multi-user.target
EOF

chmod 600 /usr/lib/systemd/system/teamsacs.service
systemctl enable teamsacs && systemctl daemon-reload

`

func InitConfig(config *config.AppConfig) error {
	// config.NBI.JwtSecret = common.UUID()
	cfgstr, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("/etc/teamsacs.yaml", cfgstr, 0644)
}

func Install(config *config.AppConfig) error {
	if !common.FileExists("/etc/teamsacs.yaml") {
		_ = InitConfig(config)
	}
	script := strings.ReplaceAll(InstallScript, "{{workdir}}", config.System.Workdir)
	cmd := "/usr/local/bin/teamsacs"
	script = strings.ReplaceAll(InstallScript, "{{command}}", cmd)
	_ = ioutil.WriteFile("/tmp/teamsacs_install.sh", []byte(script), 0777)

	// 创建用户&组
	if err := exec.Command("/bin/bash", "/tmp/teamsacs_install.sh").Run(); err != nil {
		return err
	}
	return os.Remove("/tmp/teamsacs_install.sh")
}

func Uninstall() {
	_ = os.Remove("/etc/teamsacs.yaml")
	_ = os.Remove("/usr/lib/systemd/system/teamsacs.service")
	_ = os.Remove("/usr/local/bin/teamsacs")
}
