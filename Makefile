start-homepage:
	cd services/ssr/homepage/ && templ generate --watch --cmd="go run main.go"
run-api:
	./scripts/local_start_api.sh