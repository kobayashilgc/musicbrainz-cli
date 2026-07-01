package pagination

import "testing"

func TestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		limit   int
		offset  int
		wantErr bool
	}{
		{name: "defaults", limit: DefaultLimit, offset: DefaultOffset},
		{name: "min limit", limit: MinLimit, offset: 0},
		{name: "max limit", limit: MaxLimit, offset: 0},
		{name: "limit too low", limit: 0, offset: 0, wantErr: true},
		{name: "limit too high", limit: 101, offset: 0, wantErr: true},
		{name: "negative offset", limit: 25, offset: -1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := Validate(tt.limit, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate(%d, %d) error = %v, wantErr %v", tt.limit, tt.offset, err, tt.wantErr)
			}
		})
	}
}

func TestNewPaginator(t *testing.T) {
	t.Parallel()

	p := NewPaginator(10, 20)
	if p.Limit != 10 || p.Offset != 20 {
		t.Fatalf("got limit=%d offset=%d, want 10/20", p.Limit, p.Offset)
	}
}
