# Write Buffer

This module implements a dead simple write buffer.
This only buffers writes, and reads only go to this buffer.
The goal of this is to only store delta's of whats written over.
There is no notion of a "parent", beyond a function to allow writing the write buffer into any other store.

There are no mutexes here, its up to a caller to set concurrency guidelines.
The API for reads and writes here, use `key` represented as a string.

This is NOT compatible with CacheKV store v1, or the existing store semantics.
This is because the current design of the SDK allows iterator invalidation,
a primary source of single threaded concurrency issues.
(And in fact the straw that broke the camel's back triggering this rewrite)

Issues are that namely during iterator creation you can:

* Always mutate what your iterating over (no iter_mut vs iter separation)
* Delete your current entry during an iterator call (because theres no DeleteRange API)
* New writes & deletes (even into the range your iterating over), that just don't get added into your iterator
    * The way this is designed, can also lead to deadlocks / result incorrectness if multiple iterators exist in parallel

## Restrictions

No new store writes are allowed while any iterator is open.
(TODO: We need to relax this. TBD of do we this via data structure change, changing our iterator to know if a write happened and re-seek, etc. Common use cases to allow: Updating key being iterated on, e.g. `vec.itermut()`, or deleting a range of keys. (Dumb that we use this API for it, maybe we change that))

TODO for future applications, investigate making tidwall/btree preserve allocated RAM on `Clear()`.