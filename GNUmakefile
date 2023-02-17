default: testacc

# Run acceptance tests
.PHONY: testacc
testacc: EDGE_API_KEY_ID ?= key-id
testacc: EDGE_API_KEY ?= key
testacc: EDGE_API_ENDPOINT ?= http://localhost:8018
testacc:
	TF_ACC=1 \
	EDGE_API_KEY_ID=$(EDGE_API_KEY_ID) \
	EDGE_API_KEY=$(EDGE_API_KEY) \
	EDGE_API_ENDPOINT=$(EDGE_API_ENDPOINT) \
	go test -race -parallel 3 ./... -v $(TESTARGS) -timeout 120m

.PHONY: build-local
build-local: PROVIDER_VERSION ?= 0.0.1
build-local:
	go build -o terraform-provider-edge_v$(PROVIDER_VERSION)
	mkdir -p ~/.terraform.d/plugins/local/edu/cairvine/$(PROVIDER_VERSION)/darwin_amd64
	mv terraform-provider-edge_v$(PROVIDER_VERSION)  ~/.terraform.d/plugins/local/edu/cairvine/$(PROVIDER_VERSION)/darwin_amd64
