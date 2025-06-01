BINARY_NAME=another-novel-scraper

build:
	go build
	go install
	${BINARY_NAME} || (echo "Exit command"; exit 1) || true

run:
	${BINARY_NAME} || (echo "Exit command"; exit 1) || true

release:
	goreleaser release --snapshot --clean

docker-build:
	docker buildx build --platform=linux/arm64 -t another-novel-scraper .

docker-run:
	docker run --rm -it -v $(loc):/Downloads another-novel-scraper