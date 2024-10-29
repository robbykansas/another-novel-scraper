BINARY_NAME=another-novel-scraper

run:
	go build
	go install
	${BINARY_NAME}

builds:
	goreleaser release --snapshot --clean