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

clean:
	rm -f teamsacs

gen:
	go generate

build:
	go generate
	CGO_ENABLED=0 go build -a -ldflags \
	'\
	-X "main.BuildVersion=${BUILD_VERSION}"\
	-X "main.ReleaseVersion=${RELEASE_VERSION}"\
	-X "main.BuildTime=${BUILD_TIME}"\
	-X "main.BuildName=${BUILD_NAME}"\
	-X "main.CommitID=${COMMIT_SHA1}"\
	-X "main.CommitDate=${COMMIT_DATE}"\
	-X "main.CommitUser=${COMMIT_USER}"\
	-X "main.CommitSubject=${COMMIT_SUBJECT}"\
	-s -w -extldflags "-static"\
	' \
    -o ${BUILD_NAME} ${SOURCE}

build-linux:
	go generate
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags \
	'\
	-X "main.BuildVersion=${BUILD_VERSION}"\
	-X "main.ReleaseVersion=${RELEASE_VERSION}"\
	-X "main.BuildTime=${BUILD_TIME}"\
	-X "main.BuildName=${BUILD_NAME}"\
	-X "main.CommitID=${COMMIT_SHA1}"\
	-X "main.CommitDate=${COMMIT_DATE}"\
	-X "main.CommitUser=${COMMIT_USER}"\
	-X "main.CommitSubject=${COMMIT_SUBJECT}"\
	-s -w -extldflags "-static"\
	' \
    -o ${RELEASE_DIR}/${BUILD_NAME} ${SOURCE}

pubbuild-pre:
	make build-linux
	make upx
	echo 'FROM alpine' > .build
	echo 'ARG CACHEBUST="$(shell date "+%F %T")"' >> .build
	echo 'COPY ./teamsacs /teamsacs' >> .build
	echo 'RUN chmod +x /teamsacs' >> .build
	echo 'EXPOSE 1979 1980 1981 1812/udp 1813/udp 1914/udp 1924/udp 1914/udp' >> .build
	echo 'ENTRYPOINT ["/teamsacs"]' >> .build
	scp ${RELEASE_DIR}/${BUILD_NAME} DockerServer:/tmp/teamsacs
	scp .build DockerServer:/tmp/.teamsacsbuild

pubdev:
	make pubbuild-pre
	ssh DockerServer "cd /tmp \
	&& sudo docker build -t teamsacs . -f .teamsacsbuild \
	&& sudo docker tag teamsacs alab.189csp.cn:5000/teamsacs:dev \
	&& sudo docker push alab.189csp.cn:5000/teamsacs:dev \
	&& rm -f /tmp/teamsacs \
	&& rm -f /tmp/.teamsacsbuild "
	rm -f .build

fastpub:
	make pubbuild-pre
	ssh DockerServer "cd /tmp \
	&& sudo docker build -t teamsacs . -f .teamsacsbuild \
	&& sudo docker tag teamsacs alab.189csp.cn:5000/teamsacs:latest \
	&& sudo docker push alab.189csp.cn:5000/teamsacs:latest \
	&& rm -f /tmp/teamsacs \
	&& rm -f /tmp/.teamsacsbuild "
	rm -f .build

github:
	make pubbuild-pre
	ssh DockerServer "cd /tmp \
	&& sudo docker build -t teamsacs . -f .teamsacsbuild \
	&& sudo docker tag teamsacs docker.pkg.github.com/ca17/teamsacs/teamsacs:latest \
	&& sudo docker push docker.pkg.github.com/ca17/teamsacs/teamsacs:latest \
	&& rm -f /tmp/teamsacs \
	&& rm -f /tmp/.teamsacsbuild "
	rm -f .build

upx:
	upx ${RELEASE_DIR}/${BUILD_NAME}

ci:
	@read -p "type commit message: " cimsg; \
	git ci -am "$(shell date "+%F %T") $${cimsg}"

push:
	@read -p "type commit message: " cimsg; \
	git ci -am "$(shell date "+%F %T") $${cimsg}"
	git push origin main

build-abfs:
	CGO_ENABLED=0 go build -a -ldflags '-s -w -extldflags "-static"' -o abfs commands/abfs/abfs.go

build-labfs:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w -extldflags "-static"' -o labfs commands/abfs/abfs.go

.PHONY: clean build rpccert webcert


