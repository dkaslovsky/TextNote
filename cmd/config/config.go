package config

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dkaslovsky/textnote/pkg/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type commandOptions struct {
	path bool
}

// CreateConfigCmd creates the config subcommand
func CreateConfigCmd() *cobra.Command {
	cmdOpts := commandOptions{}
	cmd := &cobra.Command{
		Use:   "config",
		Short: "show configuration",
		Long:  "displays the application's configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := config.GetConfigFilePath()

			if cmdOpts.path {
				log.Printf("configuration file path: [%s]\n", configPath)
				return nil
			}

			_, err := os.Stat(configPath)
			if os.IsNotExist(err) {
				return fmt.Errorf("cannot find configuration file [%s]", configPath)
			}
			f, err := os.Open(configPath)
			if err != nil {
				return errors.Wrapf(err, "unable to open configuration file [%s]", configPath)
			}
			c, err := io.ReadAll(f)
			if err != nil {
				return errors.Wrapf(err, "unable to read configuration file [%s]", configPath)
			}
			log.Printf("%s", c)
			return nil
		},
	}
	attachOpts(cmd, &cmdOpts)
	return cmd
}

func attachOpts(cmd *cobra.Command, cmdOpts *commandOptions) {
	flags := cmd.Flags()
	flags.BoolVarP(&cmdOpts.path, "path", "p", false, "print path to configuration file")
}
