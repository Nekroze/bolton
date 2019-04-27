package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Nekroze/bolton/pkg/boltons"
	"github.com/Nekroze/bolton/pkg/hardpoints"
	"github.com/Nekroze/bolton/pkg/tui"
)

var home = os.Getenv("HOME")
var libraryPath = filepath.Join(home, ".boltons")

func main() {
	hps := hardpoints.FindHardpoints(".")

	bl, err := boltons.LoadLibrary(libraryPath)
	if err != nil {
		fmt.Println("failed to load library: %s", err)
	}

	tui.Run(hps, bl)
}
