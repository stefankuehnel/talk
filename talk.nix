{
  lib,
  buildGoModule,
}:

buildGoModule rec {
  pname = "talk";
  version = "1.0.0";

  # See: https://nix.dev/guides/best-practices#reproducible-source-paths
  src = builtins.path {
    path = ./.;
    name = "talk";
  };

  # The vendorHash is a SHA-256 hash of the vendored dependencies, used for reproducible builds.
  # 
  # How to update:
  #   1. Run 'go mod vendor'
  #   2. Run 'nix hash path --sri vendor/' to get the new hash
  vendorHash = "sha256-9jK3jKbFp+5WSQfMbNzwIB55bC5KScZOaFHItffTF00=";

  # Inject version at build time via ldflags
  ldflags = [
    "-s"
    "-w"
    "-X=stefanco.de/talk/cmd.version=${version}"
  ];

  meta = with lib; {
    description = "A command-line interface in Go for sending messages to Nextcloud Talk.";
    homepage = "https://github.com/stefankuehnel/talk";
    license = licenses.gpl3;
    mainProgram = "talk";
    maintainers = [
      {
        name = "Stefan Kühnel";
        email = "git@stefankuehnel.com";
      }
    ];
  };
}
