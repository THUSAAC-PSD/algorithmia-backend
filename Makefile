.PHONY: install-tools
install-tools:
	@./scripts/install-tools.sh

.PHONY: run-app
run-app:
	@./scripts/run.sh  

.PHONY: build
build:
	@./scripts/build.sh 

.PHONY: install-dependencies
install-dependencies:
	@./scripts/install-dependencies.sh 

# .PHONY: integration-test
# integration-test:
# 	@./scripts/test.sh integration
#
# .PHONY: e2e-test
# e2e-test:
# 	@./scripts/test.sh e2e

#.PHONY: load-test
#load-test:
#	@./scripts/test.sh load-test

.PHONY: format
format:
	@./scripts/format.sh 

.PHONY: lint
lint:
	@./scripts/lint.sh 

.PHONY: update-dependencies
update-dependencies:
	@./scripts/update-dependencies.sh
