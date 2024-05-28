start:
	go run ./main.go

build:
	go build -o ./dist/gen-manifest ./main.go
	chmod +x ./dist/gen-manifest

install-local:
	go build -o $(GOBIN)/gen-manifest ./main.go
	chmod +x $(GOBIN)/gen-manifest

install:
	go install github.com/asolopovas/gen-manifest@latest

test:
	 go run ./main.go -c ./gen-manifest-config.json

tag-push:
	$(eval VERSION=$(shell cat version))
	git tag $(VERSION)
	git push origin $(VERSION)
	if git rev-parse latest >/dev/null 2>&1; then git tag -d latest; fi
	git tag latest
	git push origin latest --force
