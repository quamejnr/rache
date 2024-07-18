package rache

import "testing"

// Test doubly linked list
func TestDLL(t *testing.T) {
	l := DLL[int]{}
	l.insertFront(1)
	l.insertFront(2)
	l.insertFront(3)
	l.insertFront(4)

	t.Run("Test insert front", func(t *testing.T) {

		if l.head.value != 4 {
			t.Errorf("wanted %d got %d", 4, l.head.value)
		}
		if l.tail.value != 1 {
			t.Errorf("wanted %d got %d", 1, l.tail.value)
		}
		if l.head.next.value != 3 {
			t.Errorf("wanted %d got %d", 3, l.head.next.value)
		}
	})

	t.Run("Test find value that exists", func(t *testing.T) {
		want := 3
		n := l.find(want)
		if n.value != want {
			t.Errorf("wanted %d got %d", want, n.value)
		}
		if n.next.value != 2 {
			t.Errorf("wanted %d got %d", 2, n.value)
		}
		if n.prev.value != 4 {
			t.Errorf("wanted %d got %d", 4, n.value)
		}
	})

	t.Run("Test find value that does not exist", func(t *testing.T) {
		n := l.find(6)
		if n != nil {
			t.Errorf("wanted nil got %v", n)
		}
	})
	t.Run("Test remove value that exists", func(t *testing.T) {
		want := 3
		ok := l.remove(want)
		if !ok {
			t.Errorf("wanted true got %v", ok)
		}
		n := l.find(want)
		if n != nil {
			t.Errorf("wanted nil got %v", n)
		}
	})
	t.Run("Test remove value that don't exist", func(t *testing.T) {
		want := 5
		ok := l.remove(want)
		if ok {
			t.Errorf("wanted false got %v", ok)
		}
	})
	t.Run("Test delete from back of list", func(t *testing.T) {
		v, ok := l.deleteBack()
		if !ok {
			t.Errorf("wanted true got %v", ok)
		}
		if v != 1 {
			t.Errorf("wanted %d got %d", 1, v)
		}
		l.deleteBack()
		l.deleteBack()
		l.deleteBack()
		v, ok = l.deleteBack()

	})
	t.Run("Test delete from an empty list", func(t *testing.T) {
		l := DLL[int]{}
		v, ok := l.deleteBack()

		if ok {
			t.Errorf("want false got %v", ok)
		}
		if v != 0 {
			t.Errorf("wanted 0 got %d", v)
		}
	})
	t.Run("Test insert from the back of the list", func(t *testing.T) {
		l := DLL[int]{}
		l.insertBack(1)
		l.insertBack(2)
		l.insertBack(3)
		l.insertBack(4)

		if l.head.value != 1 {
			t.Errorf("wanted 1 got %d", l.head.value)
		}
		if l.head.next.value != 2 {
			t.Errorf("wanted 2 got %d", l.head.next.value)
		}
		if l.tail.value != 4 {
			t.Errorf("wanted 4 got %d", l.tail.value)
		}
		if l.tail.prev.value != 3 {
			t.Errorf("wanted 3 got %d", l.tail.prev.value)
		}
		if l.tail.next != nil {
			t.Errorf("wanted nil got %v", l.tail.next)
		}
	})
}
