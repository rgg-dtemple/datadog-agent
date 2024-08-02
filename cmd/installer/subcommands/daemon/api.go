// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package daemon provides the installer daemon commands.
package daemon

import (
	"fmt"

	"github.com/DataDog/datadog-agent/cmd/installer/command"
	"github.com/DataDog/datadog-agent/comp/core"
	"github.com/DataDog/datadog-agent/comp/core/config"
	log "github.com/DataDog/datadog-agent/comp/core/log/def"
	"github.com/DataDog/datadog-agent/comp/core/secrets"
	"github.com/DataDog/datadog-agent/comp/core/sysprobeconfig/sysprobeconfigimpl"
	"github.com/DataDog/datadog-agent/comp/updater/localapiclient"
	"github.com/DataDog/datadog-agent/comp/updater/localapiclient/localapiclientimpl"
	"github.com/DataDog/datadog-agent/pkg/util/fxutil"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

type cliParams struct {
	command.GlobalParams
	pkg     string
	version string
	catalog string
}

func apiCommands(global *command.GlobalParams) []*cobra.Command {
	setCatalogCmd := &cobra.Command{
		Use:   "set-catalog catalog",
		Short: "Sets the catalog to use",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return experimentFxWrapper(catalog, &cliParams{
				GlobalParams: *global,
				catalog:      args[0],
			})
		},
	}
	installCmd := &cobra.Command{
		Use:     "install package version",
		Aliases: []string{"install"},
		Short:   "Installs a package to the expected version",
		Args:    cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			return experimentFxWrapper(install, &cliParams{
				GlobalParams: *global,
				pkg:          args[0],
				version:      args[1],
			})
		},
	}
	startExperimentCmd := &cobra.Command{
		Use:     "start-experiment package version",
		Aliases: []string{"start"},
		Short:   "Starts an experiment",
		Args:    cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			return experimentFxWrapper(start, &cliParams{
				GlobalParams: *global,
				pkg:          args[0],
				version:      args[1],
			})
		},
	}
	stopExperimentCmd := &cobra.Command{
		Use:     "stop-experiment package",
		Aliases: []string{"stop"},
		Short:   "Stops an experiment",
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return experimentFxWrapper(stop, &cliParams{
				GlobalParams: *global,
				pkg:          args[0],
			})
		},
	}
	promoteExperimentCmd := &cobra.Command{
		Use:     "promote-experiment package",
		Aliases: []string{"promote"},
		Short:   "Promotes an experiment",
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return experimentFxWrapper(promote, &cliParams{
				GlobalParams: *global,
				pkg:          args[0],
			})
		},
	}
	return []*cobra.Command{setCatalogCmd, startExperimentCmd, stopExperimentCmd, promoteExperimentCmd, installCmd}
}

func experimentFxWrapper(f interface{}, params *cliParams) error {
	return fxutil.OneShot(f,
		fx.Supply(core.BundleParams{
			ConfigParams:         config.NewAgentParams(params.ConfFilePath),
			SecretParams:         secrets.NewEnabledParams(),
			SysprobeConfigParams: sysprobeconfigimpl.NewParams(),
			LogParams:            log.ForOneShot("INSTALLER", "off", true),
		}),
		core.Bundle(),
		fx.Supply(params),
		localapiclientimpl.Module(),
	)
}

func catalog(params *cliParams, client localapiclient.Component) error {
	err := client.SetCatalog(params.catalog)
	if err != nil {
		fmt.Println("Error setting catalog:", err)
		return err
	}
	return nil
}

func start(params *cliParams, client localapiclient.Component) error {
	err := client.StartExperiment(params.pkg, params.version)
	if err != nil {
		fmt.Println("Error starting experiment:", err)
		return err
	}
	return nil
}

func stop(params *cliParams, client localapiclient.Component) error {
	err := client.StopExperiment(params.pkg)
	if err != nil {
		fmt.Println("Error stopping experiment:", err)
		return err
	}
	return nil
}

func promote(params *cliParams, client localapiclient.Component) error {
	err := client.PromoteExperiment(params.pkg)
	if err != nil {
		fmt.Println("Error promoting experiment:", err)
		return err
	}
	return nil
}

func install(params *cliParams, client localapiclient.Component) error {
	err := client.Install(params.pkg, params.version)
	if err != nil {
		fmt.Println("Error bootstrapping package:", err)
		return err
	}
	return nil
}
