

test: parser.go
	for file in tests/*.moon; do \
		echo $$file; \
		go run . $$file; \
	done

parser.go: moonscript.peg
	pigeon -o parser.go moonscript.peg


