
KIND            := kindest/node:v1.29.2
ALPINE          := alpine:3.19
POSTGRES        := postgres:16.2

KIND_CLUSTER    := guster
NAMESPACE       := gus-system
GUS_APP         := gus
BASE_IMAGE_NAME := localhost/gradientsearch
VERSION         := 0.0.1
GUS_IMAGE       := $(BASE_IMAGE_NAME)/$(GUS_APP):$(VERSION)


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


# ==============================================================================
# Building containers

build: gus

gus:
	docker build \
		-f zarf/docker/dockerfile.gus \
		-t $(GUS_IMAGE) \
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

dev-load:
	kind load docker-image $(GUS_IMAGE) --name $(KIND_CLUSTER)

dev-apply:
	kustomize build zarf/k8s/dev/gus | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(GUS_APP) --timeout=120s --for=condition=Ready

dev-restart:
	kubectl rollout restart deployment $(GUS_APP) --namespace=$(NAMESPACE)

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