// Package lookup implements mbz lookup subcommands by MusicBrainz ID.
package lookup

import (
	"github.com/spf13/cobra"
	"go.uploadedlobster.com/mbtypes"

	"github.com/liuguancheng/musicbrainz-cli/internal/apperr"
	"github.com/liuguancheng/musicbrainz-cli/internal/output"
	"github.com/liuguancheng/musicbrainz-cli/internal/runtime"
)

// NewCommand returns the lookup command group.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lookup",
		Short: "Look up MusicBrainz entities by MBID",
	}
	cmd.AddCommand(newArtistCommand())
	cmd.AddCommand(newReleaseCommand())
	cmd.AddCommand(newReleaseGroupCommand())
	return cmd
}

func newArtistCommand() *cobra.Command {
	var includes []string
	cmd := &cobra.Command{
		Use:   "artist <mbid>",
		Short: "Look up an artist",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mbidStr := args[0]
			if err := apperr.ValidateMBID(mbidStr); err != nil {
				return err
			}
			mbid := mbtypes.MBID(mbidStr)
			artist, err := runtime.Client.LookupArtist(runtime.Context(), mbid, includes)
			if err != nil {
				return err
			}
			resp, err := output.ArtistLookup(runtime.OutputMode, mbidStr, artist)
			if err != nil {
				return err
			}
			return output.WriteJSON(runtime.Stdout, resp)
		},
	}
	cmd.Flags().StringArrayVar(&includes, "inc", nil, "附加关联数据")
	return cmd
}

func newReleaseCommand() *cobra.Command {
	var includes []string
	cmd := &cobra.Command{
		Use:   "release <mbid>",
		Short: "Look up a release",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mbidStr := args[0]
			if err := apperr.ValidateMBID(mbidStr); err != nil {
				return err
			}
			mbid := mbtypes.MBID(mbidStr)
			release, err := runtime.Client.LookupRelease(runtime.Context(), mbid, includes)
			if err != nil {
				return err
			}
			resp, err := output.ReleaseLookup(runtime.OutputMode, mbidStr, release)
			if err != nil {
				return err
			}
			return output.WriteJSON(runtime.Stdout, resp)
		},
	}
	cmd.Flags().StringArrayVar(&includes, "inc", nil, "附加关联数据")
	return cmd
}

func newReleaseGroupCommand() *cobra.Command {
	var includes []string
	cmd := &cobra.Command{
		Use:   "releasegroup <mbid>",
		Short: "Look up a release group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mbidStr := args[0]
			if err := apperr.ValidateMBID(mbidStr); err != nil {
				return err
			}
			mbid := mbtypes.MBID(mbidStr)
			releaseGroup, err := runtime.Client.LookupReleaseGroup(runtime.Context(), mbid, includes)
			if err != nil {
				return err
			}
			resp, err := output.ReleaseGroupLookup(runtime.OutputMode, mbidStr, releaseGroup)
			if err != nil {
				return err
			}
			return output.WriteJSON(runtime.Stdout, resp)
		},
	}
	cmd.Flags().StringArrayVar(&includes, "inc", nil, "附加关联数据")
	return cmd
}
