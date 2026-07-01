// Package search implements mbz search subcommands for artists and releases.
package search

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/liuguancheng/musicbrainz-cli/internal/apperr"
	"github.com/liuguancheng/musicbrainz-cli/internal/output"
	"github.com/liuguancheng/musicbrainz-cli/internal/pagination"
	"github.com/liuguancheng/musicbrainz-cli/internal/runtime"
)

// NewCommand returns the search command group.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search MusicBrainz entities",
	}
	cmd.AddCommand(newArtistCommand())
	cmd.AddCommand(newReleaseCommand())
	return cmd
}

func newArtistCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "artist <query>",
		Short: "Search artists",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			limit, offset, err := paginationFromFlags(cmd)
			if err != nil {
				return err
			}
			query := args[0]
			result, err := runtime.Client.SearchArtists(runtime.Context(), query, limit, offset)
			if err != nil {
				return err
			}
			resp, err := output.ArtistSearch(runtime.OutputMode, query, limit, offset, result)
			if err != nil {
				return err
			}
			return output.WriteJSON(runtime.Stdout, resp)
		},
	}
}

func newReleaseCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "release <query>",
		Short: "Search releases",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			limit, offset, err := paginationFromFlags(cmd)
			if err != nil {
				return err
			}
			query := args[0]
			result, err := runtime.Client.SearchReleases(runtime.Context(), query, limit, offset)
			if err != nil {
				return err
			}
			resp, err := output.ReleaseSearch(runtime.OutputMode, query, limit, offset, result)
			if err != nil {
				return err
			}
			return output.WriteJSON(runtime.Stdout, resp)
		},
	}
}

// paginationFromFlags reads and validates inherited --limit and --offset flags.
func paginationFromFlags(cmd *cobra.Command) (int, int, error) {
	limit, err := cmd.Flags().GetInt("limit")
	if err != nil {
		return 0, 0, fmt.Errorf("read limit flag: %w", err)
	}
	offset, err := cmd.Flags().GetInt("offset")
	if err != nil {
		return 0, 0, fmt.Errorf("read offset flag: %w", err)
	}
	if err := pagination.Validate(limit, offset); err != nil {
		return 0, 0, apperr.InvalidArgument(err.Error())
	}
	return limit, offset, nil
}
