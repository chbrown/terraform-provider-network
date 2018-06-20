all: terraform-provider-network

terraform-provider-network: main.go provider.go
	go build -o $@

.PHONY: pre-dist clean

pre-dist:
	gofmt -s -w *.go

clean:
	rm -f terraform-provider-network
