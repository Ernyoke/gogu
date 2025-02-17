<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# trie

```go
import "github.com/esimov/gogu/trie"
```

Package trie provides a concurrent safe implementation of the ternary search tree data structure. Trie is similar to binary search tree, but it has up to three children rather than two as of BST. Tries are used for locating specific keys from within a set or for quick lookup searches within a text like auto\-completion or spell checking.

<details><summary>Example</summary>
<p>

```go
{
	q := queue.New[string]()
	trie := New[string, int](q)
	input := []string{"cats", "cape", "captain", "foes",
		"apple", "she", "root", "shells", "the", "thermos", "foo"}

	for idx, v := range input {
		trie.Put(v, idx)
	}

	longestPref, _ := trie.LongestPrefix("capetown")
	q1, _ := trie.StartsWith("ca")

	result := []string{}
	for q1.Size() > 0 {
		val, _ := q1.Dequeue()
		result = append(result, val)
	}

	fmt.Println(trie.Size())
	fmt.Println(longestPref)
	fmt.Println(result)

}
```

#### Output

```
11
cape
[cape captain cats]
```

</p>
</details>

## Index

- [Variables](<#variables>)
- [type Item](<#type-item>)
- [type Queuer](<#type-queuer>)
- [type Trie](<#type-trie>)
  - [func New[K ~string, V any](q Queuer[K]) *Trie[K, V]](<#func-new>)
  - [func (t *Trie[K, V]) Contains(key K) bool](<#func-triek-v-contains>)
  - [func (t *Trie[K, V]) Get(key K) (v V, ok bool)](<#func-triek-v-get>)
  - [func (t *Trie[K, V]) Keys() (Queuer[K], error)](<#func-triek-v-keys>)
  - [func (t *Trie[K, V]) LongestPrefix(query K) (K, error)](<#func-triek-v-longestprefix>)
  - [func (t *Trie[K, V]) Put(key K, val V)](<#func-triek-v-put>)
  - [func (t *Trie[K, V]) Size() int](<#func-triek-v-size>)
  - [func (t *Trie[K, V]) StartsWith(prefix K) (Queuer[K], error)](<#func-triek-v-startswith>)


## Variables

```go
var ErrorNotFound = fmt.Errorf("trie node not found")
```

## type [Item](<https://github.com/esimov/gogu/blob/master/trie/trie.go#L36-L39>)

Item is a key\-value struct pair used for storing the node values.

```go
type Item[K ~string, V any] struct {
    // contains filtered or unexported fields
}
```

## type [Queuer](<https://github.com/esimov/gogu/blob/master/trie/trie.go#L19-L24>)

Queuer exposes the basic interface methods for querying the trie data structure both for searching and for retrieving the existing keys. These are generic methods having the same signature as the corresponding concrete methods from the queue package. Because both the plain array and the linked listed version of the queue package has the same method signature, each of them could be plugged in on the method invocation.

```go
type Queuer[K ~string] interface {
    Enqueue(K)
    Dequeue() (K, error)
    Size() int
    Clear()
}
```

## type [Trie](<https://github.com/esimov/gogu/blob/master/trie/trie.go#L53-L58>)

Trie is a lock\-free tree data structure having the root as the first node. It's guarded with a mutex for concurrent\-safe data access.

```go
type Trie[K ~string, V any] struct {
    // contains filtered or unexported fields
}
```

### func [New](<https://github.com/esimov/gogu/blob/master/trie/trie.go#L61>)

```go
func New[K ~string, V any](q Queuer[K]) *Trie[K, V]
```

New initializes a new Trie data structure.

### func \(\*Trie\[K, V\]\) [Contains](<https://github.com/esimov/gogu/blob/master/trie/trie.go#L77>)

```go
func (t *Trie[K, V]) Contains(key K) bool
```

Contains checks if a key exists in the symbol table.

### func \(\*Trie\[K, V\]\) [Get](<https://github.com/esimov/gogu/blob/master/trie/trie.go#L123>)

```go
func (t *Trie[K, V]) Get(key K) (v V, ok bool)
```

Get retrieves a node's value based on the key. If the key does not exist it returns false.

### func \(\*Trie\[K, V\]\) [Keys](<https://github.com/esimov/gogu/blob/master/trie/trie.go#L211>)

```go
func (t *Trie[K, V]) Keys() (Queuer[K], error)
```

Keys collects all the existing keys in the set.

### func \(\*Trie\[K, V\]\) [LongestPrefix](<https://github.com/esimov/gogu/blob/master/trie/trie.go#L158>)

```go
func (t *Trie[K, V]) LongestPrefix(query K) (K, error)
```

LongestPrefix returns the longest prefix of query in the symbol table or empty if such string does not exist.

### func \(\*Trie\[K, V\]\) [Put](<https://github.com/esimov/gogu/blob/master/trie/trie.go#L90>)

```go
func (t *Trie[K, V]) Put(key K, val V)
```

Put inserts a new node into the symbol table, overwriting the old value with the new one if the key is already in the symbol table.

### func \(\*Trie\[K, V\]\) [Size](<https://github.com/esimov/gogu/blob/master/trie/trie.go#L69>)

```go
func (t *Trie[K, V]) Size() int
```

Size returns the trie size.

### func \(\*Trie\[K, V\]\) [StartsWith](<https://github.com/esimov/gogu/blob/master/trie/trie.go#L188>)

```go
func (t *Trie[K, V]) StartsWith(prefix K) (Queuer[K], error)
```

StartsWith returns all the keys in the set that start with prefix.



