package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntervalTree_InsertAndSearch(t *testing.T) {
	tests := []struct {
		name              string
		intervalsToInsert []Interval
		searchInterval    Interval
		expectedOverlap   *Interval
		expectInsertion   bool
	}{
		{
			name: "Insert and find overlapping interval",
			intervalsToInsert: []Interval{
				{Start: []byte("10"), End: []byte("20")},
				{Start: []byte("15"), End: []byte("25")},
			},
			searchInterval:  Interval{Start: []byte("18"), End: []byte("22")},
			expectedOverlap: &Interval{Start: []byte("18"), End: []byte("22")},
			expectInsertion: true,
		},
		{
			name: "Insert and find no overlapping interval",
			intervalsToInsert: []Interval{
				{Start: []byte("10"), End: []byte("20")},
				{Start: []byte("30"), End: []byte("40")},
			},
			searchInterval:  Interval{Start: []byte("21"), End: []byte("29")},
			expectedOverlap: nil,
			expectInsertion: true,
		},
		{
			name:              "Attempt to insert beyond max depth",
			intervalsToInsert: []Interval{
				// Add intervals here such that it exceeds maxDepth
			},
			searchInterval:  Interval{Start: []byte("50"), End: []byte("60")},
			expectedOverlap: nil,
			expectInsertion: false, // Expect that the last interval(s) will not be inserted
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := IntervalTree{}

			for _, interval := range tt.intervalsToInsert {
				beforeInsertDepth := tree.depth()
				tree.Insert(interval)
				afterInsertDepth := tree.depth()

				if tt.expectInsertion {
					assert.LessOrEqual(t, afterInsertDepth, maxDepth, "Tree depth should not exceed maxDepth")
					assert.GreaterOrEqual(t, afterInsertDepth, beforeInsertDepth, "Tree depth should increase after insertion")
				} else {
					assert.Equal(t, beforeInsertDepth, afterInsertDepth, "Tree depth should remain the same when insertion is skipped")
				}
			}

			result := tree.OverlapSearch(tt.searchInterval)
			assert.Equal(t, tt.expectedOverlap, result, "Overlap search did not return the expected result")
		})
	}
}

func TestIntervalTree_InsertAndSearchWithUnionBehavior(t *testing.T) {
	tests := []struct {
		name              string
		intervalsToInsert []Interval
		searchInterval    Interval
		expectedOverlap   *Interval
	}{
		{
			name: "Insert overlapping intervals",
			intervalsToInsert: []Interval{
				{Start: []byte("10"), End: []byte("15")},
				{Start: []byte("14"), End: []byte("20")},
			},
			searchInterval:  Interval{Start: []byte("13"), End: []byte("16")},
			expectedOverlap: &Interval{Start: []byte("13"), End: []byte("16")},
		},
		{
			name: "Insert non-overlapping then overlapping interval",
			intervalsToInsert: []Interval{
				{Start: []byte("10"), End: []byte("20")},
				{Start: []byte("30"), End: []byte("40")},
				{Start: []byte("35"), End: []byte("45")},
			},
			searchInterval:  Interval{Start: []byte("32"), End: []byte("38")},
			expectedOverlap: &Interval{Start: []byte("32"), End: []byte("38")},
		},
		{
			name: "Insert overlapping intervals with unbounded end",
			intervalsToInsert: []Interval{
				{Start: []byte("50"), End: []byte{}},
				{Start: []byte("45"), End: []byte("55")},
			},
			searchInterval:  Interval{Start: []byte("48"), End: []byte("60")},
			expectedOverlap: &Interval{Start: []byte("48"), End: []byte("60")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := IntervalTree{}

			for _, interval := range tt.intervalsToInsert {
				tree.Insert(interval)
			}

			result := tree.OverlapSearch(tt.searchInterval)
			assert.Equal(t, tt.expectedOverlap, result, "Overlap search did not return the expected result")
		})
	}
}

// depth is a helper function to determine the current depth of the tree.
// This is necessary to test the maxDepth constraint.
// Implement this function based on your tree's structure.
func (t *IntervalTree) depth() int {
	return calculateDepth(t.Root, 1)
}

func calculateDepth(node *Node, currentDepth int) int {
	if node == nil {
		return currentDepth - 1
	}

	leftDepth := calculateDepth(node.Left, currentDepth+1)
	rightDepth := calculateDepth(node.Right, currentDepth+1)

	if leftDepth > rightDepth {
		return leftDepth
	}
	return rightDepth
}
