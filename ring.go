/*
Package ring implements a concurrent bounded circular buffer.  When at
maximum capacity, it drops the oldest elements to make room for the new
ones.

All operations have constant worst-case time complexity.

Internally the ring uses a fixed size buffer allocated upon
construction, proportional to its capacity.
*/
package ring

import (
	"fmt"
	"sync"
)

// Ring is a concurrent bounded circular buffer.
type Ring struct {
	mutex sync.Mutex
	buf   []interface{} // elements storage
	len   int           // how many elements are stored in the ring
	head  int           // index of the next element to be extracted
}

// Returns a new ring with the given capacity.
func New(cap int) (*Ring, error) {
	if cap < 1 {
		return nil, fmt.Errorf("ring capacity must be > 0, got %d", cap)
	}

	return &Ring{
		buf: make([]interface{}, cap),
	}, nil
}

// returns the index where the next element will be inserted.
func (r *Ring) tail() int {
	return (r.head + r.len) % cap(r.buf)
}

// Insert adds a new element to the ring. If the ring is already at
// maximum capacity, the oldest element is dropped to make room for the
// new one.
func (r *Ring) Insert(v interface{}) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// if full, make room by droppin the oldest element
	if r.len == cap(r.buf) {
		_, _ = r.extract()
	}

	r.buf[r.tail()] = v
	r.len++
}

// Extract extracts and returns the oldest element in the ring.
func (r *Ring) Extract() (interface{}, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.extract()
}

func (r *Ring) extract() (interface{}, bool) {
	if r.len == 0 {
		return nil, false
	}

	result := r.buf[r.head]
	r.head = (r.head + 1) % cap(r.buf)
	r.len--

	return result, true
}

// Peek returns the oldest element in the ring.
func (r *Ring) Peek() (interface{}, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.len == 0 {
		return nil, false
	}

	return r.buf[r.head], true
}

// Len returns the amount of elements in the ring.
func (r *Ring) Len() int {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.len
}
