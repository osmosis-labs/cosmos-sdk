# Secp256k1 keys

This folder creates an implementation of Secp256k1 keys.

Pubkey format: Secp256k1 elliptic curve point, serialized in either compressed 33 byte form, or uncompressed 65 byte form.

Address format: `RipeMD160(SHA-2(compressed pubkey bytes))` -- TODO: Its not at all clear to me that this is implemented correctly if an uncompressed pubkey is supplied..., working on test vectors for this.

Signature formats:

* `(R|S)` ECDSA signatures, with SHA-2 message hashing
* `(R|S|V)` ECDSA signatures, with Keccak message hashing, and `V` being the EC Recover byte in EVM-compatible form (oneof 27, 28)

Security considerations:

- Signature forgery requires collision attacks on SHA-2 or Keccak
- No malleability possible on the ECRecover byte