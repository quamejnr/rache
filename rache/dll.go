package rache

type Node[T comparable] struct {
	value T
	next  *Node[T]
	prev  *Node[T]
}

type DLL[T comparable] struct {
	head *Node[T]
	tail *Node[T]
	size int
}

// insert value to the front of the linked list
func (l *DLL[T]) insertFront(val T) {
	n := &Node[T]{value: val}
	if l.head == nil {
		l.head = n
		l.tail = n
	} else {
		n.next = l.head
		l.head.prev = n
		l.head = n
	}
	l.size++
}

// Find a node for a given value
func (l *DLL[T]) find(val T) *Node[T] {
	if l.head == nil {
		return nil
	}
	for c := l.head; c != nil; c = c.next {
		if c.value == val {
			return c
		}
	}
	return nil
}

// remove a node for a given value
func (l *DLL[T]) remove(val T) bool {
	n := l.find(val)
	if n == nil {
		return false
	}
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		l.head = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	} else {
		l.tail = n.prev
	}
	l.size--
	return true
}

// insert value to the back of the linked list
func (l *DLL[T]) insertBack(val T) {
	n := &Node[T]{value: val}
	if l.tail == nil {
		l.head = n
		l.tail = n
	} else {
		n.prev = l.tail
		l.tail.next = n
		l.tail = n
	}
	l.size++
}

// delete value from the back of the linked list
func (l *DLL[T]) deleteBack() (T, bool) {
	if l.head == nil {
		var zero T
		return zero, false
	}
	val := l.tail.value
	l.tail = l.tail.prev
	if l.tail == nil {
		l.head = nil
	} else {
		l.tail.next = nil
	}
	l.size--
	return val, true
}
