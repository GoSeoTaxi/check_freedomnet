PROJECT_NAME := freedomnet_checker

BIN_DIR := bin

DOCKER_IMAGE := $(PROJECT_NAME)_image
DOCKER_CONTAINER := $(PROJECT_NAME)_container

PORT := $(shell grep -m 1 'PORT' .env | cut -d '=' -f2)

GOOS_LIST := linux darwin windows
GOARCH_LIST := amd64 arm64

.PHONY: all
all: docker-build docker-run

.PHONY: docker-build
docker-build:
	@echo "Сборка Docker образа..."
	docker build -t $(DOCKER_IMAGE) .

.PHONY: docker-run
docker-run: docker-build
	@echo "Запуск Docker контейнера на порту $(PORT)..."
	docker run -d -p 8081:$(PORT) --env-file .env --name $(DOCKER_CONTAINER) $(DOCKER_IMAGE)

.PHONY: docker-stop
docker-stop:
	@echo "Остановка и удаление Docker контейнера..."
	docker stop $(DOCKER_CONTAINER) || true
	docker rm $(DOCKER_CONTAINER) || true

.PHONY: build-all
build-all:
	@echo "Компиляция бинарников для всех платформ..."
	@for goos in $(GOOS_LIST); do \
		for goarch in $(GOARCH_LIST); do \
			ext=""; \
			if [ "$$goos" = "windows" ]; then ext=".exe"; fi; \
			output_name=$(BIN_DIR)/$(PROJECT_NAME)_$$goos-$$goarch$$ext; \
			echo "Компилируем для GOOS=$$goos и GOARCH=$$goarch..."; \
			GOOS=$$goos GOARCH=$$goarch go build -o $$output_name cmd/server/main.go; \
			echo "Скомпилирован $$output_name"; \
		done \
	done

.PHONY: clean
clean:
	@echo "Очистка папки с бинарниками..."
	rm -rf $(BIN_DIR)/*
	@echo "Очистка завершена."

.PHONY: docker-clean
docker-clean: docker-stop
	@echo "Удаление Docker образа..."
	docker rmi $(DOCKER_IMAGE) || true
	@echo "Docker ресурсы очищены."