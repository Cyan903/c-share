dev:
	${MAKE} -j2 dev-api dev-web

build:
	${MAKE} build-api
	${MAKE} build-web

# dev
dev-api:
	go run cmd/c-share/*.go

dev-web:
	cd web && npm run dev

# build
build-api:
	go build -o ./bin/api ./cmd/c-share/main.go

build-web:
	cd web && npm run build
