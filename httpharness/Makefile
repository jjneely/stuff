# Dockerable static binary
httpharness: *.go
	go mod tidy
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o $@ .

docker: httpharness Dockerfile
	docker build -t httpharness:latest -f Dockerfile .
