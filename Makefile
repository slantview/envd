bindir := bin
outbin := ${bindir}/envd

build: get_deps ${outbin}
build_all: get_deps
	go build -a -o ${outbin}
${outbin}: envd.go
	go build -o ${outbin}
clean:
	rm -f ${bindir}/*
	rm .deps
get_deps: .deps
.deps: Godeps/Godeps.json
	godep restore
	touch .deps
test:
	go test -v .
fmt:
	@status=$$(git status --porcelain | grep -v '??' | grep '.go$$' 2> /dev/null); \
	if test "x$${status}" = x; then \
		gofmt -w .; \
	else \
		echo "[ERROR] gofmt changes should be committed in their own commit.\n\tDo not run fmt without committing your changes first to avoid muddying your changes with formatting corrections.\n\tEither 'git reset --hard' or 'git commit' to run 'make fmt' or run 'gofmt' manually."; \
		echo '\n\tChanged files:'; \
		git status --porcelain | grep -v '??' | grep '.go$$' | awk '{printf "\t\t%s\n", $$2}' 2> /dev/null; \
	fi