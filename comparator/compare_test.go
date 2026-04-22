package comparator

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
