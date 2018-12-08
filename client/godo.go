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
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			AddCmdImpl(args)
		},
	}

	var delFileArgs *[]string
	var delCmd = &cobra.Command{
		Use:     "del",
		Aliases: []string{"d"},
		Args:    cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if len(*delFileArgs) != 0 {
				DelBackupTaskFile(*delFileArgs)
			} else {
				DelCmdImpl(args)
			}
		},
	}
	delFileArgs = delCmd.Flags().StringArrayP("file", "f", []string{}, "delete backup file")

	// godo ls
	var listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Args:    cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			flagCount := cmd.Flags().NFlag()
			if flagCount > 1 {
				fmt.Printf("more than two flag are set, exit...")
				os.Exit(-1)
			}

			remote, err := cmd.Flags().GetBool("remote")
			if err != nil {
				fmt.Printf("parse flag error:%s\n", err)
				os.Exit(-1)
			}

			showBackupFileList, err := cmd.Flags().GetBool("backup")
			if err != nil {
				fmt.Printf("parse flag error:%s\n", err)
				os.Exit(-1)
			}

			backupFileNo, err := cmd.Flags().GetInt("file")
			if err != nil {
				fmt.Printf("parse flag error:%s\n", err)
				os.Exit(-1)
			}

			if remote {
				//print remote task file
				ListRemoteTasks(args)
			} else if showBackupFileList {
				//show backup file list
				ListBackupFiles(args)
			} else if backupFileNo > 0 {
				//print local task file
				ListBackupTasks(backupFileNo)
			} else {
				//print local task file
				ListLocalTasks(args)
			}

		},
	}
	listCmd.Flags().BoolP("remote", "r", false, "list tasks that recorded in git repo")
	listCmd.Flags().IntP("file", "f", 0, "backup file list or backup file content")
	listCmd.Flags().BoolP("backup", "b", false, "show backup file list")

	var cleanCmd = &cobra.Command{
		Use:     "tidy",
		Aliases: []string{"t"},
		Args:    cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			TidyCmdImpl(args)
		},
	}

	var backupCmd = &cobra.Command{
		Use:     "backup",
		Aliases: []string{"b"},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			BackupTaskFile()
		},
	}

	var recoverCmd = &cobra.Command{
		Use:     "recover",
		Aliases: []string{"r"},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			RecoverTaskFile()
		},
	}

	var cleanBackupFileCmd = &cobra.Command{
		Use:     "clean",
		Aliases: []string{"c"},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			//todo
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
		gitCmd, pushGitCmd, pullGitCmd, updateCmd, backupCmd, recoverCmd, cleanBackupFileCmd)
	_ = rootCmd.Execute()
}
