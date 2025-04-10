
KIND            := kindest/node:v1.29.2
ALPINE          := alpine:3.19
POSTGRES        := postgres:16.2

KIND_CLUSTER    := guster
NAMESPACE       := gus-system

VERSION         := 0.0.1

GUS_APP         := gus
AUTH_APP        := auth

BASE_IMAGE_NAME := localhost/gradientsearch

GUS_IMAGE       := $(BASE_IMAGE_NAME)/$(GUS_APP):$(VERSION)
AUTH_IMAGE      := $(BASE_IMAGE_NAME)/$(AUTH_APP):$(VERSION)



# ==================================================================================================
# Debug support

run:
	go run api/cmd/services/gus/main.go |  go run api/cmd/tooling/logfmt/main.go
	 
drun:
	  docker rm -f gus-test && docker run -p 3000:3000 --name gus-test $(GUS_IMAGE) |  go run api/cmd/tooling/logfmt/main.go 

curl-live:
	curl -il -X GET http://localhost:3000/liveness

curl-ready:
	curl -il -X GET http://localhost:3000/readiness

token:
	curl -il \
	--user "admin@example.com:gophers" http://localhost:6000/auth/token/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1


# ==============================================================================
# Building containers

build: gus auth

gus:
	docker build \
		-f zarf/docker/dockerfile.gus \
		-t $(GUS_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.
auth:
	docker build \
		-f zarf/docker/dockerfile.auth \
		-t $(AUTH_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

# ==============================================================================
# Running from within k8s/kind

dev-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner


dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

dev-status:
	watch -n 2 kubectl get pods -o wide --all-namespaces

dev-load-db:
	kind load docker-image $(POSTGRES) --name $(KIND_CLUSTER)

dev-load:
	kind load docker-image $(GUS_IMAGE) --name $(KIND_CLUSTER)
	kind load docker-image $(AUTH_IMAGE) --name $(KIND_CLUSTER)

dev-apply:
	kustomize build zarf/k8s/dev/database | kubectl apply -f -
	kubectl rollout status --namespace=$(NAMESPACE) --watch --timeout=120s sts/database

	kustomize build zarf/k8s/dev/auth | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(AUTH_APP) --timeout=120s --for=condition=Ready

	kustomize build zarf/k8s/dev/gus | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(GUS_APP) --timeout=120s --for=condition=Ready

dev-restart:
	kubectl rollout restart deployment $(GUS_APP) --namespace=$(NAMESPACE)
	kubectl rollout restart deployment $(AUTH_APP) --namespace=$(NAMESPACE)

# --------------------------------------------------------------------------------------------------

dev-describe-deployment:
	kubectl describe deployment --namespace=$(NAMESPACE) $(GUS_APP)

dev-describe-gus:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(GUS_APP)

dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(GUS_APP) --all-containers=true -f --tail=100 --max-log-requests=6 | go run api/cmd/tooling/logfmt/main.go -service=$(GUS_APP)

# ==================================================================================================
# Modules support

tidy:
	go mod tidy


# ==============================================================================
# Metrics and Tracing

metrics:
	expvarmon -ports="localhost:3010" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

statsviz:
	open -a "Google Chrome" http://localhost:3010/debug/statsviz