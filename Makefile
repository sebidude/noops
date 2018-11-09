APPNAME := noops

GITCOMMITHASH := $(shell git log --max-count=1 --pretty="format:%h" HEAD)
GITCOMMIT := -X main.revision=$(GITCOMMITHASH)

VERSIONTAG := 1.0.0
VERSION := -X main.version=$(VERSIONTAG)

BUILDTIMEVALUE := $(shell date +%Y-%m-%dT%H:%M:%S%z)
BUILDTIME := -X main.builddate=$(BUILDTIMEVALUE)

LDFLAGS := '-extldflags "-static" -d -s -w $(GITCOMMIT) $(VERSION) $(BUILDTIME)'
LDFLAGS_WINDOWS := '-extldflags "-static" -s -w $(GITCOMMIT) $(VERSION) $(BUILDTIME)'
all: info clean build

clean:
	rm -f $(APPNAME)-$(VERSIONTAG)-$(GITCOMMITHASH) $(APPNAME)

info: 
	@echo - appname:   $(APPNAME)
	@echo - verison:   $(VERSIONTAG)
	@echo - commit:    $(GITCOMMITHASH)
	@echo - buildtime: $(BUILDTIMEVALUE) 


build:
	@echo Building for linux
	@CGO_ENABLED=0 \
	GOOS=linux \
	go build -o $(APPNAME)-$(VERSIONTAG)-$(GITCOMMITHASH) -a -ldflags $(LDFLAGS) .
	cp $(APPNAME)-$(VERSIONTAG)-$(GITCOMMITHASH) $(APPNAME)
