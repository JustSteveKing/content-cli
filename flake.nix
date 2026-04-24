{
  description = "A CLI tool for managing content";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        # You can update this version or pull it from a file
        version = "0.1.0";
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "content";
          inherit version;
          src = ./.;

          # Run `nix build` and replace this with the hash provided in the error message
          vendorHash = "sha256-+hrjivOGODq9F8ILetJZj6+AvJBHOk0OJZERY4MB+Jk=";

          ldflags = [
            "-X github.com/juststeveking/content-cli/cmd.Version=${version}"
          ];

          subPackages = [ "." ];

          meta = with pkgs.lib; {
            description = "A CLI tool for managing content";
            homepage = "https://github.com/juststeveking/content-cli";
            license = licenses.mit;
            maintainers = [ ];
          };
        };

        apps.default = utils.lib.mkApp {
          drv = self.packages.${system}.default;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            goreleaser
            gnumake
          ];
        };
      });
}
