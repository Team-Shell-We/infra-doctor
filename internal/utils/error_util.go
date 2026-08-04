package utils

import (
    "fmt"
)

func PrintError() {
    fmt.Println(`
        Error: unknown command "wrong-command" for "infra-doctor"

		Run 'infra-doctor help' to see available commands.`)
}
