package rings

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

func TestSet_Empty(t *testing.T) {
	b := New[int](10)

	b.Set(0, 1)

	AssertBufferExactly(t, b, []int{1})
}

func TestSet_Middle(t *testing.T) {
	b := New[int](10)

	b.Append(1, 2, 3)
	b.Set(1, 4)

	AssertBufferExactly(t, b, []int{1, 4, 3})
}

func TestSet_End(t *testing.T) {
	b := New[int](10)

	b.Append(1, 2, 3)
	b.Set(3, -1)

	AssertBufferExactly(t, b, []int{1, 2, 3, -1})
}

func TestInsert_Empty(t *testing.T) {
	b := New[int](10)

	b.Insert(0, 1)

	AssertBufferExactly(t, b, []int{1})
}

func TestInsert_AtEnd(t *testing.T) {
	b := New[int](10)

	b.Insert(0, 1)
	b.Insert(1, 2)

	AssertBufferExactly(t, b, []int{1, 2})
}

func TestInsert_AtBeginning(t *testing.T) {
	b := New[int](10)

	b.Insert(0, 1)
	b.Insert(0, 2)

	AssertBufferExactly(t, b, []int{2, 1})
}

func TestInsert_WhenInMiddle(t *testing.T) {
	b := New[int](10)
	b.head = 5

	b.Insert(0, 1)
	b.Insert(0, 2)

	AssertBufferExactly(t, b, []int{2, 1})
}

func TestInsert_OverEdge2(t *testing.T) {
	b := New[int](10)
	b.head = 9

	b.Insert(0, 1)
	b.Insert(1, 2)
	b.Insert(0, 3)
	b.Insert(0, 4)
	b.Insert(0, 5)

	AssertBufferExactly(t, b, []int{5, 4, 3, 1, 2})
}

func TestInsert_OverEdge3(t *testing.T) {
	b := New[int](10)
	b.head = 8
	data := []int{1, 2, 3, 4, 5}

	b.Append(data...)
	b.Insert(3, 6)

	AssertBufferExactly(t, b, []int{1, 2, 3, 6, 4, 5})
}

func TestInsert_OverEdge4(t *testing.T) {
	b := New[int](10)
	b.head = 8
	data := []int{1, 2, 3, 4, 5}

	b.Append(data...)
	b.Insert(1, 6)

	AssertBufferExactly(t, b, []int{1, 6, 2, 3, 4, 5})
}

func TestInsert_OverEdge5(t *testing.T) {
	b := New[int](10)
	data := []int{1, 2, 3, 4, 5}

	b.Append(data...)

	b.Insert(1, 6)

	AssertBufferExactly(t, b, []int{1, 6, 2, 3, 4, 5})
}

func TestRemove_SingleElt(t *testing.T) {
	b := New[int](10)
	data := []int{1}

	b.Append(data...)

	elt := b.Remove(0)

	AssertEquals(t, elt, 1)
	AssertBufferEmpty(t, b)
}

func TestRemove_LastElt(t *testing.T) {
	b := New[int](10)
	data := []int{1, 2}

	b.Append(data...)

	elt := b.Remove(1)

	AssertEquals(t, elt, 2)
	AssertBufferExactly(t, b, []int{1})
}

func TestRemove_FirstElt(t *testing.T) {
	b := New[int](10)
	data := []int{1, 2}

	b.Append(data...)

	elt := b.Remove(0)

	AssertEquals(t, elt, 1)
	AssertBufferExactly(t, b, []int{2})
}

func TestRemove_WorksWhenIndexOverEnd(t *testing.T) {
	b := New[int](10)
	b.head = 8
	data := []int{1, 2, 3, 4, 5}

	b.Append(data...)

	elt := b.Remove(1)

	AssertEquals(t, elt, 2)
	AssertBufferExactly(t, b, []int{1, 3, 4, 5})
}

func TestRemove_WorksWhenIndexUnderEnd(t *testing.T) {
	b := New[int](10)
	b.head = 8
	data := []int{1, 2, 3, 4, 5}

	b.Append(data...)

	elt := b.Remove(3)

	AssertEquals(t, elt, 4)
	AssertBufferExactly(t, b, []int{1, 2, 3, 5})
}

func TestCopyFrom_WhenB2Overflow(t *testing.T) {
	b := New[int](10)
	b2 := New[int](5)
	b2.head = 4
	data := []int{1, 2, 3, 4, 5}
	data2 := []int{6, 7, 8, 9}
	b.Append(data...)
	b2.Append(data2...)

	b.CopyFrom(b2)

	AssertBufferExactly(t, b, []int{1, 2, 3, 4, 5, 6, 7, 8, 9})
}

