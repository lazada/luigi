# To see a list of commands run "make help"
##To run command use "make [name of command]", for example "make build"
##

LAST_GOPATH_DIR:=$(lastword $(subst :, ,$(GOPATH)))
GOBIN:=$(LAST_GOPATH_DIR)/bin
PROJECT_PATH:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
GLIDE_PATH:=$(LAST_GOPATH_DIR)/src/github.com/Masterminds/glide
GLIDE_VERSION:=v0.11.0
GLIDE_BIN:=$(GOBIN)/glide-$(GLIDE_VERSION)
GOHINT_CONFIG_PATH:=$(PROJECT_PATH)/deploy/go_hint_config.json
LINT_EXCLUDE:= --exclude='vendor/*'
LINT_FAST:=gometalinter $(LINT_EXCLUDE) --deadline=180s --cyclo-over=50 --min-const-length=4 --min-occurrences=10 --line-length=300 --disable-all --enable=vet --enable=vetshadow --enable=gosimple --enable=staticcheck --enable=ineffassign --enable=gocyclo --enable=lll --enable=goconst ./...
LINT_SLOW:=gometalinter $(LINT_EXCLUDE) --deadline=180s --dupl-threshold=290 --disable-all --enable=unconvert --enable=unused --enable=varcheck --enable=dupl ./...

PACKAGES_WITH_TESTS := $(shell find . -name '*_test.go' -not -path "./vendor/*" -exec dirname '{}' ';' | sort -u | sed -e 's/^\.\///')
TEST_TARGETS := $(foreach p,$(PACKAGES_WITH_TESTS),test-$(p))
TEST_USE_TAGS?=0
TEST_TAGS?=all
TEST_TMP_DIR=./tmp/ci
TEST_REPORT_DIR=$(TEST_TMP_DIR)/reports
TEST_COVER_DIR=$(TEST_TMP_DIR)/cover
TEST_COVERAGE_OUTPUT?=cover.out
TEST_COVERAGE_FILE=$(TEST_TMP_DIR)/$(TEST_COVERAGE_OUTPUT)
TEST_COVERAGE_FAIL_FILE=$(TEST_TMP_DIR)/test_coverage_fail_checker
CHECKSTYLE_GOHINT_RESULT?=$(TEST_TMP_DIR)/checkstyle-result-gohint.xml
CHECKSTYLE_GOMETALINTER_RESULT?=$(TEST_TMP_DIR)/checkstyle-result-gometalinter.xml
PARALLEL_MAKE?=$(shell getconf _NPROCESSORS_ONLN) #CPU cores


build:                   ##nothing to build in library, just getting deps
build: deps


deps:                    ##install dependencies
deps: get-glide
	$(info #Install dependencies...)
	$(GLIDE_BIN) install --force


ci-lint:                 ##run gohint and gometalinter
ci-lint: get-lint | $(TEST_TMP_DIR)
	$(info #Run gohint and gometalinter...)
	gohint -config=$(GOHINT_CONFIG_PATH) -reporter=checkstyle > $(CHECKSTYLE_GOHINT_RESULT) || true
	$(LINT_FAST) --concurrency=1 --checkstyle > $(CHECKSTYLE_GOMETALINTER_RESULT) || true


ci-test-all:             ##run all tests
ci-test-all: _test-deps _test-install | $(TEST_TMP_DIR)
	$(info #Running tests with coverage. Parallel: $(PARALLEL_MAKE). Output will be in $(TEST_COVERAGE_FILE)...)
	@mkdir -p $(TEST_REPORT_DIR) $(TEST_COVER_DIR)
	@if [ -f $(TEST_COVERAGE_FAIL_FILE) ] ; then \
		  rm $(TEST_COVERAGE_FAIL_FILE) ;\
		fi;
	@$(MAKE) -j $(PARALLEL_MAKE) $(TEST_TARGETS) TEST_TAGS=$(TEST_TAGS) # -l $(LOAD_AVERAGE)
	@echo "mode: set" > $(TEST_COVERAGE_FILE)
	@cat $(TEST_COVER_DIR)/*.out | grep -v "mode: set"| grep -v "mode: atomic" >> $(TEST_COVERAGE_FILE)
	@cat $(TEST_REPORT_DIR)/*.out | go2xunit -fail > $(TEST_TMP_DIR)/report.xml
	@gocover-cobertura < $(TEST_COVERAGE_FILE) > $(TEST_TMP_DIR)/coverage.xml
	@cat $(TEST_REPORT_DIR)/*.out
	@rm -rf $(TEST_REPORT_DIR) $(TEST_COVER_DIR)
	@if [ -f $(TEST_COVERAGE_FAIL_FILE) ] ; then \
		  rm $(TEST_COVERAGE_FAIL_FILE) ;\
		  exit 1 ;\
	fi ;

_test-install:
	go test -i -tags=$(TEST_TAGS) `$(GLIDE_BIN) nv`


$(TEST_TMP_DIR):
	@mkdir -p $(TEST_TMP_DIR)

$(TEST_TARGETS):
	$(eval $@_files=1)
	$(eval $@_package := $(subst test-,,$@))
	$(eval $@_fname := $(subst /,_,$($@_package)))
	$(eval $@_coverprofile := $(subst .,_,$($@_fname)))
	$(eval $@_reportfile := $(subst .,_,$($@_fname)))

ifeq ($(TEST_USE_TAGS),1)
	$(eval $@_TEST_ARGS := -tags='$(TEST_TAGS)')
	$(eval $@_files := $(shell ls -d $($@_package)/* | grep _test.go | xargs grep "+build" | grep '$(TEST_TAGS)' | grep -v '!$(TEST_TAGS)' | head -1 | wc -l | awk '{ print $1}'))
else
	$(eval $@_TEST_ARGS :=)
endif
	$(eval $@_TEST_ARGS := $($@_TEST_ARGS) -v  -coverprofile $(TEST_COVER_DIR)/$($@_coverprofile).out)
	@echo "== Directory $($@_package) == \n" > $(TEST_REPORT_DIR)/$($@_reportfile).out
	@if [ "$($@_files)" -eq "1" ] ; then \
		go test ./$($@_package) $($@_TEST_ARGS) >> $(TEST_REPORT_DIR)/$($@_reportfile).out || ( echo 'fail $($@_package)' >> $(TEST_COVERAGE_FAIL_FILE); exit 0); \
	fi;

_test-deps: deps
	@go get godep.lzd.co/gocover-cobertura
	@go get godep.lzd.co/go2xunit


get-glide:               ##install glide
ifeq ($(wildcard $(GLIDE_BIN)),)
	$(info #Installing glide version $(GLIDE_VERSION)...)
ifeq ($(wildcard $(GLIDE_PATH)),)
	mkdir -p $(GLIDE_PATH) && cd $(GLIDE_PATH) ;\
	git clone https://github.com/Masterminds/glide.git .
endif
	cd $(GLIDE_PATH) && git fetch --tags && git checkout $(GLIDE_VERSION) ;\
	git reset --hard && git clean -fd ;\
	make clean && make build && mv glide $(GLIDE_BIN)
else
	$(info #Found glide $(GLIDE_VERSION) in $(GLIDE_BIN)...)
endif


get-lint:                ##install gometalinter
	go get -u godep.lzd.co/gometalinter
	gometalinter --install
	go get -u godep.lzd.co/hint/gohint


help:                    ##show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'


.PHONY: build deps ci-lint ci-test-all _test-install ci-deps get-glide get-lint help
