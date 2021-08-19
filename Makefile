user	:=	$(shell whoami)
rev 	:= 	$(shell git rev-parse --short HEAD)

# GOBIN > GOPATH > INSTALLDIR
GOBIN	:=	$(shell echo ${GOBIN} | cut -d':' -f1)
GOPATH	:=	$(shell echo $(GOPATH) | cut -d':' -f1)
BIN		:= 	""

# check GOBIN
ifneq ($(GOBIN),)
	BIN=$(GOBIN)
else
	# check GOPATH
	ifneq ($(GOPATH),)
		BIN=$(GOPATH)/bin
	endif
endif

all:
	@cd kirito && go build && cd - &> /dev/null
	@cd protoc-gen-go-kirito && go build && cd - &> /dev/null

.PHONY: install
.PHONY: uninstall
.PHONY: clean
.PHONY: fmt

install: all
ifeq ($(user),root)
#root, install for all user
	@cp ./kirito/kirito /usr/bin
	@cp ./protoc-gen-go-kirito/protoc-gen-go-kirito /usr/bin
else
#!root, install for current user
	$(shell if [ -z $(BIN) ]; then read -p "Please select installdir: " REPLY; mkdir -p $${REPLY};\
	cp ./kirito/kirito $${REPLY}/;cp ./protoc-gen-go-kirito/protoc-gen-go-kirito $${REPLY}/;else mkdir -p $(BIN);\
	cp ./kirito/kirito $(BIN);cp ./protoc-gen-go-kirito/protoc-gen-go-kirito $(BIN); fi)
endif
	@which protoc-gen-xgo &> /dev/null || go get github.com/go-kirito/protobuf-go/cmd/protoc-gen-xgo
	@which protoc-gen-validate  &> /dev/null || go get github.com/envoyproxy/protoc-gen-validate
	@echo "install finished"

uninstall:
	$(shell for i in `which -a kirito | grep -v '/usr/bin/kirito' 2>/dev/null | sort | uniq`; do read -p "Press to remove $${i} (y/n): " REPLY; if [ $${REPLY} = "y" ]; then rm -f $${i}; fi; done)
	$(shell for i in `which -a protoc-gen-validate | sort | uniq`; do read -p "Press to remove $${i} (y/n): " REPLY; if [ $${REPLY} = "y" ]; then rm -f $${i}; fi; done)
	@echo "uninstall finished"

clean:
	@go mod tidy
	@echo "clean finished"

fmt:
	@gofmt -s -w .