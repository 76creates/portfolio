package cmd

import (
	"context"
	"fmt"
	"github.com/76creates/portfolio/dsn.onl/server"
	"github.com/spf13/cobra"
	"time"
)

var serverCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"s", "serv", "srv"},
	Short:   "Subcommand for managing the portfolio web server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var serverStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the portfolio web server",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		appConf := server.Conf{
			IP:           "0.0.0.0",
			Port:         "3000",
			Concurrency:  100,
			BodyLimit:    1000000,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
			Logger: server.LoggerConf{
				File:        "server.log",
				Level:       "debug",
				LogCaller:   true,
				LogToStdout: true,
				StaticExtraFields: map[string]string{
					"app": "portfolio",
				},
			},
		}
		app := appConf.GetApp(ctx)

		return app.Listen(fmt.Sprintf("%s:%s", appConf.IP, appConf.Port))
	},
}

type Conf struct {
	IP           string        `mapstructure:"ip" yaml:"ip"`
	Port         string        `mapstructure:"port" yaml:"port"`
	Concurrency  int           `mapstructure:"concurency" yaml:"concurency"`
	BodyLimit    int           `mapstructure:"bodyLimit" yaml:"bodyLimit"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout" yaml:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout" yaml:"writeTimeout"`
	IdleTimeout  time.Duration `mapstructure:"idleTimeout" yaml:"idleTimeout"`
	Logger       LoggerConf    `mapstructure:"logger" yaml:"logger"`
}

type LoggerConf struct {
	File              string            `mapstructure:"file" yaml:"file"`
	Level             string            `mapstructure:"level" yaml:"level"`
	LogCaller         bool              `mapstructure:"logCaller" yaml:"logCaller"`
	LogToStdout       bool              `mapstructure:"logToStdout" yaml:"logToStdout"`
	StaticExtraFields map[string]string `mapstructure:"staticExtraFields" yaml:"staticExtraFields"`
}

func init() {
	serverCmd.AddCommand(serverStartCmd)
}
