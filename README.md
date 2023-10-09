## Usage

### Build
```shell
go mod download
go build -o user-collector
```

### Run
```shell
./user-collector crawl-devto-users \
--from=0 \
--to=1100000 \
--concurrent=20 \
--proxy=http://localhost:8888
```

### Docker compose
```shell
docker compose up -d
```