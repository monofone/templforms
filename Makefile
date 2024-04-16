LOCALBIN=$(shell pwd)/bin
TEMPL=$(LOCALBIN)/templ
AIR=$(LOCALBIN)/air


$(LOCALBIN):
	mkdir -p $(LOCALBIN)

$(TEMPL):
	test -f $(TEMPL) || GOBIN=$(LOCALBIN) go install github.com/a-h/templ/cmd/templ@latest

$(AIR):
	test -f $(AIR) || GOBIN=$(LOCALBIN) go install github.com/cosmtrek/air@latest

templ-generate: $(TEMPL)
	$(TEMPL) generate -path ./ &

templ-watch: $(TEMPL)
	$(TEMPL) generate -watch -path ./ &

tailwind-watch:
	npx tailwindcss build -o static/css/templ-daisy-ui.css --watch=always < /dev/null &

tests: templ-generate
	go test ./...

watch: templ-watch tailwind-watch

.PHONY: storybook
storybook: $(AIR) watch
	$(AIR)