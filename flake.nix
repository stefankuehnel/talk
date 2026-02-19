{
  description = "A Nix-flake-based development environment";

  inputs = {
    nixpkgs = {
      url = "github:NixOS/nixpkgs/nixos-unstable";
    };
  };

  outputs =
    {
      nixpkgs,
      ...
    }:
    let
      supportedSystems = nixpkgs.lib.systems.flakeExposed;

      forAllSystems =
        function:
        nixpkgs.lib.genAttrs supportedSystems (
          system:
          function {
            pkgs = nixpkgs.legacyPackages.${system};

            inherit system;
          }
        );

      installNixPackages = pkgs: [
        pkgs.nix
        pkgs.busybox
        pkgs.git
        pkgs.go
        pkgs.go-task
        pkgs.golangci-lint

        pkgs.nixd # Nix Language Server
        pkgs.nixfmt-rfc-style # Nix Formatter
      ];

      installNixFormatter = pkgs: pkgs.nixfmt-tree;
    in
    {
      formatter = forAllSystems ({ pkgs, ... }: installNixFormatter pkgs);

      devShells = forAllSystems (
        { pkgs, ... }:
        {
          default = pkgs.mkShellNoCC {
            packages = installNixPackages pkgs;
          };
        }
      );

      packages = forAllSystems (
        { pkgs, ... }:
        {
          default = pkgs.buildEnv {
            name = "profile";
            paths = installNixPackages pkgs;
          };

          talk = pkgs.callPackage ./talk.nix { };
        }
      );

      apps = forAllSystems (
        { pkgs, ... }:
        {
          "talk" = {
            type = "app";
            program = pkgs.lib.getExe (pkgs.callPackage ./talk.nix { });
          };
        }
      );
    };
}
