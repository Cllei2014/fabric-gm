/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tw-bc-group/fabric-gm/peer/chaincode"
	"github.com/tw-bc-group/fabric-gm/peer/channel"
	"github.com/tw-bc-group/fabric-gm/peer/clilogging"
	"github.com/tw-bc-group/fabric-gm/peer/common"
	"github.com/tw-bc-group/fabric-gm/peer/node"
	"github.com/tw-bc-group/fabric-gm/peer/version"
)

// The main command describes the service and
// defaults to printing the help message.
var mainCmd = &cobra.Command{
	Use: "peer"}

func main() {

	// For environment variables.
	viper.SetEnvPrefix(common.CmdRoot)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// Set GM Provider by:
	// CORE_GM_PROVIDER=ALIYUN_KMS
	// CORE_GM_PROVIDER=ZHONGHUAN
	viper.SetDefault("gm.provider", "SW")

	// Define command-line flags that are valid for all peer commands and
	// subcommands.
	mainFlags := mainCmd.PersistentFlags()

	mainFlags.String("logging-level", "", "Legacy logging level flag")
	viper.BindPFlag("logging_level", mainFlags.Lookup("logging-level"))
	mainFlags.MarkHidden("logging-level")

	mainCmd.AddCommand(version.Cmd())
	mainCmd.AddCommand(node.Cmd())
	mainCmd.AddCommand(chaincode.Cmd(nil))
	mainCmd.AddCommand(clilogging.Cmd(nil))
	mainCmd.AddCommand(channel.Cmd(nil))

	// On failure Cobra prints the usage message and error string, so we only
	// need to exit with a non-0 status
	if mainCmd.Execute() != nil {
		os.Exit(1)
	}
}
