package cmd

import (
	"fmt"

	"im/pkg/config"
	"im/pkg/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	version string
	date    string

	atom = zap.NewAtomicLevel()

	logLevel string

	rootCmd = &cobra.Command{
		Use:              "im",
		Short:            "im server application",
		Version:          version,
		TraverseChildren: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := atom.UnmarshalText([]byte(logLevel))
			if err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	cobra.OnInitialize(func() {
		config.InitConfig("im")
	})

	logger.InitLogger(atom)

	rootCmd.SetVersionTemplate(fmt.Sprintf("im version %s, build at %s\n", version, date))

	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level")
	_ = viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
