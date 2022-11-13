.PHONY: default
default: build

FAB := $(shell command -v fab 2> /dev/null)

%:
ifndef FAB
	@go install github.com/sheik/fab/cmd/fab@latest
endif
	@fab $@

