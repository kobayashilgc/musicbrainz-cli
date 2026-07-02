package runtime

import (
	"context"
	"testing"
	"time"

	"go.uploadedlobster.com/mbtypes"
	"go.uploadedlobster.com/musicbrainzws2"
)

func TestStartEndCommandContext(t *testing.T) {
	t.Cleanup(ResetForTest)

	StartCommandContext(time.Second)
	ctx := Context()

	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected command context deadline")
	}
	if time.Until(deadline) <= 0 {
		t.Fatalf("deadline = %v, want future time", deadline)
	}

	EndCommandContext()

	select {
	case <-ctx.Done():
	default:
		t.Fatal("expected command context to be cancelled")
	}
	if Context() != context.Background() {
		t.Fatal("expected Context() to return background after EndCommandContext")
	}
}

func TestEndCommandContextIsIdempotent(t *testing.T) {
	t.Cleanup(ResetForTest)

	StartCommandContext(time.Second)
	EndCommandContext()
	EndCommandContext()
}

func TestCloseClientClearsReference(t *testing.T) {
	t.Cleanup(ResetForTest)

	closed := false
	Client = &stubClient{onClose: func() { closed = true }}

	if err := CloseClient(); err != nil {
		t.Fatalf("CloseClient() error = %v", err)
	}
	if Client != nil {
		t.Fatal("expected Client to be nil after CloseClient")
	}
	if !closed {
		t.Fatal("expected underlying client Close to run")
	}
}

type stubClient struct {
	onClose func()
}

func (s *stubClient) SearchArtists(context.Context, string, int, int) (musicbrainzws2.SearchArtistsResult, error) {
	panic("not implemented")
}

func (s *stubClient) SearchReleases(context.Context, string, int, int) (musicbrainzws2.SearchReleasesResult, error) {
	panic("not implemented")
}

func (s *stubClient) SearchReleaseGroups(context.Context, string, int, int) (musicbrainzws2.SearchReleaseGroupsResult, error) {
	panic("not implemented")
}

func (s *stubClient) LookupArtist(context.Context, mbtypes.MBID, []string) (musicbrainzws2.Artist, error) {
	panic("not implemented")
}

func (s *stubClient) LookupRelease(context.Context, mbtypes.MBID, []string) (musicbrainzws2.Release, error) {
	panic("not implemented")
}

func (s *stubClient) LookupReleaseGroup(context.Context, mbtypes.MBID, []string) (musicbrainzws2.ReleaseGroup, error) {
	panic("not implemented")
}

func (s *stubClient) Close() error {
	if s.onClose != nil {
		s.onClose()
	}
	return nil
}
