# State model

The existing SDK state model is broken with regard to whats the role of iterators.
We propose the mental model to be that every "Store" call be a claim on a collection (range) of data.
There can be multiple concurrent "read" claims on a range at a time.
There may only be one open "write" claim on a range, and if there are any open write claims then there can be no other read claims.

## Motivation

Right now the mental model of state is that it is one (massive) Key-Value map per module.
When you get a store, you are claiming read-write access over it. (either the whole module, or some range in the case of a prefix store)

We currently buffer all writes to this module in a "CacheKV store".
But our model for dealing with these writes is quite unsafe & unspecified.
You can have many readers and writers to the same data in parallel.
Each iterator-reader reads the state of the data at time of reader-open, despite whatever its call trace may do.
You may also want to write somewhere outside of this range in the store.
This should in principle be fine, however it is not right now.
