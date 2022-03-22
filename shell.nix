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

    # example clients to debug with
    (weston.overrideAttrs (old: {
      mesonFlags = old.mesonFlags ++ ["-Ddemo-clients=true" "-Dsimple-clients=all"];
    }))
  ];
}
