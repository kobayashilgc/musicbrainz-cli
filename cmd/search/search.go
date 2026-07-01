// Package search implements mbz search subcommands for artists and releases.
package search

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/liuguancheng/musicbrainz-cli/internal/apperr"
	"github.com/liuguancheng/musicbrainz-cli/internal/output"
	"github.com/liuguancheng/musicbrainz-cli/internal/pagination"
	releasequery "github.com/liuguancheng/musicbrainz-cli/internal/search"
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
			limit, pageNo, offset, err := paginationFromFlags(cmd)
			if err != nil {
				return err
			}
			query := args[0]
			result, err := runtime.Client.SearchArtists(runtime.Context(), query, limit, offset)
			if err != nil {
				return err
			}
			resp, err := output.ArtistSearch(runtime.OutputMode, query, limit, pageNo, result)
			if err != nil {
				return err
			}
			return output.WriteJSON(runtime.Stdout, resp)
		},
	}
}

func newReleaseCommand() *cobra.Command {
	var artistMBID string
	cmd := &cobra.Command{
		Use:   "release [query]",
		Short: "Search releases",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			limit, pageNo, offset, err := paginationFromFlags(cmd)
			if err != nil {
				return err
			}

			textQuery := ""
			if len(args) > 0 {
				textQuery = args[0]
			}

			query, err := releasequery.BuildReleaseQuery(textQuery, artistMBID)
			if err != nil {
				return err
			}

			result, err := runtime.Client.SearchReleases(runtime.Context(), query, limit, offset)
			if err != nil {
				return err
			}
			resp, err := output.ReleaseSearch(runtime.OutputMode, query, limit, pageNo, result)
			if err != nil {
				return err
			}
			return output.WriteJSON(runtime.Stdout, resp)
		},
	}
	cmd.Flags().StringVar(&artistMBID, "artist-mbid", "", "Filter releases by artist MBID (arid)")
	return cmd
}

// paginationFromFlags reads and validates inherited --limit and --pageno flags.
func paginationFromFlags(cmd *cobra.Command) (limit, pageNo, offset int, err error) {
	limit, err = cmd.Flags().GetInt("limit")
	if err != nil {
		return 0, 0, 0, fmt.Errorf("read limit flag: %w", err)
	}
	pageNo, err = cmd.Flags().GetInt("pageno")
	if err != nil {
		return 0, 0, 0, fmt.Errorf("read pageno flag: %w", err)
	}
	if err := pagination.Validate(limit, pageNo); err != nil {
		return 0, 0, 0, apperr.InvalidArgument(err.Error())
	}
	return limit, pageNo, pagination.Offset(limit, pageNo), nil
}
