PKGS := $(shell go list ./... | grep -v vendor)

.PHONY: dep
dep:
	go mod download
	go mod verify

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

.PHONY: tools
tools:
	go get github.com/reviewdog/reviewdog/cmd/reviewdog
	go get golang.org/x/lint/golint

.PHONY: reviewdog
reviewdog: tools
	reviewdog -reporter=github-pr-review
