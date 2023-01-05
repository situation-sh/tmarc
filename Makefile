PROJECT := github.com/situation-sh/tmarc

BUILD := build
EXTRA := extra

.PHONY: generate build clean fullclean reset
.DEFAULT_GOAL = build

$(EXTRA)/rua.xsd:
	@mkdir -p $(@D)
	wget -qO $@ https://dmarc.org/dmarc-xml/0.1/rua.xsd 

dmarc.go: $(EXTRA)/rua.xsd
	xgen -i $^ -o $@ -l Go -p main >/dev/null

go.mod:
	go mod init $(PROJECT)

go.sum: go.mod
	go mod tidy

$(BUILD)/tmarc: dmarc.go $(shell find . -name "*.go")
	go build -o $@ $^

build: go.sum $(BUILD)/tmarc

generate: dmarc.go

clean:
	rm -rf $(EXTRA) $(BUILD) dmarc.go 

fullclean: clean
	rm -f go.*

reset: fullclean build