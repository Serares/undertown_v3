build:
	sam build -t template.yaml

build-local:
	sam build -t template_local.yaml

build-Homepage:
	cd services/ssr/homepage && CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" lambda/main.go

run-Homepage:
	cd services/ssr/homepage && go run *.go