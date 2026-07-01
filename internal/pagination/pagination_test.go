package pagination

import "testing"

func TestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		limit   int
		pageNo  int
		wantErr bool
	}{
		{name: "defaults", limit: DefaultLimit, pageNo: DefaultPageNo},
		{name: "min limit", limit: MinLimit, pageNo: 1},
		{name: "max limit", limit: MaxLimit, pageNo: 1},
		{name: "limit too low", limit: 0, pageNo: 1, wantErr: true},
		{name: "limit too high", limit: 101, pageNo: 1, wantErr: true},
		{name: "pageno too low", limit: 25, pageNo: 0, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := Validate(tt.limit, tt.pageNo)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate(%d, %d) error = %v, wantErr %v", tt.limit, tt.pageNo, err, tt.wantErr)
			}
		})
	}
}

func TestOffset(t *testing.T) {
	t.Parallel()

	if got := Offset(10, 1); got != 0 {
		t.Fatalf("Offset(10, 1) = %d, want 0", got)
	}
	if got := Offset(10, 2); got != 10 {
		t.Fatalf("Offset(10, 2) = %d, want 10", got)
	}
}

func TestNewPaginator(t *testing.T) {
	t.Parallel()

	p := NewPaginator(10, 3)
	if p.Limit != 10 || p.Offset != 20 {
		t.Fatalf("got limit=%d offset=%d, want 10/20", p.Limit, p.Offset)
	}
}
