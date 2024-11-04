.PHONY: test-%

# tests/assign.moon - test-assign
test: parser.go $(patsubst tests/%.moon,test-%,$(wildcard tests/*.moon))
	@echo "pass"

parser.go: moonscript.peg
	@pigeon -o parser.go moonscript.peg

test-%: parser.go
	@go run . tests/$*.moon


