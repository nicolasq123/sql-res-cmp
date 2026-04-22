package comparator

import (
	"fmt"
	"sort"
	"strings"
)

type Diff struct {
	LeftOnly  [][]string
	RightOnly [][]string
	Modified  []RowDiff
	Columns   []string
}

type RowDiff struct {
	Left  []string
	Right []string
}

type Comparator struct{}

func NewComparator() *Comparator {
	return &Comparator{}
}

func (c *Comparator) Compare(r1, r2 [][]string, cols []string) *Diff {
	diff := &Diff{Columns: cols}

	if len(r1) != len(r2) {
		diff.LeftOnly = r1
		diff.RightOnly = r2
		return diff
	}

	s1 := make([][]string, len(r1))
	copy(s1, r1)
	s2 := make([][]string, len(r2))
	copy(s2, r2)

	sort.Slice(s1, func(i, j int) bool {
		return compareRows(s1[i], s1[j]) < 0
	})
	sort.Slice(s2, func(i, j int) bool {
		return compareRows(s2[i], s2[j]) < 0
	})

	for i := range s1 {
		if !rowsEqual(s1[i], s2[i]) {
			diff.Modified = append(diff.Modified, RowDiff{Left: s1[i], Right: s2[i]})
		}
	}

	return diff
}

func (c *Comparator) CompareByKey(r1, r2 [][]string, keyCols []string, allCols []string) *Diff {
	diff := &Diff{Columns: allCols}

	keyIdx := getColIndices(keyCols, allCols)

	key := func(row []string) string {
		var parts []string
		for _, i := range keyIdx {
			if i < len(row) {
				parts = append(parts, row[i])
			}
		}
		return strings.Join(parts, "\x00")
	}

	grouped1 := make(map[string][][]string)
	for _, row := range r1 {
		grouped1[key(row)] = append(grouped1[key(row)], row)
	}

	grouped2 := make(map[string][][]string)
	for _, row := range r2 {
		grouped2[key(row)] = append(grouped2[key(row)], row)
	}

	for k, rows1 := range grouped1 {
		rows2, ok := grouped2[k]
		if !ok {
			diff.LeftOnly = append(diff.LeftOnly, rows1...)
			continue
		}
		if len(rows1) != len(rows2) {
			diff.LeftOnly = append(diff.LeftOnly, rows1...)
			diff.RightOnly = append(diff.RightOnly, rows2...)
			continue
		}
		sort.Slice(rows1, func(i, j int) bool {
			return compareRows(rows1[i], rows1[j]) < 0
		})
		sort.Slice(rows2, func(i, j int) bool {
			return compareRows(rows2[i], rows2[j]) < 0
		})
		for i := range rows1 {
			if !rowsEqual(rows1[i], rows2[i]) {
				diff.Modified = append(diff.Modified, RowDiff{Left: rows1[i], Right: rows2[i]})
			}
		}
		delete(grouped2, k)
	}
	for _, rows2 := range grouped2 {
		diff.RightOnly = append(diff.RightOnly, rows2...)
	}

	return diff
}

func getColIndices(cols []string, allCols []string) []int {
	colMap := make(map[string]int, len(allCols))
	for i, c := range allCols {
		colMap[c] = i
	}
	idx := make([]int, len(cols))
	for i, col := range cols {
		if j, ok := colMap[col]; ok {
			idx[i] = j
		}
	}
	return idx
}

func compareRows(r1, r2 []string) int {
	min := len(r1)
	if len(r2) < min {
		min = len(r2)
	}
	for i := 0; i < min; i++ {
		if r1[i] < r2[i] {
			return -1
		}
		if r1[i] > r2[i] {
			return 1
		}
	}
	return len(r1) - len(r2)
}

func rowsEqual(r1, r2 []string) bool {
	if len(r1) != len(r2) {
		return false
	}
	for i := range r1 {
		if r1[i] != r2[i] {
			return false
		}
	}
	return true
}

func (d *Diff) IsEmpty() bool {
	return len(d.LeftOnly) == 0 && len(d.RightOnly) == 0 && len(d.Modified) == 0
}

func (d *Diff) String() string {
	if d.IsEmpty() {
		return "PASS: 结果一致"
	}

	var sb strings.Builder
	sb.WriteString("FAIL: 结果不一致\n")

	if len(d.LeftOnly) > 0 {
		sb.WriteString(fmt.Sprintf("\n仅左侧有 (%d 行):\n", len(d.LeftOnly)))
		for _, row := range d.LeftOnly {
			sb.WriteString(fmt.Sprintf("  %s\n", strings.Join(row, ", ")))
		}
	}

	if len(d.RightOnly) > 0 {
		sb.WriteString(fmt.Sprintf("\n仅右侧有 (%d 行):\n", len(d.RightOnly)))
		for _, row := range d.RightOnly {
			sb.WriteString(fmt.Sprintf("  %s\n", strings.Join(row, ", ")))
		}
	}

	if len(d.Modified) > 0 {
		sb.WriteString(fmt.Sprintf("\n差异 (%d 行):\n", len(d.Modified)))
		for _, m := range d.Modified {
			sb.WriteString(fmt.Sprintf("  左侧: %s\n", strings.Join(m.Left, ", ")))
			sb.WriteString(fmt.Sprintf("  右侧: %s\n", strings.Join(m.Right, ", ")))
		}
	}

	return sb.String()
}
