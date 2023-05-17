BINARY_NAME=GoldWatcher.exe
APP_NAME=GoldWatcher
VERSION=1.0.0
BUILD_NO=1
APP_ID=com.github.agustincou.fyne.goldwatcher

## build: build binary and package app
build:
	del ${BINARY_NAME}
	fyne package -appID ${APP_ID} -appVersion ${VERSION} -appBuild ${BUILD_NO} -name ${APP_NAME} -release

## run: builds and runs the application
run:
	env DB_PATH="./sql.db" go run .

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@go clean
	@del ${BINARY_NAME}
	@echo "Cleaned!"

## test: runs all tests
test:
	go test -v ./...