package internal

import (
	"bytes"
)

// Maximum depth of the interval tree. We refuse to insert if it causes a depth larger than this.
// Bounds any worst case situations. Realistically, this is safe at quite large parameterizations.
// (Safe as relative to the current expensive dirty items operations that are possible.)
const maxDepth = 32

// IntervalTree represents an interval tree, a data structure that allows
// storing intervals and efficiently querying which intervals overlap with a given interval or point.
// It primarily supports operations like insertion of intervals and searching for overlapping intervals.
//
// The tree maintains a maximum depth as defined by the maxDepth constant to prevent excessive depth
// which could lead to inefficiency.
type IntervalTree struct {
	Root *Node
}

// Interval represents a range of values defined by Start and End.
// Start is the beginning of the interval, and End is the end of the interval.
// If End is a blank slice ([]byte{}), it represents an "unbounded" range, extending infinitely.
type Interval struct {
	Start, End []byte
}

// Node represents a single node in the IntervalTree.
// It contains an Interval and the maximum End value (Max) in the subtree rooted at this node.
// It also includes pointers to the left and right child nodes.
type Node struct {
	Interval
	Max   []byte
	Left  *Node
	Right *Node
}

func (i Interval) Contains(point []byte) bool {
	if len(i.End) == 0 {
		return bytes.Compare(i.Start, point) <= 0
	}
	return bytes.Compare(i.Start, point) <= 0 && bytes.Compare(i.End, point) >= 0
}

// Insert adds a new interval to the tree. It inserts the interval in a position
// that maintains the tree's properties and adheres to the maximum depth constraint.
// If the insertion would increase the tree's depth beyond maxDepth, the interval is not inserted,
// and the tree remains unchanged. This method does not return an error in such cases.
//
// Arguments:
// interval: The interval to be inserted into the tree.
//
// Complexity: O(log n) on average, where n is the number of intervals in the tree,
// assuming the tree's depth does not exceed maxDepth.
func (t *IntervalTree) Insert(interval Interval) {
	t.Root = t.insert(t.Root, interval, 1) // Start with depth 1
}

func (t *IntervalTree) insert(node *Node, interval Interval, depth int) *Node {
	if depth > maxDepth {
		// If the depth exceeds maxDepth, do not insert and leave the tree unchanged
		return node
	}

	if node == nil {
		return &Node{Interval: interval, Max: interval.End}
	}

	// Check if the start or end of the new interval falls within the current node's interval
	if node.Interval.Contains(interval.Start) || node.Interval.Contains(interval.End) {
		// Update the current node's interval to the union of the current interval and the new interval
		node.Interval = union(node.Interval, interval)
	}

	// Decide whether to insert in the left or right subtree
	if bytes.Compare(interval.Start, node.Interval.Start) < 0 {
		node.Left = t.insert(node.Left, interval, depth+1)
	} else {
		node.Right = t.insert(node.Right, interval, depth+1)
	}

	// Update the Max value if needed
	if bytes.Compare(node.Max, interval.End) < 0 {
		node.Max = interval.End
	}

	return node
}

// union returns the union of two intervals.
func union(i1, i2 Interval) Interval {
	start := min(i1.Start, i2.Start)
	end := max(i1.End, i2.End)
	return Interval{Start: start, End: end}
}

// OverlapSearch finds an interval in the tree that overlaps with the given interval.
// It returns the first such interval found. If no overlapping interval is found,
// it returns nil. An interval with a blank End is treated as unbounded and
// overlaps with all intervals starting before or at its Start.
//
// Complexity: O(log n) on average, where n is the number of intervals in the tree.
func (t *IntervalTree) OverlapSearch(interval Interval) *Interval {
	node := overlapSearch(t.Root, interval)
	if node == nil {
		return node
	}
	return intersect(*node, interval)
}

