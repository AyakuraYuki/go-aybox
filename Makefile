NAME=go-aybox
BASE_BUILD_DIR=build
BUILD_NAME=$(GOOS)-$(GOARCH)$(GOARM)
BUILD_DIR=$(BASE_BUILD_DIR)/$(BUILD_NAME)
VERSION?=dev

ifeq ($(GOOS),windows)
  ext=.exe
  archiveCmd=zip -9 -r $(NAME)-$(BUILD_NAME)-$(VERSION).zip $(BUILD_NAME)
else
  ext=
  archiveCmd=tar czpvf $(NAME)-$(BUILD_NAME)-$(VERSION).tar.gz $(BUILD_NAME)
endif

test:
	go test -race -v -bench=. ./...
