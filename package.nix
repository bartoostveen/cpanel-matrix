{
  buildGoModule,
  lib,
}:

buildGoModule (finalAttrs: {
  pname = "cpanel-matrix";
  version = "1.0.6";

  src = ./.;

  vendorHash = "sha256-QiUiu2T0zIPSb978NkxL/n9qGlG5KkBSqR6oPvYnKX8=";

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