// OverlapSearch searches for an interval in the tree that overlaps with the given interval.
// It utilizes the properties of the interval tree to efficiently find an overlapping interval, if one exists.
//
// Arguments:
// interval: The interval to search for overlaps with. This interval can have a bounded or unbounded end.
//
// Returns:
// *Interval: A pointer to the first Interval found that overlaps with the given interval.
// If no overlapping interval is found, nil is returned.
//
// How it works:
//
//  1. Starting at the root of the interval tree, the function traverses the tree.
//     The decision to move left or right at each node is based on the comparison of the interval to be searched
//     with the interval stored at the current node, as well as the max value of the current node's left child.
//
//  2. At each node, it first checks if the current node's interval overlaps with the given interval using the
//     doOverlap function. If an overlap is found, that interval is immediately returned.
//
// 3. If no overlap is found at the current node, the search proceeds to the left or right child:
//
//   - If the left child's max value is greater than or equal to the start of the searching interval,
//     there might be an overlapping interval in the left subtree, so the search moves to the left child.
//
//   - Otherwise, the search moves to the right child, as any potential overlaps can only exist in the right subtree.
//
//     4. The search continues recursively until either an overlap is found, or all potential nodes where an
//     overlap could exist have been checked (i.e., when a leaf node or a null pointer is reached).
//
// Note:
//   - The efficiency of this search lies in the fact that not all nodes in the tree are visited. The max value
//     of the left child at each node allows skipping entire subtrees that cannot contain an overlapping interval.
//   - If the interval to be searched has an unbounded end, the doOverlap function is designed to handle this case,
//     considering the unbounded interval as extending infinitely.
//
// Example Usage:
// tree := IntervalTree{}
// // Insert intervals into the tree...
// searchInterval := Interval{Start: []byte("10"), End: []byte("20")}
// overlap := tree.OverlapSearch(searchInterval)
//
//	if overlap != nil {
//	    fmt.Println("Overlap found:", overlap)
//	} else {
//
//	    fmt.Println("No overlap found")
//	}
func overlapSearch(node *Node, interval Interval) *Interval {
	if node == nil {
		return nil
	}

	if doOverlap(node.Interval, interval) {
		return &node.Interval
	}

	if node.Left != nil && (len(node.Left.Max) == 0 || bytes.Compare(node.Left.Max, interval.Start) >= 0) {
		return overlapSearch(node.Left, interval)
	}

	return overlapSearch(node.Right, interval)
}

func doOverlap(i1, i2 Interval) bool {
	// Check if either interval has an unbounded end.
	i1Unbounded := len(i1.End) == 0
	i2Unbounded := len(i2.End) == 0

	// If both intervals are unbounded, they overlap if i1 starts before or at the start of i2.
	if i1Unbounded && i2Unbounded {
		return bytes.Compare(i1.Start, i2.Start) <= 0
	}

	// If i1 is unbounded, they overlap if i1 starts before or at the end of i2.
	if i1Unbounded {
		return bytes.Compare(i1.Start, i2.End) <= 0
	}

	// If i2 is unbounded, they overlap if i2 starts before or at the end of i1.
	if i2Unbounded {
		return bytes.Compare(i2.Start, i1.End) <= 0
	}

	// If neither is unbounded, use the original comparison.
	return bytes.Compare(i1.Start, i2.End) <= 0 && bytes.Compare(i2.Start, i1.End) <= 0
}

func intersect(i1, i2 Interval) *Interval {
	start := max(i1.Start, i2.Start)
	end := min(i1.End, i2.End)

	return &Interval{Start: start, End: end}
}

func max(a, b []byte) []byte {
	if len(a) == 0 || len(b) == 0 {
		return []byte{}
	}
	if bytes.Compare(a, b) > 0 {
		return a
	}
	return b
}

func min(a, b []byte) []byte {
	if len(a) == 0 {
		return b
	}
	if len(b) == 0 {
		return a
	}
	if bytes.Compare(a, b) < 0 {
		return a
	}
	return b
}
