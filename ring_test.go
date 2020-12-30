package ring_test

import (
	"fmt"
	"testing"

	"github.com/alcortesm/ring"
)

func TestRing(t *testing.T) {
	t.Parallel()

	subtests := map[string]func(*testing.T){
		"invalid capacity":               ringInvalidCapacity,
		"new ring is empty":              ringNewIsEmpty,
		"single insert":                  ringSingleInsert,
		"few inserts":                    ringFewInserts,
		"alternate inserts and extracts": ringAlternateInsertsAndExtracts,
		"forgets oldest":                 ringForgetsOldest,
		"extracts OK after forgetting":   ringExtractsOKAfterForgetting,
		"all together":                   ringAllTogether,
	}

	for name, testFn := range subtests {
		testFn := testFn
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			testFn(t)
		})
	}
}

// asserts that the ring has the length we want.
func assertLen(t *testing.T, r *ring.Ring, want int) {
	t.Helper()

	got := r.Len()
	if got != want {
		t.Fatalf("wrong length, want %d, got %d", want, got)
	}
}

// asserts that a ring is empty.
func assertEmpty(t *testing.T, r *ring.Ring) {
	t.Helper()

	assertLen(t, r, 0)

	got, ok := r.Peek()
	if ok || got != nil {
		t.Fatalf("want empty ring, but peek returned %T, %t", got, ok)
	}

	got, ok = r.Extract()
	if ok || got != nil {
		t.Fatalf("want empty ring, but extract returned %T, %t", got, ok)
	}
}

// asserts that peeking at the ring returns the expected value.
func assertPeek(t *testing.T, r *ring.Ring, want int) {
	t.Helper()

	v, ok := r.Peek()
	if !ok {
		t.Fatalf("when peeking %d: unexpected empty ring", want)
	}

	got, ok := v.(int)
	if !ok {
		t.Fatalf("cannot type assert peeked value to int, got: %T (%#[1]v)", v)
	}

	if got != want {
		t.Errorf("unexpected peeked value, want %d, got %d", want, got)
	}
}

// asserts that the value v has been successfully extracted from the
// ring r.
func assertExtract(t *testing.T, r *ring.Ring, want int) {
	t.Helper()

	v, ok := r.Extract()
	if !ok {
		t.Fatalf("when extracting %d: unexpected empty ring", want)
	}

	got, ok := v.(int)
	if !ok {
		t.Fatalf("cannot type assert to extracted value to int, got: %T (%#[1]v)", v)
	}

	if got != want {
		t.Errorf("unexpected extracted value, want %d, got %d", want, got)
	}
}

// tests that capacities smaller than 1 are invalid
func ringInvalidCapacity(t *testing.T) {
	for _, cap := range []int{0, -1, -42} {
		cap := cap
		name := fmt.Sprintf("cap=%d", cap)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := ring.New(cap)
			if err == nil {
				t.Fatal("unexpected success")
			}
		})
	}
}

// tests that new rings start empty.
func ringNewIsEmpty(t *testing.T) {
	for _, cap := range []int{1, 2, 42} {
		cap := cap
		name := fmt.Sprintf("cap=%d", cap)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			r, err := ring.New(cap)
			if err != nil {
				t.Fatalf("creating ring: %v", err)
			}

			assertEmpty(t, r)
		})
	}
}

// tests that you can extract a value after inserting it.
func ringSingleInsert(t *testing.T) {
	for _, cap := range []int{1, 2, 42} {
		cap := cap
		name := fmt.Sprintf("cap=%d", cap)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			r, err := ring.New(cap)
			if err != nil {
				t.Fatalf("creating ring: %v", err)
			}

			r.Insert(42)
			assertLen(t, r, 1)

			assertPeek(t, r, 42)
			assertLen(t, r, 1)

			assertExtract(t, r, 42)
			assertEmpty(t, r)
		})
	}
}

// tests that you can extract a bunch of values after inserting them if
// there is enough capacity for them.  They will be extraced in the
// same order they were inserted.
func ringFewInserts(t *testing.T) {
	for _, cap := range []int{4, 5, 42} {
		cap := cap
		name := fmt.Sprintf("cap=%d", cap)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			r, err := ring.New(cap)
			if err != nil {
				t.Fatalf("creating ring: %v", err)
			}

			r.Insert(1)
			r.Insert(2)
			r.Insert(30)
			r.Insert(42)
			assertLen(t, r, 4)

			assertPeek(t, r, 1)
			assertLen(t, r, 4)
			assertExtract(t, r, 1)
			assertLen(t, r, 3)

			assertPeek(t, r, 2)
			assertLen(t, r, 3)
			assertExtract(t, r, 2)
			assertLen(t, r, 2)

			assertPeek(t, r, 30)
			assertLen(t, r, 2)
			assertExtract(t, r, 30)
			assertLen(t, r, 1)

			assertPeek(t, r, 42)
			assertLen(t, r, 1)
			assertExtract(t, r, 42)
			assertEmpty(t, r)
		})
	}
}

