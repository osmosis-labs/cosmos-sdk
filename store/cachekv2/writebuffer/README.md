# Write Buffer

This module implements a simple write buffer, satisfying the current cosmos iteration API.
We make a btree map of KV pairs that we write to.
All reads are reads off the write buffer, this has no notion of parents.

There are no mutexes here, its up to a caller to set concurrency guidelines.
The API for reads and writes here, use `key` represented as a string.

## Restrictions

No new store writes are allowed while any iterator is open.
(TODO: We need to relax this. TBD of do we this via data structure change, changing our iterator to know if a write happened and re-seek, etc. Common use cases to allow: Updating key being iterated on, e.g. `vec.itermut()`, or deleting a range of keys. (Dumb that we use this API for it, maybe we change that))

TODO for future applications, investigate making tidwall/btree preserve allocated RAM on `Clear()`.