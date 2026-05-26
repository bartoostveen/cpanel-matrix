{
  description = "Simple, but beatiful Matrix webhook handler for cPanel notifications";

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
        { pkgs, ... }:

        {
          treefmt = {
            programs.nixfmt.enable = true;
            programs.gofmt.enable = true;
          };

          packages.default = pkgs.callPackage ./package.nix { };

          devShells.default = pkgs.mkShell {
            packages = with pkgs; [
              go
              gopls
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
