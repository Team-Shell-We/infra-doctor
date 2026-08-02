package utils

import (
    "fmt"
)

func PrintHelp() {
    fmt.Println(`
        Infra Doctor v1.0.0

        Usage
        infra-doctor <command> [arguments]

        Core Commands

        doctor
            Analyze project readiness

        visualize architecture
            Show infrastructure architecture

        visualize flow
            Show deployment flow

        explain <topic>
            Explain infrastructure concepts

        recommend
            Recommend deployment strategy

        generate <target>
            Generate configuration files

        export
            Export project report

        Run

        infra-doctor help <command>

        for detailed information.`)
}
