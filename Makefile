#----------------------------------------------------------------------------
# Makefile
# use the following command to install on a Mac:
# brew install pre-commit
#----------------------------------------------------------------------------

.PHONY: help
.DEFAULT_GOAL := help

VERSION := 1.1.2
SHA256SUM_MAC=$(shell sha256sum boilr-${VERSION}-darwin_amd6.tgz | cut -d' ' -f1 2>&1)

define FORMULA
class Boilr < Formula
  desc "Boilerplate template manager that generates files or directories from template repositories"
  homepage "https://github.com/solaegis/boilr"
  url "https://github.com/solaegis/boilr/releases/download/v${VERSION}/boilr-${VERSION}-darwin_amd64.tgz"
  version "${VERSION}"
  sha256 "${SHA256SUM_MAC}"

  def install
    bin.install "boilr"
  end
end
endef
export FORMULA

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| awk 'BEGIN {FS = ":.*?## "}; \
	{printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

precommit-install: ## install pre-commit
	@pre-commit install

precommit-run: ## run all pre-commit hooks
	@pre-commit run -a

build-mac: ## build the boilr executable
	@go mod tidy
	@go build
	@tar czf boilr-${VERSION}-darwin_amd6.tgz ./boilr
	@mv boilr boilr-mac
	
build-linux: ## build the boilr executable
	@go mod tidy
	@GOOS=linux go build -o boilr-linux

brew-formula-mac: ## create a brew formula for mac
	@echo "$${FORMULA}" > boilr.rb

clean: ## run clean up
	@pre-commit clean
	@rm boilr boilr-mac boilr.rb 2>&1
