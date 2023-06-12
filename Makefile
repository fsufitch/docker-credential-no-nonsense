GOARGS ?= 
SOURCES := go.mod go.sum $(shell find . -name '*.go')
CMDPKG := ./cmd/docker-credential-no-nonsense

# See possibilities with "go tool dist list"
DISTS ?=

default:
	echo ${SOURCES}
	DISTS='${DISTS}'; \
	TARGETS=''; \
	if [ -z "$$DISTS"]; then \
		DISTS=$$(eval "$$(go tool dist env) echo \$$GOOS/\$$GOARCH"); \
	fi; \
	for dist in $$(echo "$$DISTS" | grep -Eo '\S+'); do \
		if [ -z "$$(echo $$dist | grep '^windows/')" ]; then \
			SUFFIX=''; \
		else \
			SUFFIX='.exe'; \
		fi; \
		TARGETS="$$TARGETS dist/$$dist/docker-credential-no-nonsense$$SUFFIX"; \
	done; \
	echo DISTS $$DISTS && echo TARGETS $$TARGETS && make $$TARGETS;
.PHONY: default


all-dists:
	DISTS="$$(go tool dist list)" make default
.PHONY: all-dists

clean:
	rm -rf dist
.PHONY: clean

dist/%/docker-credential-no-nonsense: ${SOURCES}
	@echo '$@' | ( grep -v 'dist/windows' > /dev/null ) || ( echo "ERROR: Windows build must have .exe suffix" && exit 1 )
	@\
		export GOOS="$$(echo '$@' | sed -E 's@^dist/([^/]+)/.*@\1@')"; \
		export GOARCH="$$(echo '$@' | sed -E 's@^dist/[^/]+/([^/]+)/.*@\1@')"; \
		echo "Building... GOOS=$$GOOS GOARCH=$$GOARCH"; \
		go build ${GOARGS} -o '$@' '${CMDPKG}'

dist/%/docker-credential-no-nonsense.exe: ${SOURCES}
	@echo '$@' | ( grep '^dist/windows/' > /dev/null ) || (echo "ERROR: non-Windows may not have .exe suffix" >&2 && exit 1)
	@\
		export GOOS="$$(echo '$@' | sed -E 's@^dist/([^/]+)/.*@\1@')"; \
		export GOARCH="$$(echo '$@' | sed -E 's@^dist/[^/]+/([^/]+)/.*@\1@')"; \
		echo "Building... GOOS=$$GOOS GOARCH=$$GOARCH"; \
		go build ${GOARGS} -o '$@' '${CMDPKG}'
