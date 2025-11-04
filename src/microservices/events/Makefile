.PHONY: generate build

# Путь к OpenAPI спецификации (можно задать через переменную окружения OPENAPI_SPEC)
OPENAPI_SPEC ?= ../../../api-specification.yaml

# Генерация кода из OpenAPI спецификации
generate:
	@echo "Generating code from OpenAPI specification: $(OPENAPI_SPEC)"
	@which oapi-codegen > /dev/null || go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	@mkdir -p generated
	@oapi-codegen -generate types,server -package generated -o generated/api.gen.go $(OPENAPI_SPEC)
	@echo "Code generated successfully!"

# Сборка приложения
build:
	@echo "Building events service..."
	@go build -o bin/events .

# Установка зависимостей
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Запуск тестов
test:
	@echo "Running tests..."
	@go test -v ./...

# Очистка сгенерированных файлов
clean:
	@echo "Cleaning generated files..."
	@rm -rf generated/*.go
	@rm -rf bin/

