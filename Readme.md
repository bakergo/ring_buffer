# Ring Buffer

This package implements a non-resizing, unsynchronized Ring Buffer. Inserts at either end are O(1), but in the middle are O(n) as the buffer needs to be shifted.

## Why?

There's few implementations of a ring buffer yet readily available on Github and taking advantage of go 1.18's generics. Some of these Slice operations are a little tricky and I wanted to ensure that I had tests for them.

## Installation

`go get github.com/bakergo/ring_buffer`

## Example Usage

```
func demo() {

  // Defines a new int RingBuffer of capacity 10
  rb := ring_buffer.New[int](10)

  // Append/Prepend methods exist
  rb.Append(1, 2, 3)
  rb.Prepend(4, 5, 6) // 4,5,6,1,2,3

  // Return an item at index
  x:= rb.Get(3) // x==1
  rb.Set(4, 3) // 4,5,6,1,4,3

  // Insert and remove
  rb.Remove(3) // 4,5,6,4,3
  rb.Insert(2, 1) // 4,5,1,6,4,3
}
```

