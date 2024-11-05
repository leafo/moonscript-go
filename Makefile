.PHONY: test test-% write-tests

ALL_TESTS = $(patsubst tests/%.moon,test-%,$(wildcard tests/*.moon)) tests/empty.lua tests/empty_comments.lua tests/assign.lua tests/functions.lua tests/tables.lua tests/functions.lua tests/chain.lua tests/statements.lua tests/simple.lua tests/loops.lua

# run all tests
test: parser.go $(ALL_TESTS)
	@echo "pass"

# write the output of all tests to tests dir, so they can be diff'd
write-tests: parser.go $(ALL_TESTS)
	@echo "pass"

moonscript-go: *.go
	@go build -o moonscript-go .

parser.go: moonscript.peg
	@pigeon -o parser.go moonscript.peg

test-%: moonscript-go
	@./moonscript-go -json tests/$*.moon

tests/%.json: moonscript-go tests/%.moon
	./moonscript-go -json tests/$*.moon | jq -S > tests/$*.json

tests/%.lua: moonscript-go tests/%.moon
	./moonscript-go tests/$*.moon > tests/$*.lua

