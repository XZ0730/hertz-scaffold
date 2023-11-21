MODULE = github.com/XZ0730/hertz-scaffold

.PHONY: target
target:
	sh build.sh

.PHONY: clean
clean:
	@find . -type d -name "output" -exec rm -rf {} + -print