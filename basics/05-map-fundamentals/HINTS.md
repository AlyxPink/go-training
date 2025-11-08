# Hints for Map Fundamentals

## Level 1: Getting Started

- Create maps with `make(map[K]V)` or `map[K]V{k1: v1, k2: v2}`
- Add/update: `m[key] = value`
- Delete: `delete(m, key)`
- Check existence: `val, ok := m[key]`

## Level 2: Word Frequency

- Initialize an empty map: `make(map[string]int)`
- Range over the input slice
- Increment count for each word: `m[word]++`
- Zero value for int is 0, so this works even for new keys

## Level 3: Map Inversion

- Create a new map with swapped key/value types
- Range over original map: `for k, v := range m`
- Set inverted[v] = k
- Consider: what if multiple keys have the same value?

## Level 4: Map Merging

- Create a new map for the result
- Copy all entries from first map
- Copy all entries from second map (overwrites duplicates)
- Range over each map and add entries

## Level 5: Map Gotchas

Important things to know:
- Nil map reads return zero value but writes panic
- Map iteration order is random (non-deterministic)
- Maps are not safe for concurrent use without synchronization
- You cannot take the address of a map element: `&m[key]` is invalid
