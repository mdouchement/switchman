package main

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const version = "0.0.1"

var (
	ver    bool
	cnf    string
	config Config
)

func main() {
	c := &cobra.Command{
		Use:   "switchman",
		Short: "Proxy for dev stack",
		Long:  "Proxy for dev stack",
		RunE:  action,
	}
	c.Flags().BoolVarP(&ver, "version", "v", false, "Print the version")
	c.Flags().StringVarP(&cnf, "config", "c", "", "Configuration file")

	if err := c.Execute(); err != nil {
		fmt.Println(err)
	}
}

func action(c *cobra.Command, args []string) error {
	if ver {
		fmt.Printf("Switchman %s\n", version)
		return nil
	}

	if cnf == "" {
		return errors.New("No configuration file provided")
	}

	data, err := ioutil.ReadFile(cnf)
	if err != nil {
		return errors.Wrap(err, "config:")
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return errors.Wrap(err, "parsing config:")
	}

	engine := echo.New()
	engine.Use(middleware.Recover())

	Dispatch(engine, config.Rules)

	if err := engine.Start(config.Listen); err != nil {
		return errors.Wrap(err, "server:")
	}
	return nil
}
