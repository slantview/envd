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
	./test.sh