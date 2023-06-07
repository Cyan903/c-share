# Makefile configs

APP = "c-share" # App name
WEB = ${PWD}/web # Ensure this points to web
UPX = upx2 # https://upx.github.io/ (optional)
PKG = npm # pnpm, yarn, etc...

format:
	@printf "Formating Go...\n"
	@go fmt ./...
	@go mod tidy
	@printf "Formating Web...\n"
	@cd $(WEB) && $(PKG) run format

update:
	@printf "Updating API...\n"
	@go mod tidy
	@printf "Updating web...\n"
	@cd $(WEB) && $(PKG) install

## Production
.PHONY: build
build:
	@printf "Building API & web...\n\n"
	@${MAKE} build-api
	@${MAKE} build-web

.PHONY: build-api
build-api:
	@printf "Building API...\n"
	@go build -o ./bin/$(APP) ./cmd/c-share/main.go
	@printf "Do you want to minify the API with upx? (y/N) "; \
		read answer_minify; \
		if [ $$answer_minify == "y" ]; then $(UPX) ./bin/$(APP); fi

.PHONY: build-web
build-web:
	@printf "Building web...\n"
	@cd $(WEB) && $(PKG) run build

## Development
.PHONY: dev
dev:
	@printf "Starting API & web without building...\n\n"
	@${MAKE} -j2 dev-api dev-web

.PHONY: dev-api
dev-api:
	@printf "Running API...\n"
	@go run cmd/c-share/*.go

.PHONY: dev-web
dev-web:
	@printf "Running web...\n"
	@cd $(WEB) && $(PKG) run dev
