GO ?= go

TARGETS := gateway productsvc
TARGETDIR := bin

.PHONY: build
build:
	@for target in $(TARGETS) ; do \
		$(GO) build -o "$(TARGETDIR)/$$target" ./cmd/$$target.go ; \
	done

.PHONY: clean
clean:
	@$(RM) -r $(TARGETDIR)