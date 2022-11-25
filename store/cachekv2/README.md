# CacheKV

A `CacheKVStore` is cache wrapper for a `KVStore`. It extends the operations of the `KVStore` to work with a cache, allowing for reduced I/O operations and more efficient disposing of changes (e.g. after processing a failed transaction).

The core goals the CacheKVStore seeks to solve are:

* Buffer all writes to the parent store, so they can be dropped if they needs to be reverted
* Allow iteration over contiguous spans of keys (all SDK stores)
* Act as a cache, so we don't repeat I/O to disk for reads we've already done
* Make subsequent reads, account for prior buffered writes
* Write all changes to the parent store

We build this as two stores, one for "dirty write storage", and a second for read caching.
We then make a default constructor that creates both.

This is NOT compatible with CacheKV store v1, or the existing store semantics.
This is because the current design of the SDK allows for (and actively abuses) iterator invalidation,
a primary source of single threaded concurrency issues.
(This is in fact the straw that broke the camel's back triggering this rewrite)

Issues are that namely during iterator creation you can:

* Always mutate what your iterating over (no iter_mut vs iter separation)
* Delete your current entry during an iterator call (because theres no DeleteRange API)
* New writes & deletes (even into the range your iterating over), that just don't get added into your iterator
    * The way this is designed, can also lead to deadlocks / result incorrectness if multiple iterators exist in parallel

## Intended later changes

More API breaks are needed for these, but would like:

* Express separate cache retention policies for:
    * One-off read caching (e.g. only LRU cache for 1000 entries)
    * Iterator interval read caching
* Charging a separate gas amount for a read in cache vs not in cache.

An SDK design choice that should happen is:

* BeginBlock and EndBlock should not have the write buffer enabled, and should use a different cache for minimal LRU's on reads