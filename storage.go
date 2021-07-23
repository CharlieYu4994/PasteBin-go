package main

type node struct {
	prev *node
	key  string
	data *paste
	next *node
}

type LinkedHashMap struct {
	table map[string]*node
	head  *node
	tail  *node
	cap   int
}

func NewLHM(length int) *LinkedHashMap {
	return &LinkedHashMap{
		table: make(map[string]*node, length),
		head:  nil,
		tail:  nil,
		cap:   length,
	}
}

func (l *LinkedHashMap) Add(key string, data *paste) {
	tmp := &node{
		prev: l.tail,
		key:  key,
		data: data,
		next: nil,
	}

	if l.cap <= 0 {
		l.Delete(l.head.key)
	}
	l.cap--

	l.table[key] = tmp
	if l.head == nil && l.tail == nil {
		l.head = tmp
		l.tail = tmp
	} else {
		l.appendToTail(tmp)
	}
}

func (l *LinkedHashMap) Delete(key string) bool {
	tmp, ok := l.table[key]
	if !ok {
		return ok
	}

	delete(l.table, key)
	l.remove(tmp)
	l.cap++

	return ok
}

func (l *LinkedHashMap) Get(key string) (*paste, bool) {
	tmp, ok := l.table[key]
	if !ok {
		return nil, ok
	}

	if tmp != l.tail {
		l.remove(tmp)
		l.appendToTail(tmp)
	}

	return tmp.data, ok
}

func (l *LinkedHashMap) remove(n *node) {
	if l.head == l.tail {
		l.head = nil
		l.tail = nil
	} else if n == l.head {
		l.head = n.next
		l.head.prev = nil
	} else if n == l.tail {
		l.tail = n.prev
		l.tail.next = nil
	} else {
		n.prev.next = n.next
		n.next.prev = n.prev
	}
}

func (l *LinkedHashMap) appendToTail(n *node) {
	n.prev = l.tail
	n.next = nil
	l.tail.next = n
	l.tail = n
}
