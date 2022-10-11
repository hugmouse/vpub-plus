{
  description = "vpub-plus-plus";

  # Nixpkgs / NixOS version to use.
  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let

      # System types to support.
      supportedSystems =
        [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system:
        import nixpkgs {
          inherit system;
          overlays = [ self.overlays.default ];
        });
    in
    {

      # A Nixpkgs overlay.
      overlays.default = final: prev: {
        vpub-plus-plus = with final;

          buildGoModule rec {

            pname = "vpub-plus-plus";
            version = "0.1";

            src = ./.;
            vendorSha256 =
              "sha256-1mTXrD/U8KGUlKe40QzYe1xaHaEWjn6u/AvkXohqpd0=";
            # subPackages = [ "." ];

            meta = with lib; {
              description = "TODO";
              homepage = "https://github.com/pinpox/vpub-plus-plus";
              license = licenses.gpl3;
              maintainers = with maintainers; [ pinpox ];
              # platforms = platforms.linux;
            };
          };
      };

      # Package
      packages = forAllSystems (system: {
        inherit (nixpkgsFor.${system}) vpub-plus-plus;
        default = self.packages.${system}.vpub-plus-plus;
      });

      # Nixos module
      nixosModules.vpub-plus-plus = { pkgs, lib, config, ... }:
        with lib;
        let cfg = config.services.vpub-plus-plus;
        in {
          imports = [ ./module.nix ];
          config = mkIf cfg.enable {
            nixpkgs.overlays = [ self.overlays.default ];
          };
        };

      # Tests run by 'nix flake check' and by Hydra.
      checks = forAllSystems
        (system:
          with nixpkgsFor.${system};

          lib.optionalAttrs stdenv.isLinux {
            # A VM test of the NixOS module.
            vmTest =
              with import (nixpkgs + "/nixos/lib/testing-python.nix")
                {
                  inherit system;
                };

              (makeTest {
                name = "vpub-plus-plus-test";
                nodes = {
                  server = {
                    imports = [ self.nixosModules.vpub-plus-plus ];
                    services.vpub-plus-plus = {
                      port = "1234";
                      envFile = "env.example";
                      enable = true;
                      title = "test-forum";
                    };
                  };
                };

                testScript =
                  ''
                    start_all()
                    server.wait_for_unit("multi-user.target")
                    server.wait_for_unit("vpub-plus-plus.service")
                  '';
              }).test;
          }
        );
    };
}
