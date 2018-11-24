package godo

import (
	"github.com/spf13/cobra"
	cmdImpl "github.com/thesunnysky/godo/cmd/cmdImpl1"
)

func Run() {
	var addCmd = &cobra.Command{
		Use:     "add [jobs]",
		Aliases: []string{"a"},
		Short:   "a [jobs]",
		Long:    "add [jobs]",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmdImpl.AddCmdImpl(args)
		},
	}

	var delCmd = &cobra.Command{
		Use:     "del [jobs_index]",
		Aliases: []string{"d"},
		Short:   "d [jobs_index]",
		Long:    "del [jobs_index]",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmdImpl.DelCmdImpl(args)
		},
	}

	var listCmd = &cobra.Command{
		Use:     "list [jobs_num]",
		Aliases: []string{"l"},
		Short:   "list jobs",
		Long:    "list jobs",
		Args:    cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cmdImpl.ListCmdImpl(args)
		},
	}

	var cleanCmd = &cobra.Command{
		Use:     "tidy",
		Aliases: []string{"t"},
		Args:    cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cmdImpl.CleanCmdImpl(args)
		},
	}

	var rootCmd = &cobra.Command{Use: "godo"}
	rootCmd.AddCommand(addCmd, delCmd, listCmd, cleanCmd)
	rootCmd.Execute()

}
