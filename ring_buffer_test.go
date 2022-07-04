package main

import "testing"

func TestCap_EqualsSz(t *testing.T) {
	b := New[int](10)
	if b.Cap() != 10 {
		t.Errorf("Expected cap == 10")
	}

	b = New[int](5)
	if b.Cap() != 5 {
		t.Errorf("Expected cap == 5")
	}
}

func TestLen_Init0(t *testing.T) {
	b := New[int](10)
	if b.Len() != 0 {
		t.Errorf("Expected len == 0")
	}
}

func TestRingBuffer_Append_Appends(t *testing.T) {
	b := New[int](10)

	data := []int{1, 2, 3, 4, 5}

	b.Append(data...)

	for i, d := range data {
		if b.Get(i) != d {
			t.Errorf("expected %d got %d", d, b.Get(i))
		}
	}
}

func TestRingBuffer_AppendOverEnd(t *testing.T) {
	b := New[int](10)

	data := []int{1, 2, 3, 4, 5}
	data2 := []int{6, 7, 8, 9, 10, 11, 12, 13}
	expected := []int{4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

	b.Append(data...)
	b.PopFirst(3)
	b.Append(data2...)

	for i, d := range expected {
		if b.Get(i) != d {
			t.Errorf("expected %d got %d", d, b.Get(i))
		}
	}
}

func TestRingBuffer_Prepend(t *testing.T) {
	b := New[int](10)

	data := []int{1, 2, 3, 4, 5}
	data2 := []int{6, 7, 8}
	expected := []int{6, 7, 8, 1, 2, 3, 4, 5}

	b.Prepend(data...)
	b.Prepend(data2...)

	for i, d := range expected {
		if b.Get(i) != d {
			t.Errorf("expected %d got %d", d, b.Get(i))
		}
	}
}

func TestRingBuffer_PrependOverEnd(t *testing.T) {
	b := New[int](10)

	data := []int{1, 2, 3, 4, 5}
	data2 := []int{6, 7, 8, 9, 10}
	expected := []int{6, 7, 8, 9, 10, 3, 4, 5}

	b.Append(data...)
	b.PopFirst(2)
	b.Prepend(data2...)

	for i, d := range expected {
		if b.Get(i) != d {
			t.Errorf("expected %d got %d", d, b.Get(i))
		}
	}
}

func TestPopLast_RemovesLast(t *testing.T) {
	b := New[int](10)

	data := []int{1, 2, 3, 4, 5}
	popExpected := []int{4, 5}
	expected := []int{1, 2, 3}

	b.Append(data...)
	pop := b.PopLast(2)

	for i, d := range expected {
		if b.Get(i) != d {
			t.Errorf("expected %d got %d", d, b.Get(i))
		}
	}
	for i, d := range popExpected {
		if pop[i] != d {
			t.Errorf("expected %d got %d", d, pop[i])
		}
	}
}

func TestPopLast_Truncates(t *testing.T) {
	b := New[int](10)

	data := []int{1, 2, 3, 4, 5}
	popExpected := data

	b.Append(data...)
	pop := b.PopLast(11)
	if b.Len() != 0 {
		t.Error("expected array to be clear now")
	}
	if len(pop) != len(popExpected) {
		t.Errorf("Expected len to be %d was %d", len(popExpected), len(pop))
	}
	for i, d := range popExpected {
		if pop[i] != d {
			t.Errorf("expected %d got %d", d, pop[i])
		}
	}
}

func TestPopLastOverEnd(t *testing.T) {
	b := New[int](10)
	b.head = 7

	data := []int{1, 2, 3, 4, 5, 6}
	expected := []int{1, 2}
	popExpected := []int{3, 4, 5, 6}

	b.Append(data...)
	pop := b.PopLast(4)

	if b.Len() != len(expected) {
		t.Errorf("expected buffer to contain %d, got %d", len(expected),
			b.Len())
	}
	for i, d := range expected {
		if b.Get(i) != d {
			t.Errorf("expected %d got %d", d, b.Get(i))
		}
	}
	if len(pop) != len(popExpected) {
		t.Errorf("Expected popped len to be %d was %d", len(popExpected), len(pop))
	}
	for i, d := range popExpected {
		if pop[i] != d {
			t.Errorf("expected %d got %d", d, pop[i])
		}
	}
}

func TestRingBuffer_AppendRemoveAppend_Appends(t *testing.T) {
	b := New[int](10)
	data := []int{1, 2, 3, 4, 5}
	expected := []int{1, 2, 4, 5}

	b.Append(data...)
	b.Remove(2)

	for i, d := range expected {
		if b.Get(i) != d {
			t.Errorf("expected %d got %d", d, b.Get(i))
		}
	}
}

func TestRingBuffer_AppendRemoveAll(t *testing.T) {
	b := New[int](10)
	data := []int{1, 2, 3, 4, 5}

	b.Append(data...)
	b.Remove(2)
	b.Remove(2)
	b.Remove(2)
	b.PopLast(2)

	if b.Len() != 0 {
		t.Errorf("expected len == 0")
	}
}

func TestRingBuffer_PopLeft(t *testing.T) {
	b := New[int](10)
	data := []int{1, 2, 3, 4, 5}
	popExpected := []int{1, 2}
	arrExpected := []int{3, 4, 5}

	b.Append(data...)
	pop := b.PopFirst(2)

	for i, d := range popExpected {
		if pop[i] != d {
			t.Errorf("expected %d got %d", d, pop[i])
		}
	}
	for i, d := range arrExpected {
		if b.Get(i) != d {
			t.Errorf("expected %d got %d", d, b.Get(i))
		}
	}
}
