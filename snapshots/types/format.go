package types

// CurrentFormat is the currently used format for snapshots. Snapshots using the same format
// must be identical across all nodes for a given height, so this must be bumped when the binary
// snapshot output changes.
// Format 1 was the original format.
// Format 2 serailizes app version in addition to the original 
// snapshot data.
const CurrentFormat uint32 = 2
