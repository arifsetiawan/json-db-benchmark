
run-bench:
	export $$(grep -v ^\# .env | grep . | cat | xargs) && go test -timeout 30m -bench=.

run-test:
	export $$(grep -v ^\# .env | grep . | cat | xargs) && go test -v -timeout 30s github.com/arifsetiawan/json-db-benchmark -run TestDropInit