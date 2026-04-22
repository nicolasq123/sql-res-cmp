package cmp

import "testing"

func TestCompare_Identical(t *testing.T) {
	c := NewComparator()
	r1 := [][]string{
		{"1", "a"},
		{"2", "b"},
	}
	r2 := [][]string{
		{"2", "b"},
		{"1", "a"},
	}
	cols := []string{"id", "name"}
	diff := c.Compare(r1, r2, cols)
	if !diff.IsEmpty() {
		t.Errorf("expected empty diff, got %v", diff)
	}
}

func TestCompare_DifferentRows(t *testing.T) {
	c := NewComparator()
	r1 := [][]string{{"1", "a"}}
	r2 := [][]string{{"1", "a"}, {"2", "b"}}
	cols := []string{"id", "name"}
	diff := c.Compare(r1, r2, cols)
	if diff.IsEmpty() {
		t.Error("expected non-empty diff")
	}
	if len(diff.LeftOnly) != 1 || len(diff.RightOnly) != 2 {
		t.Errorf("unexpected diff: LeftOnly=%d, RightOnly=%d", len(diff.LeftOnly), len(diff.RightOnly))
	}
}

func TestCompare_Modified(t *testing.T) {
	c := NewComparator()
	r1 := [][]string{{"1", "a"}}
	r2 := [][]string{{"1", "b"}}
	cols := []string{"id", "name"}
	diff := c.Compare(r1, r2, cols)
	if diff.IsEmpty() {
		t.Error("expected non-empty diff")
	}
	if len(diff.Modified) != 1 {
		t.Errorf("expected 1 modified row, got %d", len(diff.Modified))
	}
}

func TestCompareByKey(t *testing.T) {
	c := NewComparator()
	r1 := [][]string{
		{"1", "a"},
		{"2", "b"},
	}
	r2 := [][]string{
		{"2", "x"},
		{"3", "c"},
	}
	cols := []string{"id", "name"}
	diff := c.CompareByKey(r1, r2, []string{"id"}, cols)
	if len(diff.LeftOnly) != 1 {
		t.Errorf("expected 1 left-only row, got %d", len(diff.LeftOnly))
	}
	if len(diff.RightOnly) != 1 {
		t.Errorf("expected 1 right-only row, got %d", len(diff.RightOnly))
	}
	if len(diff.Modified) != 1 {
		t.Errorf("expected 1 modified row, got %d", len(diff.Modified))
	}
}

func TestDiff_String(t *testing.T) {
	d := &Diff{
		LeftOnly:  [][]string{{"1", "a"}},
		RightOnly: [][]string{{"2", "b"}},
		Modified:  []RowDiff{{Left: []string{"3", "x"}, Right: []string{"3", "y"}}},
		Columns:   []string{"id", "name"},
	}
	s := d.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}

func TestCompare_Empty(t *testing.T) {
	c := NewComparator()
	var r1, r2 [][]string
	cols := []string{"id", "name"}
	diff := c.Compare(r1, r2, cols)
	if !diff.IsEmpty() {
		t.Error("expected empty diff for empty inputs")
	}
}

func TestCompare_SingleRow(t *testing.T) {
	c := NewComparator()
	r1 := [][]string{{"1", "a"}}
	r2 := [][]string{{"1", "a"}}
	cols := []string{"id", "name"}
	diff := c.Compare(r1, r2, cols)
	if !diff.IsEmpty() {
		t.Error("expected empty diff")
	}
}

func TestCompareByKey_Empty(t *testing.T) {
	c := NewComparator()
	var r1, r2 [][]string
	cols := []string{"id", "name"}
	diff := c.CompareByKey(r1, r2, []string{"id"}, cols)
	if !diff.IsEmpty() {
		t.Error("expected empty diff")
	}
}

func TestCompareByKey_SameKey(t *testing.T) {
	c := NewComparator()
	r1 := [][]string{{"1", "a"}, {"1", "b"}}
	r2 := [][]string{{"1", "a"}, {"1", "b"}}
	cols := []string{"id", "name"}
	diff := c.CompareByKey(r1, r2, []string{"id"}, cols)
	if !diff.IsEmpty() {
		t.Error("expected empty diff")
	}
}

func TestDiff_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		diff *Diff
		want bool
	}{
		{
			name: "all empty",
			diff: &Diff{LeftOnly: nil, RightOnly: nil, Modified: nil},
			want: true,
		},
		{
			name: "has LeftOnly",
			diff: &Diff{LeftOnly: [][]string{{"1"}}},
			want: false,
		},
		{
			name: "has RightOnly",
			diff: &Diff{RightOnly: [][]string{{"1"}}},
			want: false,
		},
		{
			name: "has Modified",
			diff: &Diff{Modified: []RowDiff{{}}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.diff.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRowsEqual(t *testing.T) {
	tests := []struct {
		name string
		r1   []string
		r2   []string
		want bool
	}{
		{"same", []string{"a", "b"}, []string{"a", "b"}, true},
		{"different length", []string{"a"}, []string{"a", "b"}, false},
		{"different value", []string{"a", "b"}, []string{"a", "c"}, false},
		{"empty", []string{}, []string{}, true},
		{"nil vs empty", nil, []string{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rowsEqual(tt.r1, tt.r2); got != tt.want {
				t.Errorf("rowsEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompareRows(t *testing.T) {
	tests := []struct {
		name string
		r1   []string
		r2   []string
		want int
	}{
		{"same", []string{"a", "b"}, []string{"a", "b"}, 0},
		{"less", []string{"a"}, []string{"b"}, -1},
		{"greater", []string{"b"}, []string{"a"}, 1},
		{"shorter", []string{"a"}, []string{"a", "b"}, -1},
		{"longer", []string{"a", "b"}, []string{"a"}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareRows(tt.r1, tt.r2); got != tt.want {
				t.Errorf("compareRows() = %v, want %v", got, tt.want)
			}
		})
	}
}
