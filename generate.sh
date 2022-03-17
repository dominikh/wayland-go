#!/bin/sh
go run ./cmd/wayland-scanner -prefix=wl_ $(pkg-config --variable=pkgdatadir wayland-server)/wayland.xml > wlclient/protocols/wayland/wayland.go
go run ./cmd/wayland-scanner -mode=server -prefix=wl_ $(pkg-config --variable=pkgdatadir wayland-server)/wayland.xml > wlserver/protocols/wayland/wayland.go

go run ./cmd/wayland-scanner -prefix=xdg_ -i ./wlclient/protocols/wayland $(pkg-config --variable=pkgdatadir wayland-protocols)/stable/xdg-shell/xdg-shell.xml  > ./wlclient/protocols/xdg-shell/xdg-shell.go
go run ./cmd/wayland-scanner -mode=server -prefix=xdg_ -i ./wlclient/protocols/wayland $(pkg-config --variable=pkgdatadir wayland-protocols)/stable/xdg-shell/xdg-shell.xml  > ./wlserver/protocols/xdg-shell/xdg-shell.go
