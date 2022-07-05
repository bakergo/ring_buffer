package rings

type RingBuffer[T any] struct {
	noCopy noCopy
	head   int
	len    int
	buf    []T
}

// Creates a new RingBuffer with a capacity of `sz` items.
func New[T any](sz int) *RingBuffer[T] {
	return &RingBuffer[T]{
		head: 0,
		len:  0,
		buf:  make([]T, sz),
	}
}

// Returns the utilized capacity of the buffer.
func (r *RingBuffer[T]) Len() int {
	return r.len
}

// Returns the maximum capacity of the buffer.
func (r *RingBuffer[T]) Cap() int {
	return cap(r.buf)
}

// Inserts `d` items after the last item in the buffer.
// If the ring buffer does not have sufficient storage for the items, panic.
func (r *RingBuffer[T]) Append(d ...T) {
	if r.len+len(d) > cap(r.buf) {
		panic("insert over capacity")
	}

	// Copy to the end
	next := (r.head + r.len) % cap(r.buf)
	// Number of items we can copy
	items := cap(r.buf) - next
	if items >= len(d) {
		copy(r.buf[next:], d)
	} else {
		copy(r.buf[next:], d[:items])
		copy(r.buf, d[items:])
	}
	r.len += len(d)
}

// Inserts items before the first item in the buffer.
// If the buffer does not have sufficient storage for the items, panic.
func (r *RingBuffer[T]) Prepend(d ...T) {
	if r.len+len(d) > cap(r.buf) {
		panic("insert over capacity")
	}
	head := (r.head + cap(r.buf) - len(d)) % cap(r.buf)
	items := cap(r.buf) - head
	if items >= len(d) {
		copy(r.buf[head:], d)
	} else {
		copy(r.buf[head:], d[:items])
		copy(r.buf, d[items:])
	}
	r.head = head
	r.len += len(d)
}

// Removes and returns `num` items with the highest indices.
// If `num` is higher than the number of items, truncates and returns.
func (r *RingBuffer[T]) PopLast(num int) []T {
	if num > r.len {
		num = r.len
	}
	b := make([]T, num)
	if num <= 0 {
		return b
	}
	first := r.idx(r.len - num)
	last := r.idx(r.len)

	if first < last {
		copy(b, r.buf[first:last])
	} else if first >= last {
		stride := cap(r.buf) - first
		copy(b, r.buf[first:])
		copy(b[stride:], r.buf[:last])
	}
	r.len -= num
	return b
}

// Removes and returns the first `num` items from the buffer.
// If `num` is greater than the number of items, truncates completely.
func (r *RingBuffer[T]) PopFirst(num int) []T {
	if num > r.len {
		num = r.len
	}
	b := make([]T, num)
	if num <= 0 {
		return b
	}

	first := r.idx(0)
	last := r.idx(num)
	if first < last {
		copy(b, r.buf[first:last])
	} else if first >= last {
		stride := cap(r.buf) - first
		copy(b, r.buf[first:])
		copy(b[stride:], r.buf[:last])
	}
	r.head = last
	r.len -= num
	return b
}

// Returns an item at the specified index `i`
func (r *RingBuffer[T]) Get(i int) T {
	if i >= r.len || i < 0 {
		panic("index out of bounds")
	}
	return r.buf[r.idx(i)]
}

// Replaces an item at the index `i` with the value `v`.
// If i == len, grows the RingBuffer
func (r *RingBuffer[T]) Set(i int, v T) {
	if i > r.len || i < 0 || i >= cap(r.buf) {
		panic("index out of bounds")
	}
	r.buf[r.idx(i)] = v
	if i == r.len {
		r.len++
	}
}

func (r *RingBuffer[T]) idx(i int) int {
	return (r.head + cap(r.buf) + i) % cap(r.buf)
}

// Inserts an item before the specified index.
// If an item is already present at `i`, performs an in-order shift.
func (r *RingBuffer[T]) Insert(i int, v T) {
	if i > r.len || i < 0 {
		panic("index out of bounds")
	}
	idx := r.idx(i)
	if i == r.len {
		r.buf[idx] = v
		r.len++
		return
	}

	// Try to minimize copies, and try to creep to the left if possible.
	// There are several segments to consider
	// If idx > r.head
	//   If r.head > 0
	//     Shift r.head:idx to r.head-1:idx-1
	//   If r.head == 0
	//     Shift idx:end to idx+1:end+1
	// If idx < r.head
	//   Shift idx:end to idx+1:end+1
	// Insert at idx
	last := r.idx(r.len)
	if idx > r.head {
		if r.head > 0 {
			copy(r.buf[r.head-1:idx], r.buf[r.head:idx+1])
			r.head--
			idx--
		} else {
			// r.head == 0
			copy(r.buf[idx+1:], r.buf[idx:])
		}
	} else if idx < r.head {
		// We're split over 0 -- Shift out.
		copy(r.buf[idx+1:last+1], r.buf[idx:last])
	} else {
		// idx == r.head
		r.head = r.idx(-1)
		idx = r.head
	}
	r.buf[idx] = v
	r.len++
}

// Removes an item at the specified index.
// Other items will be shifted so the collection remains in order.
func (r *RingBuffer[T]) Remove(i int) T {
	if i >= r.len || i < 0 {
		panic("index out of bounds")
	}
	idx := r.idx(i)
	last := r.idx(r.len)

	v := r.buf[idx]

	// Try to minimize copies, and try to creep to the left if possible.
	// If idx < end
	//   Shift idx+1:end+1 to idx:end
	// If idx > end
	//   Shift head:idx-1 to head+1:idx

	if idx < last {
		copy(r.buf[idx:last-1], r.buf[idx+1:last])
	} else if idx > last {
		copy(r.buf[r.head+1:idx+1], r.buf[r.head:idx])
		r.head++
	}
	r.len--
	return v
}

// Appends the data from r2 to the receiver.
func (r *RingBuffer[T]) CopyFrom(r2 *RingBuffer[T]) {
	if r.Len()+r2.Len() > r.Cap() {
		panic("cannot concatenate arrays (not enough cap)")
	}

	r2end := r2.head + r2.len

	if r2end > cap(r2.buf) {
		r.Append(r2.buf[r2.head:]...)
		r.Append(r2.buf[:(r2.head+r2.len)%cap(r2.buf)]...)
	} else {
		r.Append(r2.buf[r2.head : r2.head+r2.len]...)
	}
}

// Removes the last `num` data in the ring buffer.
func (r *RingBuffer[T]) TruncLast(num int) {
	if num > r.len {
		num = r.len
	}
	if num <= 0 {
		return
	}
	r.len -= num
}

// Removes the first `num` data in the ring buffer.
func (r *RingBuffer[T]) TruncFirst(num int) {
	if num > r.len {
		num = r.len
	}
	if num <= 0 {
		return
	}

	last := r.idx(num)
	r.head = last
	r.len -= num
}

// Returns a shallow copy of the data contained in the RingBuffer as a slice.
func (r *RingBuffer[T]) AsSlice() []T {
	b := make([]T, r.len)
	first := r.idx(0)
	last := r.idx(r.len)

	if r.len == 0 {
		return b
	}
	if first < last {
		copy(b, r.buf[first:])
		return b
	}
	stride := cap(r.buf) - r.head
	copy(b[:stride], r.buf[first:])
	copy(b[stride:], r.buf[:last])
	return b
}

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
