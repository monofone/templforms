LOCALBIN=$(shell pwd)/bin
TEMPL=$(LOCALBIN)/templ

$(LOCALBIN):
	mkdir -p $(LOCALBIN)

$(TEMPL):
	test -f $(TEMPL) || GOBIN=$(LOCALBIN) go install github.com/a-h/templ/cmd/templ@latest

templ-generate: $(TEMPL)
	$(TEMPL) generate -watch -path ./ &

templ-watch: $(TEMPL)
	$(TEMPL) generate -path ./ &

tests: templ-generate
	go test ./...
