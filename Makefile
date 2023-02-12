BUILD_ORG   := ca17
BUILD_VERSION   := latest
BUILD_TIME      := $(shell date "+%F %T")
BUILD_NAME      := teamsacs
RELEASE_VERSION := v1.0.1
SOURCE          := main.go
RELEASE_DIR     := ./release
COMMIT_SHA1     := $(shell git show -s --format=%H )
COMMIT_DATE     := $(shell git show -s --format=%cD )
COMMIT_USER     := $(shell git show -s --format=%ce )
COMMIT_SUBJECT     := $(shell git show -s --format=%s )

buildpre:
	echo "${BUILD_VERSION} ${RELEASE_VERSION} ${BUILD_TIME}" > assets/buildver.txt
	echo "BuildVersion=${BUILD_VERSION}" > assets/build.txt
	echo "ReleaseVersion=${RELEASE_VERSION}" >> assets/build.txt
	echo "BuildTime=${BUILD_TIME}" >> assets/build.txt
	echo "BuildName=${BUILD_NAME}" >> assets/build.txt
	echo "CommitID=${COMMIT_SHA1}" >> assets/build.txt
	echo "CommitDate=${COMMIT_DATE}" >> assets/build.txt
	echo "CommitUser=${COMMIT_USER}" >> assets/build.txt
	echo "CommitSubject=${COMMIT_SUBJECT}" >> assets/build.txt

fastpub:
	make buildpre
	make build
	docker build --build-arg BTIME="$(shell date "+%F %T")" -t teamsacs . -f docker/teamsacs/Dockerfile
	docker tag teamsacs ${BUILD_ORG}/teamsacs:latest
	docker push ${BUILD_ORG}/teamsacs:latest

fastpubm1:
	make buildpre
	make build
	docker buildx build --platform=linux/amd64 --build-arg BTIME="$(shell date "+%F %T")" -t teamsacs . -f docker/teamsacs/Dockerfile
	docker tag teamsacs ${BUILD_ORG}/teamsacs:latest
	docker push ${BUILD_ORG}/teamsacs:latest

build:
	test -d ./release || mkdir -p ./release
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags  '-s -w -extldflags "-static"'  -o ./release/teamsacs main.go
	upx ./release/teamsacs


syncdev:
	@read -p "提示:同步操作尽量在完成一个完整功能特性后进行，请输入提交描述 (develop):  " cimsg; \
	make buildpre
	git commit -am "$(shell date "+%F %T") : $${cimsg}" || echo "no commit"
	# 切换主分支并更新
	git checkout main
	git pull origin main
	# 切换开发分支变基合并提交
	git checkout develop
	git rebase -i main
	# 切换回主分支并合并开发者分支，推送主分支到远程，方便其他开发者合并
	git checkout main
	git merge --no-ff develop
	git push origin main
	# 切换回自己的开发分支继续工作
	git checkout develop

tr069crt:
	# 1 Generate CA private key
	test -f assets/ca.key || openssl genrsa -out assets/ca.key 4096
	# 2 Generate CA certificate
	test -f assets/ca.crt || openssl req -x509 -new -nodes -key assets/ca.key -days 3650 -out assets/ca.crt -subj \
	"/C=CN/ST=Shanghai/O=teamsacs/CN=TeamsacsCA/emailAddress=master@teamsacs.cc"
	# 3 Generate server private key
	openssl genrsa -out assets/server.key 2048
	# 4 Generate a certificate request file
	openssl req -new -key assets/server.key -out assets/server.csr -subj \
	"/C=CN/ST=Shanghai/O=teamsacs/CN=*.teamsacs.cc/emailAddress=master@teamsacs.cc"
	# 5 Generate a server certificate based on the CA's private key and the above certificate request file
	openssl x509 -req -in assets/server.csr -CA assets/ca.crt -CAkey assets/ca.key -CAcreateserial -out assets/server.crt -days 7300
	mv assets/server.key assets/cwmp.tls.key
	mv assets/server.crt assets/cwmp.tls.crt

clicrt:
	# 1 生成client私钥
	openssl genrsa -out assets/client.key 2048
	# 2 生成client请求文件
	openssl req -new -key assets/client.key -subj "/CN=*.teamsacs.cc" -out assets/client.csr
	# 3 生成client证书
	openssl x509 -req -in assets/client.csr -CA assets/ca.crt -CAkey assets/ca.key -CAcreateserial -out assets/client.crt -days 7300
	mv assets/client.key assets/client.tls.key
	mv assets/client.crt assets/client.tls.crt

updev:
	make buildpre
	make build
	scp ${RELEASE_DIR}/${BUILD_NAME} trdev-server:/tmp/teamsacs
	ssh trdev-server "systemctl stop teamsacs && /tmp/teamsacs -install && systemctl start teamsacs"

.PHONY: clean build tr069crt radseccrt


