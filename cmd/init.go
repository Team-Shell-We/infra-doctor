package cmd

import (
	"fmt"
	"os"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
	"github.com/Team-Shell-We/infra-doctor/internal/utils"
	"github.com/spf13/cobra"
)

func initCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize Infra Doctor",
		Args:  cobra.NoArgs,

		RunE: func(cmd *cobra.Command, args []string) error {
			currentDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf(
					"failed to get current directory: %w",
					err,
				)
			}

			detection, err := project.DetectSpringBoot(currentDir)
			if err != nil {
				return err
			}

			if !detection.Detected {
				return fmt.Errorf(
					"Spring Boot project not found in %s; "+
						"build.gradle, build.gradle.kts, or pom.xml with a Spring Boot marker is required",
					currentDir,
				)
			}

			fmt.Println("Spring Boot project:", currentDir)

			result, err := utils.Initialize(currentDir)
			if err != nil {
				return err
			}

			for _, createdPath := range result.CreatedPaths {
				fmt.Println("Created:", createdPath)
			}

			return nil
		},
	}
}

func init() {
	rootCmd.AddCommand(initCommand())
}
