# Install
```bash
pip install pre-commit
pre-commit install
pre-commit run --all-files
```
```bash
go install
go mod tidy
```

# dev
generate openapi schema
```bash
go generate ./...
```

# TODO
- github CI
- golangcilint
- .pre-commit-config
- OpenAPI schema
- Dockerfile/docker-compose
