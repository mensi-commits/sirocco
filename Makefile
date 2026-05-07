# Sirocco Makefile

APP_NAME_SWITCH=switch
APP_NAME_WORKER=worker

SWITCH_DIR=./switch
WORKER_DIR=./worker

SWITCH_PORT=8080
WORKER_PORT_1=9001
WORKER_PORT_2=9002

.PHONY: all switch worker run clean

# ------------------------
# Run everything
# ------------------------
all: worker switch

# ------------------------
# Run switch (router)
# ------------------------
switch:
	@echo "Starting Sirocco Switch on :$(SWITCH_PORT)"
	@cd $(SWITCH_DIR) && go run main.go

# ------------------------
# Run worker 1
# ------------------------
worker:
	@echo "Starting Worker 1 on :$(WORKER_PORT_1)"
	@cd $(WORKER_DIR) && PORT=$(WORKER_PORT_1) go run main.go

# ------------------------
# Optional: run worker 2
# ------------------------
worker2:
	@echo "Starting Worker 2 on :$(WORKER_PORT_2)"
	@cd $(WORKER_DIR) && PORT=$(WORKER_PORT_2) go run main.go

# ------------------------
# Run full system (recommended)
# ------------------------
run:
	@echo "Starting full Sirocco system..."

	@make worker & \
	make worker2 & \
	sleep 2 && \
	make switch

# ------------------------
# Clean processes (optional helper)
# ------------------------
clean:
	@echo "Killing Go processes..."
	@pkill -f "go run"



	https://gist.github.com/CodingKoopa/3b30afe8c91e3950f6b124cd2abe3b6b