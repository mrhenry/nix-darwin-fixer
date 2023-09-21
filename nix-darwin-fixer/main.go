package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

const LaunchDaemonPath = "/Library/LaunchDaemons/be.mrhenry.nix-darwin-fixer.plist"
const NixGcRootPath = "/nix/var/nix/gcroots/nix-darwin-fixer"

const LaunchDaemonPlist = `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>be.mrhenry.nix-darwin-fixer</string>
    <key>ProgramArguments</key>
    <array>
		<string>/bin/sh</string>
		<string>-c</string>
		<string>/bin/wait4path [[NIX_STORE_PATH]]/bin/nix-darwin-fixer &amp;&amp; exec [[NIX_STORE_PATH]]/bin/nix-darwin-fixer fix</string>  
    </array>
	<key>RunAtLoad</key>
    <true/>
	<key>StandardErrorPath</key>
    <string>/var/log/nix-fixer.log</string>
    <key>StandardOutPath</key>
    <string>/var/log/nix-fixer.log</string>
</dict>
</plist>
`

func main() {
	app := &cli.App{
		Version: "0.0.1",
		Commands: []*cli.Command{
			{
				Name:  "install",
				Usage: "Install the LaunchDaemon",
				Action: func(cCtx *cli.Context) error {
					return install()
				},
			},
			{
				Name:  "uninstall",
				Usage: "Uninstall the LaunchDaemon",
				Action: func(cCtx *cli.Context) error {
					return uninstall()
				},
			},
			{
				Name:  "fix",
				Usage: "Run the fixer just once",
				Action: func(cCtx *cli.Context) error {
					// Try to fix the NIX setup
					return tryFix()
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var filesToFix = []string{
	"/etc/zshrc",
	"/etc/bashrc",
}

const SNIPPET = `
# Nix
if [ -e '/nix/var/nix/profiles/default/etc/profile.d/nix-daemon.sh' ]; then
  . '/nix/var/nix/profiles/default/etc/profile.d/nix-daemon.sh'
fi
# End Nix
`

// Try to fix all files
func tryFix() error {
	for _, path := range filesToFix {
		if err := fixFile(path); err != nil {
			return err
		}
	}

	return nil
}

// Fix a single file
func fixFile(path string) error {
	fmt.Printf("Checking %s\n", path)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		fmt.Println("File does not exist, skipping")
		return nil
	}
	if err != nil {
		return err
	}

	text := string(data)

	if strings.Contains(text, "/nix/var/nix/profiles/default/etc/profile.d/nix-daemon.sh") {
		fmt.Println("File already fixed, skipping")
		return nil
	}

	newText := text + SNIPPET

	err = os.WriteFile(path+".backup-before-nix", []byte(text), 0444)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, []byte(newText), 0444)
	if err != nil {
		return err
	}

	fmt.Println("File fixed")

	return nil
}

func install() error {
	self_path := os.Getenv("SELF_NIX_STORE_PATH")
	if self_path == "" {
		return fmt.Errorf("SELF_NIX_STORE_PATH is not set")
	}

	fmt.Printf("Installing LaunchDaemon to %s\n", LaunchDaemonPath)

	// Create the LaunchDaemon plist
	plist := strings.ReplaceAll(LaunchDaemonPlist, "[[NIX_STORE_PATH]]", self_path)

	// Write the LaunchDaemon plist
	err := os.WriteFile(LaunchDaemonPath, []byte(plist), 0444)
	if err != nil {
		return err
	}

	fmt.Printf("Creating GC root at %s\n", NixGcRootPath)

	// Create the GC root
	os.Remove(NixGcRootPath)
	err = os.Symlink(self_path, NixGcRootPath)
	if err != nil {
		return err
	}

	exec.Command("launchctl", "load", LaunchDaemonPath).Run()

	return nil
}

func uninstall() error {
	exec.Command("launchctl", "unload", LaunchDaemonPath).Run()

	fmt.Printf("Removing LaunchDaemon at %s\n", LaunchDaemonPath)
	err := os.Remove(LaunchDaemonPath)
	if os.IsNotExist(err) {
		fmt.Println("LaunchDaemon does not exist, skipping")
		err = nil
	}
	if err != nil {
		return err
	}

	fmt.Printf("Removing GC root at %s\n", NixGcRootPath)
	err = os.Remove(NixGcRootPath)
	if os.IsNotExist(err) {
		fmt.Println("GC root does not exist, skipping")
		err = nil
	}
	if err != nil {
		return err
	}

	return nil
}
