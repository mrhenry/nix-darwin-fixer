# Keep Nix working on macOS after system updates

This is a fix for when Nix breaks after macOS system updates.
It installs a LaunchDaemon that will run on system boot and fix the
broken `/etc/zshrc` and `/etc/bashrc` files. You can also run it manually.

## How to use the fixer

### Fix manually without installing.

```sh
# If you already had Nix installed and it is already broken,
# you will need to first source nix into your current shell.
. '/nix/var/nix/profiles/default/etc/profile.d/nix-daemon.sh'

# Then you can run the fixer. Note this will prompt for your password
# as it needs to run as root.
nix run github:mrhenry/nix-darwin-fixer -- fix

# Now open a new shell and you should be good to go.
nix --version
```

### Install the fixer as a LaunchDaemon

```sh
# If you already had Nix installed and it is already broken,
# you will need to first source nix into your current shell.
. '/nix/var/nix/profiles/default/etc/profile.d/nix-daemon.sh'

# Then you can run the fixer. Note this will prompt for your password
# as it needs to run as root.
nix run github:mrhenry/nix-darwin-fixer -- install

# Now open a new shell and you should be good to go.
nix --version
```

### Uninstall the fixer

```sh
# This will remove all files installed by the installer.
nix run github:mrhenry/nix-darwin-fixer -- uninstall
```

## Reading Material

- https://developer.apple.com/library/archive/documentation/MacOSX/Conceptual/BPSystemStartup/Chapters/CreatingLaunchdJobs.html#//apple_ref/doc/uid/10000172i-SW7-BCIEDDBJ
- https://github.com/NixOS/nix/issues/3616
- https://checkoway.net/musings/nix/#fixing-shell-integration
- https://github.com/DeterminateSystems/nix-installer/issues/593
