HARBOR_REGISTRY=harbor.amihan.net
APP=product-service
PROJECT=bdo/commons
VERSION=latest
BUILD_IMAGE=${HARBOR_REGISTRY}/${PROJECT}/${APP}:${VERSION}
CHART_REPO=https://${HARBOR_REGISTRY}/chartrepo/bdo
CHARTDIR=bdo
CHART_NAME=product-service

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

pull-image:
	@echo '========================================'
	@echo ' Getting latest image'
	@echo '========================================'
	@docker pull ${HARBOR_REGISTRY}/${PROJECT}/${APP}:${VERSION} || true

package-image:
	@echo '========================================'
	@echo ' Packaging docker image'
	@echo '========================================'
	docker build --network host -t ${BUILD_IMAGE} .
	@echo 'Done.'

package-chart:
	@echo '========================================'
	@echo ' Packaging chart'
	@echo '========================================'
	@mkdir -p build/chart/${APP}/files
	@helm dep update helm/audit-service
	@cp -r helm/${APP} build/chart
	@cp config/rbac.yaml build/chart/${APP}/files/
	@helm package  --app-version ${VERSION} -u -d build/chart build/chart/${APP}
	@echo 'Done.'

package: package-image
	
publish-image: package-image
	@echo '========================================'
	@echo ' Publishing image'
	@echo '========================================'
	docker push ${BUILD_IMAGE}
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
	@echo '${HARBOR_REGISTRY}'
	@echo '========================================'
	@echo ${HARBOR_PASS} | docker login ${HARBOR_REGISTRY} --username ${HARBOR_USER} --password-stdin
	@echo 'Done.'

devops-setup-helm: 
	@echo '========================================'
	@echo ' Setting up Helm'
	@echo '========================================'
	helm init --client-only
	helm repo add msme ${CHART_REPO} --username ${HARBOR_USER} --password ${HARBOR_PASS} 
	helm repo add codecentric https://codecentric.github.io/helm-charts
	helm repo update
	@echo 'Done.'

devops-setup-cluster: 
	kubectl config set-cluster ${CLUSTER_NAME} --server="${CLUSTER_API_URL}"
	kubectl config set clusters.${CLUSTER_NAME}.certificate-authority-data ${CLUSTER_CA}
	kubectl config set-credentials ${CLUSTER_USER} --token="${CLUSTER_USER_TOKEN}"
	kubectl config set-context ${CLUSTER_CONTEXT} --cluster=${CLUSTER_NAME} --user=${CLUSTER_USER}
	kubectl config use-context ${CLUSTER_CONTEXT}
	kubectl config set-context ${CLUSTER_CONTEXT} --namespace=${CLUSTER_NAMESPACE}
	kubectl config view

devops-deploy-chart:
	@echo '========================================'
	@echo ' Deploying application'
	@echo '========================================'
	helm upgrade ${CHART_NAME} --install \
		--namespace ${NAMESPACE} \
		--set image.tag=${VERSION} \
		--set registries[0].url=${HARBOR_REGISTRY} \
		--set registries[0].username=${HARBOR_USER} \
		--set registries[0].password=${HARBOR_PASS} \
		--set extraLabels.git_hash=${CI_COMMIT_SHORT_SHA} \
		--values ${VALUES_FILE} msme/product-service
	@echo 'Done.'

test-quality:
	@echo '========================================'
	@echo ' Running tests: Lint (golangci-lint)'
	@echo '========================================'
	@docker run -it  -v $(PWD):/app --rm golangci/golangci-lint:v1.23.6 bash -c "cd /app && /usr/bin/golangci-lint --verbose -c config/.golangci.yml run"
	@echo NO ERRORS FOUND.
