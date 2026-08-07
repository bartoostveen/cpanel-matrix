{
  buildGoModule,
  lib,
  ...
}:

buildGoModule (finalAttrs: {
  pname = "cpanel-matrix";
  version = "1.0.0";

  src = ./.;

  vendorHash = "sha256-GPn1Wt5euaEw1B2e5mYKxDwu+4e4jLDK52F00zE8LcM=";

  meta = {
    description = "Simple, but beatiful Matrix webhook handler for cPanel notifications";
    homepage = "https://git.bartoostveen.nl/bart/cpanel-matrix";
    license = lib.licenses.gpl3Only;
    mainProgram = "cpanel-matrix";
  };
})
