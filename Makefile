VERSION ?= 1.5.5-release-test
GIT_HASH = $(shell git rev-parse --short HEAD)
DELIVERY_ENGINEERING_GPG_KEY = 4306B2215879C678

.PHONY: vendor test install

vendor:
	dep ensure

test:
	go test -v -json ./... | go-passe

fmt:
	go fmt ./...

install:
	go install -ldflags "-X main.version=${VERSION}-${GIT_HASH}-dev"

release:
	git tag -a "v${VERSION}" -m "Releasing version ${VERSION}"
	git push --tags
	goreleaser --rm-dist

release-deb:
	deb-s3 upload --s3-region=ap-southeast-2 --bucket=apt.redbubble.com --sign=${DELIVERY_ENGINEERING_GPG_KEY} bin/yak_${VERSION}_linux_amd64.deb
