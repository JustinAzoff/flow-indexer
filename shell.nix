let
  pkgs = import <nixpkgs> {};
in
pkgs.mkShell {
  buildInputs = [
    pkgs.go
    pkgs.nfdump
    pkgs.argus pkgs.argus-clients
  ];
}
