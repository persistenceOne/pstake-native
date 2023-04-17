CHAIN_ID := localnet

CHAIN_DIR ?= /tmp/trash
CHAIN_BIN ?= ./../../build/pstaked
CHAIN_DATA_DIR ?= .pstaked

all: clean setup start

.bash:
	CHAIN_ID=$(CHAIN_ID) \
	HOME=$(CHAIN_DIR) \
	CHAIN_BIN=$(CHAIN_BIN) \
	CHAIN_DATA_DIR=$(CHAIN_DATA_DIR) \
	/bin/bash $(SCRIPT_FILE)

setup:
	$(MAKE) .bash SCRIPT_FILE=setup.sh

clean:
	rm -rf $(CHAIN_DIR)/$(CHAIN_DATA_DIR)

start:
	HOME=$(CHAIN_DIR) $(CHAIN_BIN) start $(ARGS)



