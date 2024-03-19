SERVICE_NAME="rdf"

build:
	go build -ldflags="-s -w" -o ${SERVICE_NAME}
run:
	./${SERVICE_NAME} 16 testfileV2 testfileV1
test:
	go test ./delta -v
bench:
	go test -benchtime=1x -bench=. ./delta
