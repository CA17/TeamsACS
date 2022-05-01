BUILD_ORG   := metaslink
BUILD_VERSION   := latest
BUILD_TIME      := $(shell date "+%F %T")
BUILD_NAME      := metaslink
RELEASE_VERSION := v2.0.1
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
	docker build --build-arg BTIME="$(shell date "+%F %T")" -t metaslink . -f Dockerfile
	docker tag metaslink metaslink/metaslink:latest
	docker push metaslink/metaslink:latest

build-ctl:
	make buildpre
	test -d ./release || mkdir -p ./release
	CGO_ENABLED=0 GOARCH=amd64 go build -a -ldflags  '-s -w -extldflags "-static"'  -o ./release/teamsctl teamsctl/teamsctl.go
	upx ./release/teamsctl

ci:
	@read -p "type commit message: " cimsg; \
	git ci -am "$(shell date "+%F %T") $${cimsg}"

syncdev:
	@read -p "提示:同步操作尽量在完成一个完整功能特性后进行，请输入提交描述 (develop):  " cimsg; \
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


gitlog:
	git log --oneline


.PHONY: clean build


