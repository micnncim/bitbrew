PKGS := $(shell go list ./... | grep -v vendor)

.PHONY: test
test:
	go test -v -parallel=4 $(PKGS)

.PHONY: vet
vet:
	go vet $(PKGS)

.PHONY: lint
lint:
	golint $(PKGS)

.PHONY: coverage
coverage:
	go test -v -race -covermode=atomic -coverpkg=./... -coverprofile=coverage.txt ./...

.PHONY: reviewdog
reviewdog:
	reviewdog -reporter=github-pr-review