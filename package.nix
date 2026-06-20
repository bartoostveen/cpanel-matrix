{
  buildGoModule,
  lib,
  ...
}:

buildGoModule (finalAttrs: {
  pname = "cpanel-matrix";
  version = "0.0.1";

  src = ./.;

  vendorHash = "sha256-IY20o8ch/WcJNSH1s2JgmdHIVYZykTri8G9TeMzKrUE=";

  meta = {
    description = "Simple, but beatiful Matrix webhook handler for cPanel notifications";
    homepage = "https://git.bartoostveen.nl/bart/cpanel-matrix";
    license = lib.licenses.gpl3Only;
    mainProgram = "cpanel-matrix";
  };
})
