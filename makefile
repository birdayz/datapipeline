.PHONY: clean build test bench lint
build:
	@go build -o csvconverter
clean:
	@rm -f ./csvconverter
test:
	@go test $$(go list github.com/birdayz/datapipeline/... | grep -v vendor) -cover
bench:
	@go test $$(go list github.com/birdayz/datapipeline/... | grep -v vendor) -bench=.
lint:
	@go get github.com/golang/lint/golint
	@go get github.com/gordonklaus/ineffassign
	@go vet $$(go list datapipeline/... | grep -v vendor)
	@go list datapipeline/... | grep -v /vendor/ | xargs -L1 golint
	@ineffassign ${CURDIR}

