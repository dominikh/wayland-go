{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  buildInputs = with pkgs; [
    libGL
    libdrm
    libinput
    pkg-config
    wayland
    wayland.debug
    wayland-protocols
    wayland-utils
  ];
}
