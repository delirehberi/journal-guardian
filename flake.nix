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
            libnotify
            dpkg
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
          
          nativeBuildInputs = [ pkgs.libnotify ];
        };

        packages.deb = pkgs.stdenv.mkDerivation {
          name = "log-watcher-deb";
          src = ./.;
          
          buildInputs = [ packages.default pkgs.dpkg ];
          
          buildPhase = ''
            mkdir -p deb/usr/bin
            mkdir -p deb/DEBIAN
            mkdir -p deb/usr/lib/systemd/user
            
            cp ${packages.default}/bin/log_watcher deb/usr/bin/
            cp debian/control deb/DEBIAN/
            cp debian/postinst deb/DEBIAN/
            chmod 755 deb/DEBIAN/postinst
            cp log_watcher.service deb/usr/lib/systemd/user/

            # Update Version in control file
            sed -i "s/Version: .*/Version: ${packages.default.version}/g" deb/DEBIAN/control
          '';
          
          installPhase = ''
            mkdir -p $out
            dpkg-deb --build deb $out/log_watcher_${packages.default.version}_amd64.deb
          '';
        };
      }
    );
}
