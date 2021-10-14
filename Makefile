# Definitions
ROOT                    := $(PWD)
GO_HTML_COV             := ./coverage.html
GO_TEST_OUTFILE         := ./c.out
GOLANG_DOCKER_IMAGE     := golang:1.17.2
GOLANG_DOCKER_CONTAINER := ferfabricio/vwap-test


#   Deletes container if exists
#   Usage:
#       make clean
clean:
	docker rm -f ${GOLANG_DOCKER_CONTAINER} || true

#   Usage:
#       make test
test:
	docker run -w /app -v ${ROOT}:/app ${GOLANG_DOCKER_IMAGE} go test ./... -coverprofile=${GO_TEST_OUTFILE}
	docker run -w /app -v ${ROOT}:/app ${GOLANG_DOCKER_IMAGE} go tool cover -html=${GO_TEST_OUTFILE} -o ${GO_HTML_COV}

#   Usage:
#       make run
run:
	docker run -w /app -v ${ROOT}:/app ${GOLANG_DOCKER_IMAGE} go run main.go
