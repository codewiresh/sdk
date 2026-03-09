OPENAPI_SRC := ../codewire/openapi.yaml
GOBIN       := $(shell go env GOPATH)/bin
DATAMODEL   := $(shell command -v datamodel-codegen 2>/dev/null || echo datamodel-codegen)

.PHONY: generate generate-go generate-ts generate-py check

generate: openapi-sdk.json generate-go generate-ts generate-py

openapi-sdk.json: $(OPENAPI_SRC) scripts/filter-openapi.py
	python3 scripts/filter-openapi.py $(OPENAPI_SRC) openapi-sdk.json

generate-go: openapi-sdk.json
	cd go && $(GOBIN)/oapi-codegen -generate types -package api -o internal/api/types_gen.go ../openapi-sdk.json

generate-ts: openapi-sdk.json
	cd typescript && npx openapi-typescript ../openapi-sdk.json -o src/types.gen.ts

generate-py: openapi-sdk.json
	$(DATAMODEL) --input openapi-sdk.json --output python/codewire/models.py \
	  --output-model-type pydantic_v2.BaseModel --target-python-version 3.10

check:
	cd go && go build ./...
	cd typescript && npx tsc --noEmit
	cd python && python3 -m py_compile codewire/models.py \
	  && python3 -m py_compile codewire/client.py \
	  && python3 -m py_compile codewire/resources.py \
	  && python3 -m py_compile codewire/environment.py \
	  && python3 -m py_compile codewire/secrets.py
