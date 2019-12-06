test:
	go test -race -v github.com/stobias123/gosolar/...
docker:
	docker build -t stobias123/solarcmd:latest . && docker push stobias123/solarcmd:latest
local:
	rm comamnd/solarcmd || true
	cd command
	go build -o solarcmd .