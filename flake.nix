{
  description = "Simple Matrix webhook handler for cPanel notifications";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";

    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "aarch64-darwin"
      ];

      imports = [
        inputs.treefmt-nix.flakeModule
      ];

      perSystem =
        { pkgs, self', ... }:

        {
          treefmt = {
            programs.nixfmt.enable = true;
            programs.gofmt.enable = true;
          };

          packages = {
            default = pkgs.callPackage ./package.nix { };
            cpanel-matrix = self'.packages.default;
            static = self'.packages.default.overrideAttrs {
              GCO_ENABLED = 0;
              ldflags = [
                "-linkmode external"
                "-extldflags '-static -L${pkgs.glibc.static}/lib'"
              ];
            };
            towncrier-build = pkgs.writeShellApplication {
              name = "towncrier-build";
              runtimeInputs = [ pkgs.python314Packages.towncrier ];
              text = ''
                towncrier build --version "${self'.packages.default.version}" --yes "$@"
              '';
            };
            get-changelog = pkgs.writeShellApplication {
              name = "get-changelog";
              runtimeInputs = [
                pkgs.coreutils
                pkgs.git
              ];
              text = ''
                git diff HEAD~ HEAD -- CHANGELOG.md | grep '^[+]' | sed 's/^+//' | tail -n +3
              '';
            };
            tag-release = pkgs.writeShellApplication {
              name = "tag-release";
              runtimeInputs = [
                pkgs.git
                self'.packages.towncrier-build
              ];
              text = ''
                towncrier-build
                git commit -m "chore: Release ${self'.packages.default.version}"
                git tag 'v${self'.packages.default.version}' -m "Release ${self'.packages.default.version}"
              '';
            };
          };

          devShells.default = pkgs.mkShell {
            packages = with pkgs; [
              go
              gopls
              python314Packages.towncrier
            ];

            shellHook = ''
              export GOPATH=$PWD/.gopath
              export PATH=$GOPATH/bin:$PATH
              mkdir -p $GOPATH
              go telemetry off # This doesn't restore original state
            '';
          };
        };
    };
}
