.PHONY: test-%

# tests/assign.moon - test-assign
test: parser.go $(patsubst tests/%.moon,test-%,$(wildcard tests/*.moon))
	@echo "pass"

test-write: parser.go $(patsubst tests/%.moon,tests/%.json,$(wildcard tests/*.moon))
	@echo "pass"

parser.go: moonscript.peg
	@pigeon -o parser.go moonscript.peg

test-%: parser.go
	@go run . tests/$*.moon

tests/%.json: parser.go tests/%.moon
	go run . tests/$*.moon | jq -S > tests/$*.json

