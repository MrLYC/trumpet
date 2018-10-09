VERSION = 0.0.1

ROOTDIR = $(shell pwd)
APPNAME = trumpet
APPPATH = github.com/mrlyc/trumpet
GODIR = /tmp/gopath
SRCDIR = ${GODIR}/src/${APPPATH}
TARGET = bin/${APPNAME}

GOENV = GOPATH=${GODIR} GO15VENDOREXPERIMENT=1

GO = ${GOENV} go
DEP = ${GOENV} dep

LDFLAGS = -X ${APPPATH}/config.Version=${VERSION} -X ${APPPATH}/config.AppName=${APPNAME}
DEBUGLDFLAGS = ${LDFLAGS} -X ${APPPATH}/config.Mode=debug
RELEASELDFLAGS = -w ${LDFLAGS} -X ${APPPATH}/config.Mode=release

GOFINDSH = find . -type f -name "*.go" -not -path "./vendor/*"
GOFILES = $(shell ${GOFINDSH})

.PHONY: release
release: ${SRCDIR}
	${GO} build -i -ldflags="${RELEASELDFLAGS} -X ${APPPATH}/config.BuildHash=`git rev-parse HEAD`" -o ${TARGET} ${APPPATH}

.PHONY: build
build: ${SRCDIR}
	${GO} build -i -ldflags="${DEBUGLDFLAGS}" -o ${TARGET} ${APPPATH}

${SRCDIR}:
	mkdir -p bin
	mkdir -p `dirname "${SRCDIR}"`
	ln -s ${ROOTDIR} ${SRCDIR}

.PHONY: init
init: ${SRCDIR} update

.PHONY: update
update: ${SRCDIR}
	cd ${SRCDIR} && ${DEP} ensure -v

.PHONY: test
test: ${SRCDIR}
	$(eval package ?= $(patsubst ./%,${APPPATH}/%,$(shell find "." -name "*_test.go" -not -path "./vendor/*" -not -path "./.*" -exec dirname {} \; | uniq)))
	${GOENV} go test ${package}

.PHONY: lint
lint:
	${GOENV} ${GOFINDSH} -exec golint -set_exit_status {} \;

.PHONY: fmt
fmt:
	${GOENV} gofmt -w ${GOFILES}

.PHONY: go-env
go-env:
	@${GOENV} go env