// tests that you can alternate inserts and extracts while keeping the
// stored elements below the capacity.
func ringAlternateInsertsAndExtracts(t *testing.T) {
	for _, cap := range []int{3, 4, 42} {
		cap := cap
		name := fmt.Sprintf("cap=%d", cap)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			r, err := ring.New(cap)
			if err != nil {
				t.Fatalf("creating ring: %v", err)
			}

			r.Insert(1)
			assertExtract(t, r, 1)
			assertEmpty(t, r)

			r.Insert(2)
			r.Insert(3)
			assertLen(t, r, 2)
			assertPeek(t, r, 2)
			assertExtract(t, r, 2)
			assertLen(t, r, 1)
			assertPeek(t, r, 3)
			assertExtract(t, r, 3)
			assertEmpty(t, r)

			r.Insert(4)
			r.Insert(5)
			r.Insert(6)
			assertLen(t, r, 3)
			assertPeek(t, r, 4)
			assertExtract(t, r, 4)
			assertLen(t, r, 2)
			assertPeek(t, r, 5)
			assertExtract(t, r, 5)
			assertLen(t, r, 1)
			assertPeek(t, r, 6)
			assertExtract(t, r, 6)
			assertEmpty(t, r)
		})
	}
}

// tests that the ring forgets the oldest values when its capacity is
// reached and we keep inserting values.
func ringForgetsOldest(t *testing.T) {
	r, err := ring.New(4)
	if err != nil {
		t.Fatalf("creating ring: %v", err)
	}

	r.Insert(1)
	r.Insert(2)
	r.Insert(3)
	r.Insert(4)

	r.Insert(5)
	assertLen(t, r, 4)
	assertPeek(t, r, 2)

	r.Insert(6)
	r.Insert(7)
	assertLen(t, r, 4)
	assertPeek(t, r, 4)
}

// tests that you can extracts values fine when the ring has drop some
// values due to reaching its maximum capacity.
func ringExtractsOKAfterForgetting(t *testing.T) {
	r, err := ring.New(4)
	if err != nil {
		t.Fatalf("creating ring: %v", err)
	}

	r.Insert(1)
	r.Insert(2)
	r.Insert(3)
	r.Insert(4)

	r.Insert(5) // drops 1
	assertExtract(t, r, 2)
	assertLen(t, r, 3)
	assertPeek(t, r, 3)

	r.Insert(6)
	r.Insert(7) // drops 3
	r.Insert(8) // drops 4
	assertLen(t, r, 4)
	assertPeek(t, r, 5)
	assertExtract(t, r, 5)
	assertExtract(t, r, 6)
	assertExtract(t, r, 7)
	assertExtract(t, r, 8)
}

// tests all cases together.
func ringAllTogether(t *testing.T) {
	r, err := ring.New(4)
	if err != nil {
		t.Fatalf("creating ring: %v", err)
	}

	r.Insert(1)             // [1]
	assertExtract(t, r, 1)  // []
	assertEmpty(t, r)       //
	r.Insert(2)             // [2]
	r.Insert(3)             // [2, 3]
	assertExtract(t, r, 2)  // [3]
	assertLen(t, r, 1)      //
	r.Insert(4)             // [3, 4]
	assertLen(t, r, 2)      //
	assertPeek(t, r, 3)     //
	assertExtract(t, r, 3)  // [4]
	r.Insert(5)             // [4, 5]
	r.Insert(6)             // [4, 5, 6]
	assertLen(t, r, 3)      //
	assertPeek(t, r, 4)     //
	r.Insert(7)             // [4, 5, 6, 7]
	r.Insert(8)             // [5, 6, 7, 8]
	assertLen(t, r, 4)      //
	assertPeek(t, r, 5)     //
	r.Insert(9)             // [6, 7, 8, 9]
	assertLen(t, r, 4)      //
	assertPeek(t, r, 6)     //
	r.Insert(10)            // [7, 8, 9, 10]
	r.Insert(11)            // [8, 9, 10, 11]
	r.Insert(12)            // [9, 10, 11, 12]
	assertLen(t, r, 4)      //
	r.Insert(13)            // [10, 11, 12, 13]
	r.Insert(14)            // [11, 12, 13, 14]
	r.Insert(15)            // [12, 13, 14, 15]
	assertLen(t, r, 4)      //
	r.Insert(16)            // [13, 14, 15, 16]
	r.Insert(17)            // [14, 15, 16, 17]
	assertLen(t, r, 4)      //
	assertPeek(t, r, 14)    //
	assertExtract(t, r, 14) // [15, 16, 17]
	r.Insert(18)            // [15, 16, 17, 18]
	assertExtract(t, r, 15) // [16, 17, 18]
	assertExtract(t, r, 16) // [17, 18]
	assertExtract(t, r, 17) // [18]
	assertExtract(t, r, 18) // []
	assertEmpty(t, r)       //
}
