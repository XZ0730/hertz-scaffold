include ./common.mk

SERVICE_NAME = runFzu

.PHONY: build
build:
	sh build.sh

.PHONY: new
new:
	hz new \
	-module $(MODULE) \
	hz update -idl ./idl/api.thrift

.PHONY: gen
gen:
	hz update -idl ./idl/api.thrift

.PHONY: server
server:
	make build
	sed -i 's/\r//' ./output/bootstrap.sh
	cd output && sh bootstrap.sh