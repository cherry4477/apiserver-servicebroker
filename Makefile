define GITIGNORE
apiserver.local.config
bin
BUILD.bazel
config
default.etcd
docs/build
docs/static_includes
glide.lock
glide.yaml
kubeconfig
pkg/client
pkg/**/*generated*
WORKSPACE

endef
export GITIGNORE

##############################

.PHONY: default
default: build

.PHONY: clean
clean:
	rm -rf bin pkg/client
	find pkg -type f -name "*generated*" -delete

.PHONY: test
test:
	go test ./pkg/... ./cmd/...

.PHONY: docs
docs:
	apiserver-boot build docs

# copy gen to non-gen
.PHONY: g2ng
g2ng:
	for file in \
		pkg/apis/prd/util/util.go \
		pkg/apis/prd/v1/backingserviceinstance_types.go \
		pkg/apis/prd/v1/backingservice_types.go \
		pkg/apis/prd/v1/binding_backingserviceinstance_types.go \
		pkg/apis/prd/v1/servicebroker_types.go \
		pkg/controller/backingservice/controller.go \
		pkg/controller/backingserviceinstance/controller.go \
		pkg/controller/servicebroker/controller.go \
		pkg/controller/sharedinformers/informers.go \
		; \
		do yes | \cp -rf $$file 0-non-gen/$$file; \
	done

# copy non-gen to gen
.PHONY: ng2g
ng2g:
	yes | \cp -rf 0-non-gen/pkg/* pkg
	yes | \cp -rf 0-non-gen/vendor/* vendor

.PHONY: gen
gen: ng2g
	apiserver-boot build generated

.PHONY: build
build: gen
	apiserver-boot build executables --generate=false

.PHONY: build2
build2: ng2g
	apiserver-boot build executables --generate=false

.PHONY: image
image:
	apiserver-boot build container \
		--image apiserver-servicebroker \
		--generate false

.PHONY: local-config
local-config:
	apiserver-boot build config \
		--name apiserver-servicebroker \
		--namespace default \
		--local-minikube

.PHONY: cluster-config
cluster-config:
	apiserver-boot build config \
		--name apiserver-servicebroker \
		--namespace default \
		--image apiserver-servicebroker:lastest

.PHONY: run-local
run-local:
	apiserver-boot run local --generate=false
	# --apiserver=
	# --controller-manager=

.PHONY: run-cluster
run-cluster:
	apiserver-boot run in-cluster --generate=false \
		--name nameofservicetorun --namespace default \
		--image apiserver-servicebroker:lastest

.PHONY: init
init: # to create the skeleton of this project.
	if [ ! -e .gitignore ]; then yes | echo "$$GITIGNORE" > .gitignore; fi
	yes | \rm -rf bin BUILD.bazel cmd docs glide.lock glide.yaml pkg sample vendor WORKSPACE boilerplate.go.txt
	if [ -e 0-non-gen/source-license-head ]; then yes | \cp 0-non-gen/source-license-head boilerplate.go.txt; else touch boilerplate.go.txt; fi
	apiserver-boot init repo --domain asiainfo.com
	apiserver-boot create group version resource --non-namespaced=true \
		--group prd --version v1 --kind ServiceBroker
	apiserver-boot create group version resource --non-namespaced=true \
		--group prd --version v1 --kind BackingService
	apiserver-boot create group version resource --non-namespaced=false \
		--group prd --version v1 --kind BackingServiceInstance
	apiserver-boot create subresource --subresource binding \
		--group prd --version v1 --kind BackingServiceInstance
	apiserver-boot build generated
	
