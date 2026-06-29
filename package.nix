{
  buildGoModule,
  lib,
  ...
}:

buildGoModule (finalAttrs: {
  pname = "cpanel-matrix";
  version = "0.0.1";

  src = ./.;

  vendorHash = "sha256-h3F4wo9zTuSUOasP6oKYQavfORhRNqJrbEa1rH0MCkE=";

  meta = {
    description = "Simple, but beatiful Matrix webhook handler for cPanel notifications";
    homepage = "https://git.bartoostveen.nl/bart/cpanel-matrix";
    license = lib.licenses.gpl3Only;
    mainProgram = "cpanel-matrix";
  };
})
