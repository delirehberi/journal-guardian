.PHONY: build run nix-build nix-shell check clean

PROGName = log_watcher

build:
	CGO_ENABLED=0 go build -ldflags "-s -w" -o $(PROGName) .

run:
	go run .

nix-build:
	nix build .

nix-shell:
	nix develop

check:
	nix flake check

clean:
	rm -f $(PROGName)
