{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  buildInputs = with pkgs; [
    libGL
    libdrm
    libinput
    pkg-config
    wayland
    wayland.debug # debug symbols
    wayland-protocols
    wayland-utils
    hello-wayland # simple test client
  ];
}
