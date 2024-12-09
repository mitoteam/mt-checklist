# Set up variables
EXECUTABLE_NAME := mt-checklist
SUBMODULES := ./internal/dhtml ./internal/mtweb

#do the job
include internal/goappbase/Makefile.inc.mk

# additional targets
.PHONY: run
run:
	clear
	go build
	${EXECUTABLE_NAME} run
