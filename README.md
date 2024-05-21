# IoT Project - autonomous heating system

## Architecture

![architecture image](./imgs/architectre.svg)

## Install
```bash
pip install pre-commit
pre-commit install
pre-commit run --all-files
```
```bash
go install
go mod tidy
```

## Development
generate OpenAPI schema
```bash
go generate ./...
```

run all services
```bash
docker-compose up
```
