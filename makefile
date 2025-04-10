

GUS_APP       := gus
BASE_IMAGE_NAME := localhost/gradientsearch
VERSION       := 0.0.1
GUS_IMAGE     := $(BASE_IMAGE_NAME)/$(GUS_APP):$(VERSION)



# ==================================================================================================
# Debug support


run:
	go run apis/services/sales/main.go

drun:
	docker run -p 3000:3000 $(GUS_IMAGE)


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

# ==================================================================================================
# Modules support

tidy:
	go mod tidy