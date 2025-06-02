# Another Novel Scraper

## Standalone bundle
**windows** [another-novel-scraper](https://github.com/robbykansas/another-novel-scraper/releases/download/v1.0.1/another-novel-scraper_Windows_x86_64.zip)
> extract zip files and use run.bat

## Local or Docker
Local or Docker is convenient way to run it
- Clone the project
```
git clone https://github.com/robbykansas/another-novel-scraper
```
- Open directory
```
cd another-novel-scraper
```
- Build locally or Build with docker using Makefile
```
make build
```
or 
```
make docker-build
make docker-run loc=/Users/yourname/Downloads
```
- Build apps without makefile
```
go build
go install
another-novel-scraper
```
- Build docker without makefile
```
docker buildx build --platform=linux/arm64 -t another-novel-scraper .
```
- Docker run
```
docker run --rm -it -v /Users/yourname/Downloads:/Downloads another-novel-scraper
```
> must mounting your download local files when running docker

## Description and example
[cobra]: https://github.com/spf13/cobra
[viper]: https://github.com/spf13/viper
[bubbletea]: https://github.com/charmbracelet/bubbletea

Scrape novel to epub using golang cli by [cobra], [viper] and [bubbletea] to make it faster and interactive, because, why not?

<p>
<img width="100%" src="./assets/ans-example.gif">
</p>
