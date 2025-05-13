APP        := updoc
PKG_CMD    := ./cmd/updoc
BUILD_DIR  := ./bin

# Export every key=value pair from .env for this command only.
# Using "." works in POSIX sh; "source" caused the error you saw.
load-env = set -o allexport; . .env 2>/dev/null || true; set +o allexport;

.PHONY: run build tidy clean help

run: tidy                          ## Build deps, load .env, then run with optional ARGS="…"
	@echo "==> running $(APP)…";
	@$(load-env) go run $(PKG_CMD) $(ARGS)

build: tidy | $(BUILD_DIR)         ## Compile binary into ./bin/
	go build -o $(BUILD_DIR)/$(APP) $(PKG_CMD)

 tidy:                             ## Ensure go.mod/go.sum are tidy
	go mod tidy

clean:                             ## Remove build artefacts
	rm -rf $(BUILD_DIR)

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

help:                              ## Show available make targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
	  awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'