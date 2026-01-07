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
      rec {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools
          ] ++ pkgs.lib.optionals pkgs.stdenv.isLinux [
            pkgs.libnotify
            pkgs.dpkg
          ];
          
          shellHook = ''
            echo "Go Log Watcher dev shell"
            echo "Go version: $(go version)"
          '';
        };

        packages.default = pkgs.buildGoModule {
          pname = "log_watcher";
          version = "1.0.0";
          src = ./.;
          vendorHash = null; # No dependencies yet
          
          nativeBuildInputs = pkgs.lib.optionals pkgs.stdenv.isLinux [ pkgs.libnotify ];
        };

        packages.deb = pkgs.stdenv.mkDerivation {
          name = "log-watcher-deb";
          src = ./.;
          buildInputs = [ packages.default pkgs.nfpm ];
          buildPhase = ''
            export VERSION="${packages.default.version}"
            mkdir -p result/bin
            cp ${packages.default}/bin/log_watcher result/bin/
            nfpm pkg --packager deb --target .
          '';
          installPhase = ''
            mkdir -p $out
            cp *.deb $out/
          '';
        };

        packages.rpm = pkgs.stdenv.mkDerivation {
          name = "log-watcher-rpm";
          src = ./.;
          buildInputs = [ packages.default pkgs.nfpm ];
          buildPhase = ''
            export VERSION="${packages.default.version}"
            mkdir -p result/bin
            cp ${packages.default}/bin/log_watcher result/bin/
            nfpm pkg --packager rpm --target .
          '';
          installPhase = ''
            mkdir -p $out
            cp *.rpm $out/
          '';
        };

        packages.arch = pkgs.stdenv.mkDerivation {
          name = "log-watcher-arch";
          src = ./.;
          buildInputs = [ packages.default pkgs.nfpm ];
          buildPhase = ''
            export VERSION="${packages.default.version}"
            mkdir -p result/bin
            cp ${packages.default}/bin/log_watcher result/bin/
            nfpm pkg --packager archlinux --target .
          '';
          installPhase = ''
            mkdir -p $out
            cp *.pkg.tar.zst $out/
          '';
        };
      }
    );
}
