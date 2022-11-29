package lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	mFront *ListItem
	mBack  *ListItem
	mSize  int
}

func (l *list) Len() int {
	return l.mSize
}

func (l *list) Front() *ListItem {
	return l.mFront
}

func (l *list) Back() *ListItem {
	return l.mBack
}

func (l *list) PushFront(v interface{}) *ListItem {
	l.mSize++
	node := new(ListItem)
	node.Value = v
	node.Next = l.mFront
	node.Prev = nil
	if l.mFront != nil {
		l.mFront.Prev = node
		l.mFront = node
	} else {
		l.mFront = node
		l.mBack = node
	}

	return node
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.mSize++
	node := new(ListItem)
	node.Value = v
	node.Next = nil
	node.Prev = l.mBack
	if l.mBack != nil {
		l.mBack.Next = node
		l.mBack = node
	} else {
		l.mBack = node
		l.mFront = node
	}

	return node
}

func (l *list) Remove(i *ListItem) {
	l.mSize--
	switch {
	case i.Prev == nil && i.Next != nil:
		l.mFront.Next.Prev = nil
		l.mFront = l.mFront.Next
	case i.Next == nil && i.Prev != nil:
		l.mBack.Prev.Next = nil
		l.mBack = l.mBack.Prev
	case i.Next == nil && i.Prev == nil:
		l.mBack = nil
		l.mFront = nil
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
}

func (l *list) MoveToFront(i *ListItem) {
	switch {
	case i.Prev == nil && i.Next != nil:
		return
	case i.Next == nil && i.Prev != nil:
		l.mBack.Prev.Next = nil
		l.mBack = l.mBack.Prev
	case i.Next == nil && i.Prev == nil:
		return
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	l.mFront.Prev = i
	i.Prev = nil
	i.Next = l.mFront
	l.mFront = i
}

func NewList() List {
	return new(list)
}
