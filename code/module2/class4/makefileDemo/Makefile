export tag=v1.0

build:
	echo "building makefileDemo binary ${tag}"
	go build -o makefileExec makefileDemo.go

release: build
	echo "building makefileDemo container"
	echo "no docker, fake release"

push: release
	echo "pushing makefileDemo"
	# docker push的命令 但是我没有docker环境 所以省省吧