{
  description = "Log Watcher with Ollama integration";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools
            libnotify
          ];
          
          shellHook = ''
            echo "Go Log Watcher dev shell"
            echo "Go version: $(go version)"
          '';
        };

        packages.default = pkgs.buildGoModule {
          pname = "log_watcher";
          version = "0.1.0";
          src = ./.;
          vendorHash = null; # No dependencies yet
          
          nativeBuildInputs = [ pkgs.libnotify ];
        };
      }
    );
}
