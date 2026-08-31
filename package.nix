{
  buildGoModule,
  lib,
  ...
}:

buildGoModule (finalAttrs: {
  pname = "cpanel-matrix";
  version = "1.0.3";

  src = ./.;

  vendorHash = "sha256-PJ6XiQR1IGsXEMtmB7wjpl22jx0VZaBULfjWfuj6XPk=";

  meta = {
    description = "Simple Matrix webhook handler for cPanel notifications";
    homepage = "https://git.bartoostveen.nl/bart/cpanel-matrix";
    changelog = "https://git.bartoostveen.nl/bart/cpanel-matrix/src/tag/v${finalAttrs.version}/CHANGELOG.md";
    license = lib.licenses.gpl3Only;
    maintainers = with lib.maintainers; [ bartoostveen ];
    mainProgram = "cpanel-matrix";
    platforms = lib.platforms.all;
  };
})
