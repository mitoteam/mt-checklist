# Set up variables
EXECUTABLE_NAME := mt-checklist
SUBMODULES := ./dhtml ./mtweb

#do the job
include goappbase/Makefile.inc.mk

# additional targets
.PHONY: run
run:
	clear
	go run main.go run
