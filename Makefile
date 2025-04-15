BINARY_NAME=another-novel-scraper

run:
	go build
	go install
	${BINARY_NAME} || (echo "Exit command"; exit 1) || true

builds:
	goreleaser release --snapshot --clean