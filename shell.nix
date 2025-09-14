{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  nativeBuildInputs = with pkgs.buildPackages; [ go xclip pkg-config ];

  buildInputs = with pkgs; [
    xorg.libX11
    xorg.libXext
    xorg.libXrender
    xorg.libxcb
    xorg.libXau
    xorg.libXdmcp
    xorg.libXrandr
    xorg.libXxf86vm
    xorg.libXi
    xorg.libXcursor
    xorg.libXinerama

    libGL
  ];
}
