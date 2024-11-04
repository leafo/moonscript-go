.PHONY: test-%

# tests/assign.moon - test-assign
test: parser.go $(patsubst tests/%.moon,test-%,$(wildcard tests/*.moon))
	@echo "pass"

test-write: parser.go $(patsubst tests/%.moon,tests/%.json,$(wildcard tests/*.moon))
	@echo "pass"

moonscript-go: *.go
	go build -o moonscript-go .

parser.go: moonscript.peg
	@pigeon -o parser.go moonscript.peg

test-%: parser.go
	@go run . -json tests/$*.moon

tests/%.json: moonscript-go tests/%.moon
	./moonscript-go -json tests/$*.moon | jq -S > tests/$*.json

