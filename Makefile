PROJECT_DIR = $(CURDIR)
PROJECT_NAME = $(shell basename $(PROJECT_DIR))
PROJECT_BIN = $(PROJECT_DIR)/bin


.DEFAULT_GOAL := build

.PHONY: install-all
install-all: install-glide install-deps 

install-glide:
	@command -v glide >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "--> installing glide"; \
		curl https://glide.sh/get | sh; \
	fi
install-deps:
	@test -d $(PROJECT_DIR)/vendor || mkdir $(PROJECT_DIR)/vendor
	@glide install

update-deps:
	@glide update

build: $(PROJECT_BIN)
	@echo "  >  Building binary..."
#	@GOBIN=$(PROJECT_BIN) go build ./...
	@go build -o ${PROJECT_BIN}/app cmd/app/main.go 
	@go build -o ${PROJECT_BIN}/asap_worker cmd/asap_worker/main.go
	@go build -o ${PROJECT_BIN}/regular_worker cmd/regular_worker/main.go

run:
	@${PROJECT_BIN}/app -source=$(PROJECT_DIR)/ &
	@${PROJECT_BIN}/asap_worker -source=$(PROJECT_DIR)/ &
	@${PROJECT_BIN}/regular_worker &

$(PROJECT_BIN):
	@test -d $@ || mkdir $@ 
