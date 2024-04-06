# Auto Catalog REST API Server 

## Usage

Create `configs/config.toml` file. You can use `configs/config.example.toml` as example

Install dependencies:

``` bash:
make deps
```

Run server:

``` bash:
make run

```

Build binary:

``` bash:
make build
```

## Development

The server invokes GRPC procedures to fetch data from microservices, it includes interfaces generated in following microservices via create GRPC clients and make connection with microservices:
- https://github.com/maximprokopchuk/address_service
- https://github.com/maximprokopchuk/storehouse_service
- https://github.com/maximprokopchuk/auto_reference_catalog_service

## Tests

Run tests:

``` bash:
make test
```

## Linter

Run linter:

``` bash:
make lint
```