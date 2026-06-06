package main

import (
	"fmt"
	"log"
	"os"

	"github.com/deluxesande/disk-cleaner/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

const asciiBanner = `
 ____  _     _      ____ _                           
|  _ \(_)___| | __ / ___| | ___  __ _ _ __   ___ _ __ 
| | | | / __| |/ /| |   | |/ _ \/ _  | '_ \ / _ \ '__|
| |_| | \__ \   < | |___| |  __/ (_| | | | |  __/ |   
|____/|_|___/_|\_\ \____|_|\___|\__,_|_| |_|\___|_|   
                                                      
    High-Performance System Cleanup Engine
`

func main() {
	fmt.Print(asciiBanner)

	p := tea.NewProgram(ui.InitialModel())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
