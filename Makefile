BIN := bin
MAIN := $(BIN)/main

AIR := $(BIN)/air
AIR_CONF := .air.toml

$(BIN):
	mkdir -p $(BIN)

$(AIR): $(BIN)
	curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
	$(AIR) -v

air: $(AIR)

run-dev: $(AIR)
	source .env; $(AIR) -c $(AIR_CONF)

rm-air:
	rm -rf $(AIR)

$(MAIN):
	go build -o $(MAIN) main.go
	chmod +x $(MAIN)

build: $(MAIN)

rm-bin:
	rm -rf $(MAIN)

serve:
	source .env; $(MAIN)
