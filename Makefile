.PHONY: test test-% write-tests

# we list out the lua ones manually because there are some that can't compile to lua directly since they need transformations
ALL_TESTS = $(patsubst tests/%.moon,test-%,$(wildcard tests/*.moon))

LUA_TESTS = tests/empty.lua tests/empty_comments.lua tests/assign.lua tests/functions.lua tests/tables.lua tests/chain.lua tests/statements.lua tests/simple.lua tests/loops.lua

# run all tests
test: parser.go $(patsubst tests/%.moon,test-%,$(wildcard tests/*.moon))
	@echo "pass"

# write the output of all tests to tests dir, so they can be diff'd
write-tests: parser.go $(patsubst tests/%.moon,tests/%.json,$(wildcard tests/*.moon)) $(LUA_TESTS)
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

