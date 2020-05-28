let
  pkgs = import <nixpkgs> {};
in
pkgs.mkShell {
  buildInputs = with pkgs;[
    go
    #For testing
    nfdump
    argus argus-clients
  ];
}
