{
  description = "A LaunchDeamon that fixes the zshrc and bashrc files for Nix on Darwin.";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-23.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };

        src = pkgs.lib.cleanSourceWith {
          src = ./.;
          filter = path: type: builtins.match "^.*/(((nix-darwin-fixer|vendor)(/.+)?)|(go[.](mod|sum)))$" path != null;
        };

        bin = pkgs.buildGoModule {
          pname = "nix-darwin-fixer";
          version = "0.0.1";

          subPackages = [ "nix-darwin-fixer" ];
          src = src;

          vendorHash = null;
          CGO_ENABLED = "0";
        };

        binWrapped = pkgs.runCommand "nix-darwin-fixer"
          {
            buildInputs = [ pkgs.makeWrapper ];
          } ''
          mkdir -p $out/bin
          makeWrapper ${bin}/bin/nix-darwin-fixer $out/bin/nix-darwin-fixer \
            --set SELF_NIX_STORE_PATH $out
        '';

        sudoWrapper = pkgs.writeShellScriptBin "nix-darwin-fixer"
          ''
            set -e
            exec sudo ${binWrapped}/bin/nix-darwin-fixer "$@"
          '';
      in
      {
        packages = { } // (pkgs.lib.optionalAttrs pkgs.stdenv.isDarwin {
          default = sudoWrapper;
          bin = binWrapped;
        });

        devShells.default = pkgs.mkShell {
          packages = [
            pkgs.go
          ];
        };
      });
}
