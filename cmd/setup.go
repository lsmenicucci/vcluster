/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/lsmenicucci/vcluster/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup [config-file]",
	Short: "Setup virtual cluster based on configuration file",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: run,
	Args: cobra.MatchAll(cobra.ExactArgs(1)),
}


func init() {
	rootCmd.AddCommand(setupCmd)
	
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func run(cmd *cobra.Command, args []string){
	log := logrus.WithField("config-file", args[0])

	cfstat, err := os.Stat(args[0])
	if (os.IsNotExist(err)){
		log.Fatalf("File does not exists")
	}

	if (err != nil){
		log.WithError(err).Fatalf("Could not verify file %s", args[0])
	}

	if (cfstat.IsDir()){
		log.Fatalf("%s is a directory", args[0])
	}

	// cfg := &pkg.ClusterConfig{}
	// err = cfg.LoadFromFile(args[0])
	// if (err != nil){
	// 	log.Logger.WithError(err).Fatal("Error loading configuration file")
	// }

	l, err := pkg.DialLibvirt()
	if (err != nil){
		log.WithError(err).Fatalf("Could not dial to libvirt's daemon")
	}

	pkg.Test(l)
}
