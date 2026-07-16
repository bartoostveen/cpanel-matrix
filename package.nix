{
  buildGoModule,
  lib,
  ...
}:

buildGoModule (finalAttrs: {
  pname = "cpanel-matrix";
  version = "0.0.1";

  src = ./.;

  vendorHash = "sha256-6UgQ7fmiy/aIc9PG243hpzsLmcrNJKue6Ve3lqNMGrM=";

  meta = {
    description = "Simple, but beatiful Matrix webhook handler for cPanel notifications";
    homepage = "https://git.bartoostveen.nl/bart/cpanel-matrix";
    license = lib.licenses.gpl3Only;
    mainProgram = "cpanel-matrix";
  };
})
