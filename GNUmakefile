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
