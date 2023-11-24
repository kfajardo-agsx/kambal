APP=account-service
REGISTRY=harbor.amihan.net
CONTAINER_REPO = $(CI_PROJECT_NAMESPACE)/$(CI_PROJECT_NAME)
PROJECT=common
IMAGE_TAG = $(CI_COMMIT_REF_NAME)
TAG = $(CONTAINER_REPO)/$(PROJECT):$(IMAGE_TAG)
PUSH_TAG = $(REGISTRY)/$(CONTAINER_REPO)/$(PROJECT):$(IMAGE_TAG)
REGISTRY_USER = $(HARBOR_AMIHAN_ROBOT_USER)
REGISTRY_PASS = $(HARBOR_AMIHAN_ROBOT_PASS)

all: clean lint build test package

clean:
	@echo '========================================'
	@echo ' Cleaning project'
	@echo '========================================'
	@go clean
	@rm -rf build | true
	@docker-compose down
	@echo 'Done.'

deps:
	@echo '========================================'
	@echo ' Getting Dependencies'
	@echo '========================================'
	@echo 'Cleaning up dependency list...'
	@go mod tidy
	@echo 'Downloading dependencies...'
	go mod download
	@echo 'Vendorizing dependencies...'
	go mod download
	@echo 'Done.'

gen:
	@echo '========================================'
	@echo ' Generating dependencies'
	@echo '========================================'
	@go generate ./cmd

build: deps gen
	@echo '========================================'
	@echo ' Building project'
	@echo '========================================'
	go fmt ./...
	@go build -mod=vendor -o build/bin/${APP} -ldflags "-X main.version=${TAG} -w -s" .
	@echo 'Done.'

test:
	@echo '========================================'
	@echo ' Running tests'
	@echo '========================================'
	@go test ./...
	@echo 'Done.'

lint:
	@echo '========================================'
	@echo ' Running lint'
	@echo '========================================'
	@golint ./...
	@echo 'Done.'

run-deps:
	@echo '========================================'
	@echo ' Running dependencies'
	@echo '========================================'
	@docker-compose up -d

migrate:
	@echo '========================================'
	@echo ' Running migrations'
	@echo '========================================'
	@build/bin/${APP} migrate ${ARGS}


run: build
	@echo '========================================'
	@echo ' Running application'
	@echo '========================================'

	@build/bin/${APP} serve ${ARGS}
	@echo 'Done.'

run-test: 
	@echo '========================================'
	@echo ' Running tests'
	@echo '========================================'
	@go test ./... --coverprofile=coverage.out ./...
	@echo 'Done.'

run-test-all: 
	@echo '========================================'
	@echo ' Running tests'
	@echo '========================================'
	@go test ./... -v --coverpkg=./... --coverprofile=coverage-all.out ./...
	@echo '========================================'
	@echo ' Test Coverage Summary'
	@echo '========================================'
	@go tool cover -func=coverage-all.out
	@echo 'Done.'

pull-image:
	@echo '========================================'
	@echo ' Getting latest image'
	@echo '========================================'
	@docker pull ${REGISTRY}/${PROJECT}/${APP}:${IMAGE_TAG} || true

package-image:
	@echo '========================================'
	@echo ' Packaging docker image'
	@echo '========================================'
	docker build --network host -t ${PUSH_TAG} .
	@echo 'Done.'

package-chart:
	@echo '========================================'
	@echo ' Packaging chart'
	@echo '========================================'
	@mkdir -p build/chart/${APP}/files
	@helm dep update helm/audit-service
	@cp -r helm/${APP} build/chart
	@cp config/rbac.yaml build/chart/${APP}/files/
	@helm package  --app-version ${IMAGE_TAG} -u -d build/chart build/chart/${APP}
	@echo 'Done.'

package: package-image

publish-image: package-image
	@echo '========================================'
	@echo ' Publishing image'
	@echo '========================================'
	docker push ${PUSH_TAG}
	@echo 'Done.'

publish-chart: package-chart
	@echo '========================================'
	@echo ' Publishing chart'
	@echo '========================================'
	helm push build/chart/*.tgz msme
	@echo 'Done.'

publish: publish-image publish-chart

harbor-login:
	@echo '========================================'
	@echo ' Harbor Login'
	@echo '${REGISTRY}'
	@echo '========================================'
	@echo ${HARBOR_PASS} | docker login ${REGISTRY} --username ${HARBOR_USER} --password-stdin
	@echo 'Done.'

test-quality:
	@echo '========================================'
	@echo ' Running tests: Lint (golangci-lint)'
	@echo '========================================'
	@docker run -it  -v $(PWD):/app --rm golangci/golangci-lint:v1.40 bash -c "cd /app && /usr/bin/golangci-lint --verbose --timeout=24h -c lint/.golangci.yml run"
	@echo NO ERRORS FOUND.

# devops
devops-registry-login:
	@echo ${REGISTRY_PASS} | docker login ${REGISTRY} --username ${REGISTRY_USER} --password-stdin
	@echo ${DOCKER_PASS} | docker login --username ${DOCKER_USER} --password-stdin

devops-package-image:
	@echo '========================================'
	@echo ' Packaging docker image'
	@echo '========================================'
	docker build --network host -t ${PUSH_TAG} .
	@echo 'Done.'

devops-publish-image:
	@echo '========================================'
	@echo ' Publishing image'
	@echo '========================================'
	docker push ${PUSH_TAG}
	@echo 'Done.'

devops-deploy-chart:
	helm upgrade ${PROJECT}-${APP}-${ENV} --install \
		--namespace ${NAMESPACE} \
		--set image.repository=$(REGISTRY)/$(CONTAINER_REPO)/$(PROJECT) \
		--set image.tag=${IMAGE_TAG} \
		--set registries[0].url=${REGISTRY} \
		--set registries[0].username=${REGISTRY_USER} \
		--set registries[0].password=${REGISTRY_PASS} \
		--set extraLabels.git_hash=commit-${CI_COMMIT_SHORT_SHA} \
		--values ${VALUES_FILE} helm/application-service

