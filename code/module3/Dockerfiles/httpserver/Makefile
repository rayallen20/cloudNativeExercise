export tag=v1.0

build:
	echo "building httpserver binary"
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 .

release: build
	echo "building httpserver container"
	docker build -t 40486453/httpserver:${tag} .

push: release
	echo "pushing cncamp/httpserver"
	docker push 40486453/httpserver:v1.0