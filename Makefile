.PHONY: test lint run lines smoke
test:
	go test ./...
lint:
	gofmt -w $$(find . -name '*.go')
	go vet ./...
run:
	go run ./cmd/control-api
lines:
	./tools/count-lines.sh
smoke:
	./scripts/smoke.sh
