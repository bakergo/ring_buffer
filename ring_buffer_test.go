package main

import (
	"testing"

	"golang.org/x/exp/constraints"
)

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

	AssertBufferExactly(t, b, data)
}

func TestRingBuffer_AppendOverEnd(t *testing.T) {
	b := New[int](10)

	data := []int{1, 2, 3, 4, 5}
	data2 := []int{6, 7, 8, 9, 10, 11, 12, 13}
	expected := []int{4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

	b.Append(data...)
	b.PopFirst(3)
	b.Append(data2...)

	AssertBufferExactly(t, b, expected)
}

func TestRingBuffer_Prepend(t *testing.T) {
	b := New[int](10)

	data := []int{1, 2, 3, 4, 5}
	data2 := []int{6, 7, 8}
	expected := []int{6, 7, 8, 1, 2, 3, 4, 5}

	b.Prepend(data...)
	b.Prepend(data2...)

	AssertBufferExactly(t, b, expected)
}

func TestRingBuffer_PrependOverEnd(t *testing.T) {
	b := New[int](10)

	data := []int{1, 2, 3, 4, 5}
	data2 := []int{6, 7, 8, 9, 10}
	expected := []int{6, 7, 8, 9, 10, 3, 4, 5}

	b.Append(data...)
	b.PopFirst(2)
	b.Prepend(data2...)

	AssertBufferExactly(t, b, expected)
}

func TestPopLast_RemovesLast(t *testing.T) {
	b := New[int](10)

	data := []int{1, 2, 3, 4, 5}
	popExpected := []int{4, 5}
	expected := []int{1, 2, 3}

	b.Append(data...)
	pop := b.PopLast(2)

	AssertBufferExactly(t, b, expected)
	AssertSliceExactly(t, pop, popExpected)
}

func TestPopLast_Truncates(t *testing.T) {
	b := New[int](10)

	data := []int{1, 2, 3, 4, 5}
	popExpected := data

	b.Append(data...)
	pop := b.PopLast(11)

	AssertBufferExactly(t, b, []int{})
	AssertSliceExactly(t, pop, popExpected)
}

func TestPopLastOverEnd(t *testing.T) {
	b := New[int](10)
	b.head = 7

	data := []int{1, 2, 3, 4, 5, 6}
	expected := []int{1, 2}
	popExpected := []int{3, 4, 5, 6}

	b.Append(data...)
	pop := b.PopLast(4)

	AssertBufferExactly(t, b, expected)
	AssertSliceExactly(t, pop, popExpected)
}

func TestPopLast_Empty(t *testing.T) {
	b := New[int](5)

	data := []int{}
	expected := []int{}
	popExpected := []int{}

	b.Append(data...)
	pop := b.PopLast(4)

	AssertBufferExactly(t, b, expected)
	AssertSliceExactly(t, pop, popExpected)
}

func TestPopFirst_RemovesAllItems(t *testing.T) {
	b := New[int](10)

	data := []int{1, 2, 3, 4, 5}
	expected := []int{}
	popExpected := []int{1, 2, 3, 4, 5}

	b.Append(data...)
	pop := b.PopFirst(11)

	AssertBufferExactly(t, b, expected)
	AssertSliceExactly(t, pop, popExpected)
}

func TestPopFirst_RemovesSomeFirst(t *testing.T) {
	b := New[int](10)

	data := []int{1, 2, 3, 4, 5}
	expected := []int{3, 4, 5}
	popExpected := []int{1, 2}

	b.Append(data...)
	pop := b.PopFirst(2)

	AssertBufferExactly(t, b, expected)
	AssertSliceExactly(t, pop, popExpected)
}

func TestPopFirst_OverEnd(t *testing.T) {
	b := New[int](10)
	b.head = 8

	data := []int{1, 2, 3, 4, 5}
	expected := []int{4, 5}
	popExpected := []int{1, 2, 3}

	b.Append(data...)
	pop := b.PopFirst(3)

	AssertSliceExactly(t, pop, popExpected)
	AssertBufferExactly(t, b, expected)
}

func TestPopFirst_Full(t *testing.T) {
	b := New[int](5)

	data := []int{1, 2, 3, 4, 5}
	expected := []int{}
	popExpected := []int{1, 2, 3, 4, 5}

	b.Append(data...)
	pop := b.PopFirst(5)

	AssertBufferExactly(t, b, expected)
	AssertSliceExactly(t, pop, popExpected)
}

func TestPopFirst_Empty(t *testing.T) {
	b := New[int](5)

	data := []int{}
	expected := []int{}
	popExpected := []int{}

	b.Append(data...)
	pop := b.PopFirst(5)

	AssertBufferExactly(t, b, expected)
	AssertSliceExactly(t, pop, popExpected)
}

func TestRingBuffer_AppendRemoveAppend_Appends(t *testing.T) {
	b := New[int](10)
	data := []int{1, 2, 3, 4, 5}
	expected := []int{1, 2, 4, 5}

	b.Append(data...)
	b.Remove(2)

	AssertBufferExactly(t, b, expected)
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

	AssertSliceExactly(t, pop, popExpected)
	AssertBufferExactly(t, b, arrExpected)
}

func AssertSliceExactly[T constraints.Integer](t *testing.T, have []T, expect []T) {
	if len(have) != len(expect) {
		t.Errorf("expected slice to contain %d, got %d", len(expect),
			len(have))
	}
	for i, d := range expect {
		if have[i] != d {
			t.Errorf("expected %d got %d", d, have[i])
		}
	}
}

func AssertBufferExactly[T constraints.Integer](t *testing.T, b *RingBuffer[T], expected []T) {
	if b.Len() != len(expected) {
		t.Errorf("expected buffer to contain %d, got %d", len(expected),
			b.Len())
	}
	for i, d := range expected {
		if b.Get(i) != d {
			t.Errorf("expected %d got %d", d, b.Get(i))
		}
	}
}
