NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

all: test

format:
	go fmt ./...

test:
	@echo "$(OK_COLOR)==> Running...$(NO_COLOR)"
	@ginkgo -r -trace -keepGoing

.PHONY: all format test