func TestCopyFrom_WhenPlentyOfSpace(t *testing.T) {
	b := New[int](10)
	b2 := New[int](5)
	data := []int{1, 2, 3, 4, 5}
	data2 := []int{6, 7, 8, 9}
	b.Append(data...)
	b2.Append(data2...)

	b.CopyFrom(b2)

	AssertBufferExactly(t, b, []int{1, 2, 3, 4, 5, 6, 7, 8, 9})
}

func TestCopyFrom_WhenEmpty(t *testing.T) {
	b := New[int](10)
	b2 := New[int](10)

	b.CopyFrom(b2)

	AssertBufferEmpty(t, b)
}

func TestTruncLast_WhenEmpty(t *testing.T) {
	b := New[int](10)

	b.TruncLast(10)

	AssertBufferEmpty(t, b)
}

func TestTruncLast_RemovesAll(t *testing.T) {
	b := New[int](10)

	b.Append(1, 2)
	b.TruncLast(2)

	AssertBufferEmpty(t, b)
}

func TestTruncLast_RemovesNotAll(t *testing.T) {
	b := New[int](10)

	b.Append(1, 2, 3)
	b.TruncLast(1)

	AssertBufferExactly(t, b, []int{1, 2})
}

func TestTruncLast_0RemovesNone(t *testing.T) {
	b := New[int](10)

	b.Append(1, 2, 3)
	b.TruncLast(0)

	AssertBufferExactly(t, b, []int{1, 2, 3})

}

func TestTruncLast_RemovesSomeWhenFull(t *testing.T) {
	b := New[int](3)

	b.Append(1, 2, 3)
	b.TruncLast(1)

	AssertBufferExactly(t, b, []int{1, 2})

}

func TestTruncLast_RemovesAllWhenFull(t *testing.T) {
	b := New[int](3)

	b.Append(1, 2, 3)
	b.TruncLast(4)

	AssertBufferEmpty(t, b)
}

func TestTruncFirst_WhenEmpty(t *testing.T) {
	b := New[int](10)

	b.TruncFirst(10)

	AssertBufferEmpty(t, b)
}

func TestTruncFirst_RemovesAll(t *testing.T) {
	b := New[int](10)

	b.Append(1, 2)
	b.TruncFirst(2)

	AssertBufferEmpty(t, b)
}

func TestTruncFirst_RemovesNotAll(t *testing.T) {
	b := New[int](10)

	b.Append(1, 2, 3)
	b.TruncFirst(1)

	AssertBufferExactly(t, b, []int{2, 3})
}

func TestTruncFirst_0RemovesNone(t *testing.T) {
	b := New[int](10)

	b.Append(1, 2, 3)
	b.TruncFirst(0)

	AssertBufferExactly(t, b, []int{1, 2, 3})

}

func TestTruncFirst_RemovesSomeWhenFull(t *testing.T) {
	b := New[int](3)

	b.Append(1, 2, 3)
	b.TruncFirst(1)

	AssertBufferExactly(t, b, []int{2, 3})

}

func TestTruncFirst_RemovesAllWhenFull(t *testing.T) {
	b := New[int](3)

	b.Append(1, 2, 3)
	b.TruncFirst(4)

	AssertBufferEmpty(t, b)
}

func TestAsSlice_WhenSome(t *testing.T) {
	b := New[int](5)
	data := []int{1, 2, 3}

	b.Append(data...)

	AssertSliceExactly(t, b.AsSlice(), data)
}

func TestAsSlice_WhenFull(t *testing.T) {
	b := New[int](5)
	data := []int{1, 2, 3, 4, 5}

	b.Append(data...)

	AssertSliceExactly(t, b.AsSlice(), data)
}

func TestAsSlice_WhenEmpty(t *testing.T) {
	b := New[int](5)

	AssertSliceExactly(t, b.AsSlice(), []int{})
}

func TestAsSlice_WhenOverEnd(t *testing.T) {
	b := New[int](5)
	b.head = 3
	data := []int{1, 2, 3}

	b.Append(data...)

	AssertSliceExactly(t, b.AsSlice(), data)
}

func TestAsSlice_WhenOverEndFull(t *testing.T) {
	b := New[int](5)
	b.head = 3
	data := []int{1, 2, 3, 4, 5}

	b.Append(data...)

	AssertSliceExactly(t, b.AsSlice(), data)
}

func AssertEquals[T constraints.Integer](t *testing.T, have T, expect T) {
	if have != expect {
		t.Errorf("expected to have %d, got %d", expect, have)
	}
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

func AssertBufferEmpty[T constraints.Integer](t *testing.T, b *RingBuffer[T]) {
	if b.Len() != 0 {
		t.Error("expected buffer to be empty")
	}
	// TODO: Dump items in b
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
