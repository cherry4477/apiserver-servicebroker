.PHONY: default
default: build

.PHONY: clean
clean:
	rm -rf bin

.PHONY: test
test:
	go test ./pkg/... ./cmd/...

# copy non-gen to gen
.PHONY: ng2g
ng2g:
	yes | \cp -rf 0-non-gen/pkg/* pkg
	yes | \cp -rf 0-non-gen/vendor/* vendor

.PHONY: build
build: ng2g
	go build -o bin/apiserver ./cmd/apiserver
	go build -o bin/controller-manager ./cmd/controller-manager

