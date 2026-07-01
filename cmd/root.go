// Package cmd defines the Cobra command tree and CLI lifecycle hooks.
package cmd

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/liuguancheng/musicbrainz-cli/internal/apperr"
	"github.com/liuguancheng/musicbrainz-cli/internal/client"
	"github.com/liuguancheng/musicbrainz-cli/internal/output"
	"github.com/liuguancheng/musicbrainz-cli/internal/pagination"
	"github.com/liuguancheng/musicbrainz-cli/internal/runtime"
	searchcmd "github.com/liuguancheng/musicbrainz-cli/cmd/search"
	lookupcmd "github.com/liuguancheng/musicbrainz-cli/cmd/lookup"
)

// Version is embedded in the MusicBrainz User-Agent string.
const Version = "0.1.0"

// RootCmd builds the mbz root command with persistent flags and subcommands.
func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "mbz",
		Short:         "MusicBrainz command-line query tool",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := initOutputMode(cmd); err != nil {
				return err
			}
			return initClient(cmd)
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if runtime.Client != nil {
				return runtime.Client.Close()
			}
			return nil
		},
	}

	rootCmd.PersistentFlags().IntP("limit", "l", pagination.DefaultLimit, "每页条数 (1-100)")
	rootCmd.PersistentFlags().IntP("pageno", "p", pagination.DefaultPageNo, "页码 (从 1 开始)")
	rootCmd.PersistentFlags().String("output", string(output.ModeSimple), "输出模式 (simple|full)")
	rootCmd.PersistentFlags().String("user-agent", "", "HTTP User-Agent")
	rootCmd.PersistentFlags().String("contact", client.DefaultContactURL, "联系方式 URL")
	rootCmd.PersistentFlags().String("api-url", client.DefaultAPIURL, "WS2 API 根地址")

	rootCmd.AddCommand(searchcmd.NewCommand())
	rootCmd.AddCommand(lookupcmd.NewCommand())

	return rootCmd
}

// Execute runs the CLI and returns a process exit code suitable for os.Exit.
func Execute() int {
	cmd := RootCmd()
	if err := cmd.Execute(); err != nil {
		return apperr.WriteAndExitCode(runtime.Stderr, err)
	}
	return apperr.ExitSuccess
}

// initOutputMode parses --output and stores the selected mode in runtime.
func initOutputMode(cmd *cobra.Command) error {
	value, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}
	mode, err := output.ParseMode(value)
	if err != nil {
		return apperr.InvalidArgument(err.Error())
	}
	runtime.OutputMode = mode
	return nil
}

// initClient constructs the MusicBrainz API client unless tests inject a mock.
func initClient(cmd *cobra.Command) error {
	// Tests may pre-set runtime.Client to avoid real HTTP calls.
	if runtime.Client != nil {
		return nil
	}

	userAgent, _ := cmd.Flags().GetString("user-agent")
	contact, _ := cmd.Flags().GetString("contact")
	apiURL, _ := cmd.Flags().GetString("api-url")

	runtime.Client = client.New(client.Config{
		AppName:    client.DefaultAppName,
		Version:    Version,
		ContactURL: contact,
		APIURL:     apiURL,
		UserAgent:  userAgent,
	})
	return nil
}

// SetIO sets stdout/stderr writers. Intended for tests.
func SetIO(out, errOut io.Writer) {
	runtime.Stdout = out
	runtime.Stderr = errOut
}

// SetClient replaces the MusicBrainz client. Intended for tests.
func SetClient(c client.Interface) {
	runtime.Client = c
}

// ResetForTest restores default runtime state. Intended for tests.
func ResetForTest() {
	runtime.ResetForTest()
}
