{
  lib,
  buildGoModule,
  ebadf,
  pkg-config,
  sdl2-compat,
  SDL2_image,
  SDL2_mixer,
  SDL2_ttf,
  SDL2_gfx,
  xorg,

}:
buildGoModule {
  pname = "tvman";
  version = "git";
  src = builtins.filterSource ebadf.lib.filterGoSourcesRecursive ./.;
  #CGO_ENABLED = 0;
  nativeBuildInputs = [
    pkg-config
  ];
  buildInputs = [
    sdl2-compat
    SDL2_image
    SDL2_mixer
    SDL2_ttf
    SDL2_gfx
    xorg.libX11
    xorg.libxcb
  ];
  doCheck = false;
  vendorHash = "sha256-Ig2gjxULO/bf3zhZGQnc7l1VCslwf5RbNVxDC7RwdTo=";
  subPackages = [
    "cmd/tvman"
    "cmd/menuman"
  ];
}
