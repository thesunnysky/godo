package godo

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var g = &GitRepo{repoPath: ClientConfig.GithubRepo}

func Run() {

	var addCmd = &cobra.Command{
		Use:     "add [jobs]",
		Aliases: []string{"a"},
		Short:   "a [jobs]",
		Long:    "add [jobs]",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			AddCmdImpl(args)
		},
	}

	var delCmd = &cobra.Command{
		Use:     "del [jobs_index]",
		Aliases: []string{"d"},
		Short:   "d [jobs_index]",
		Long:    "del [jobs_index]",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			DelCmdImpl(args)
		},
	}

	var listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "list jobs",
		Long:    "list jobs",
		Args:    cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			remote, err := cmd.Flags().GetBool("remote")
			if err != nil {
				fmt.Println("parse flag error:%s\n", err)
				os.Exit(-1)
			}
			ListTasks(args, remote)
		},
	}
	listCmd.Flags().BoolP("remote", "r", false, "remote")

	var cleanCmd = &cobra.Command{
		Use:     "tidy",
		Aliases: []string{"t"},
		Args:    cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			TidyCmdImpl(args)
		},
	}

	//push-server
	var pushServerCmd = &cobra.Command{
		Use:  "push-server",
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			PushServerCmd(args)
		},
	}

	//pull-server
	var pullServerCmd = &cobra.Command{
		Use:  "pull-server",
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			PullServerCmd(args)
		},
	}

	//godo push
	var pushGitCmd = &cobra.Command{
		Use:  "push",
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			PushGitCmd(args)
		},
	}

	//godo pull
	var pullGitCmd = &cobra.Command{
		Use:  "pull",
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			gitCmArgs := []string{"pull"}
			_ = g.GitCmd(gitCmArgs)
		},
	}

	var gitCmd = &cobra.Command{
		Use:     "git",
		Aliases: []string{"g"},
		Args:    cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			_ = g.GitCmd(os.Args[2:])
		},
	}

	var gitAddCmd = &cobra.Command{
		Use:  "add",
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var gitArgs []string
			if addAll, _ := cmd.Flags().GetBool("all"); addAll {
				gitArgs = []string{"add", "-A"}
			} else {
				gitArgs = args
			}
			_ = g.GitCmd(gitArgs)
		},
	}
	gitAddCmd.Flags().BoolP("all", "A", false, "Add all")

	var gitPushCmd = &cobra.Command{
		Use:  "push",
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			_ = g.GitCmd(args)
		},
	}

	var gitPullCmd = &cobra.Command{
		Use:  "pull",
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			gitArgs := []string{"pull", "--rebase"}
			_ = g.GitCmd(gitArgs)
		},
	}

	var commitMsg string
	var gitCommitCmd = &cobra.Command{
		Use:     "commit",
		Aliases: []string{"cm"},
		Args:    cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			gitArgs := []string{"commit", "-m", "\"" + commitMsg + "\""}
			_ = g.GitCmd(gitArgs)
		},
	}
	gitCommitCmd.Flags().StringVarP(&commitMsg, "message", "m", "", "commit message")

	var gitAddAndCommitCmd = &cobra.Command{
		Use:  "cm",
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			message := strings.Join(args, " ")
			gitArgs := []string{"commit", "-am", "\"" + message + "\""}
			_ = g.GitCmd(gitArgs)
		},
	}

	//godo update
	var updateCmd = &cobra.Command{
		Use:     "update",
		Aliases: []string{"u"},
		Args:    cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			UpdateCmd(args)
		},
	}

	gitCmd.AddCommand(gitAddCmd, gitPushCmd, gitCommitCmd, gitPullCmd, gitAddAndCommitCmd)

	var rootCmd = &cobra.Command{Use: "godo"}
	rootCmd.AddCommand(addCmd, delCmd, listCmd, cleanCmd, pushServerCmd, pullServerCmd,
		gitCmd, pushGitCmd, pullGitCmd, updateCmd)
	_ = rootCmd.Execute()
}
