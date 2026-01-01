{
  description = "Simple dev shell with libglpk-dev";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      # Change this to your system if needed
      system = "x86_64-linux";

      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        name = "glpk-dev-shell";

        packages = with pkgs; [
          glpk         # includes both lib + headers (+ cli tool)
        ];

        hardeningDisable = [ "fortify" ];   # ← this is the magic line

        # Optional - makes pkg-config find glpk easily
        shellHook = ''
          echo "GLPK development environment loaded"
          echo "pkg-config --cflags --libs glpk"
        '';
      };
    };
}