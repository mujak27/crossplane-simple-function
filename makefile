init:
	crossplane xpkg init learn function-template-go

run:
	go run -v . --insecure --debug

test:
	go test -v -cover .

render:
	crossplane composition render example/xr.yaml example/composition.yaml example/functions.yaml

generate:
	go generate ./...

generate_xrd_from_crd:
	python3 crd_to_xrd.py \
		generated/crds/example.org_xbuckets.yaml \
		generated/xrd.yaml

generate_crd_from_types_using_kubebuilder:
	go get k8s.io/apimachinery@v0.32.3
	go mod tidy

	go run sigs.k8s.io/controller-tools/cmd/controller-gen@v0.17.2 \
		object crd paths=./api/... \
		output:crd:artifacts:config=generated/crds
