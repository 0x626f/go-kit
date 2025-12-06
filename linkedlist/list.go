package linkedlist

import (
	"github.com/0x626f/go-kit/abstract"
	"github.com/0x626f/go-kit/utils"
)

type LinkedList[D any] struct {
	LinkedListBase[int, D]
}

type LinkedListBase[I int, D any] struct {
	head, tail *LinkedNode[D]
	size       int
}

type LinkedNode[D any] struct {
	left, right *LinkedNode[D]
	Data        D
}

func NewLinkedList[D any]() *LinkedList[D] {
	return &LinkedList[D]{}
}

func (list *LinkedListBase[I, D]) insert(data D, back bool) *LinkedNode[D] {
	node := &LinkedNode[D]{Data: data}

	if list.head == nil {
		list.head = node
	} else if list.tail == nil {
		if back {
			list.tail = node
			list.tail.left = list.head
			list.head.right = list.tail
		} else {
			list.tail = list.head
			list.head = node

			list.tail.left = list.head
			list.head.right = list.tail
		}
	} else {
		if back {
			list.tail.right = node
			node.left = list.tail
			list.tail = node
		} else {
			list.head.left = node
			node.right = list.head
			list.head = node
		}
	}

	list.size++

	return node
}

func (list *LinkedListBase[I, D]) deleteByIndex(index int) {
	list.Remove(list.findNodeByIndex(index))
}

func (list *LinkedListBase[I, D]) Remove(node *LinkedNode[D]) {
	if node == nil {
		return
	}

	if node.left != nil {
		node.left.right = node.right
	}

	if node.right != nil {
		node.right.left = node.left
	}

	if list.tail == node {
		list.tail = node.left
	}

	if list.head == node {
		list.head = node.right
	}
	list.size--
}

func (list *LinkedListBase[I, D]) calcAbsoluteIndex(index int) (int, bool) {
	if list.size == 0 {
		return index, false
	}

	idx := index
	if index < 0 {
		idx = list.size + index
	}

	if idx < 0 || idx >= list.size {
		return index, false
	}

	return idx, true
}

func (list *LinkedListBase[I, D]) findNodeByIndex(index int) *LinkedNode[D] {

	idx, exists := list.calcAbsoluteIndex(index)

	if !exists {
		return nil
	}

	if idx == 0 {
		if list.head == nil {
			return nil
		}
		return list.head
	}

	med := list.size / 2

	if idx < med {
		iterator := list.head

		for idx != 0 {
			iterator = iterator.right
			idx--
		}

		return iterator
	}

	idx = list.size - 1 - idx

	iterator := list.tail

	for idx != 0 {
		iterator = iterator.left
		idx--
	}

	return iterator
}

func (list *LinkedListBase[I, D]) Size() int {
	return list.size
}

func (list *LinkedListBase[I, D]) IsEmpty() bool {
	return list.size == 0
}

func (list *LinkedListBase[I, D]) At(index int) D {
	node := list.findNodeByIndex(index)

	if node == nil {
		return utils.Zero[D]()
	}

	return node.Data
}

func (list *LinkedListBase[I, D]) Get(index int) D {
	return list.At(index)
}

func (list *LinkedListBase[I, D]) Push(data D) {
	_ = list.insert(data, true)
}

func (list *LinkedListBase[I, D]) PushFront(data D) {
	_ = list.insert(data, false)
}

func (list *LinkedListBase[I, D]) PushAll(data ...D) {
	for _, value := range data {
		_ = list.insert(value, true)
	}
}

func (list *LinkedListBase[I, D]) Insert(data D) *LinkedNode[D] {
	return list.insert(data, true)
}

func (list *LinkedListBase[I, D]) InsertFront(data D) *LinkedNode[D] {
	return list.insert(data, false)
}

func (list *LinkedListBase[I, D]) IndexOf(predicate abstract.Predicate[D]) (int, bool) {
	var index int
	iterator := list.head

	for iterator != nil {

		if predicate(iterator.Data) {
			return index, true
		}
		iterator = iterator.right
		index++
	}
	return 0, false
}

func (list *LinkedListBase[I, D]) Join(collection abstract.Collection[int, D]) {
	collection.ForEach(func(index int, data D) bool {
		_ = list.insert(data, true)
		return true
	})
}

func (list *LinkedListBase[I, D]) Merge(collection abstract.Collection[int, D]) abstract.Collection[int, D] {
	merged := NewLinkedList[D]()

	iterator := list.head

	for iterator != nil {
		merged.Push(iterator.Data)
		iterator = iterator.right
	}

	collection.ForEach(func(index int, data D) bool {
		merged.Push(data)
		return true
	})

	return merged
}

func (list *LinkedListBase[I, D]) Delete(index int) {
	list.deleteByIndex(index)
}

func (list *LinkedListBase[I, D]) DeleteBy(predicate abstract.Predicate[D]) {
	iterator := list.head

	for iterator != nil {
		if predicate(iterator.Data) {
			list.Remove(iterator)
		}
		iterator = iterator.right
	}
}

func (list *LinkedListBase[I, D]) DeleteAll() {
	iterator := list.head

	for iterator != nil {
		next := iterator.right
		iterator.left, iterator.right = nil, nil
		iterator = next
	}

	list.head, list.tail = nil, nil
	list.size = 0
}

func (list *LinkedListBase[I, D]) Some(predicate abstract.Predicate[D]) bool {
	iterator := list.head

	for iterator != nil {
		if predicate(iterator.Data) {
			return true
		}
		iterator = iterator.right
	}

	return false
}

func (list *LinkedListBase[I, D]) Find(predicate abstract.Predicate[D]) (D, bool) {
	iterator := list.head

	for iterator != nil {
		if predicate(iterator.Data) {
			return iterator.Data, true
		}
		iterator = iterator.right
	}

	return utils.Zero[D](), false
}

func (list *LinkedListBase[I, D]) Filter(predicate abstract.Predicate[D]) abstract.Collection[int, D] {
	filtered := NewLinkedList[D]()

	iterator := list.head
	for iterator != nil {
		if predicate(iterator.Data) {
			filtered.Push(iterator.Data)
		}
		iterator = iterator.right
	}

	return filtered
}

func (list *LinkedListBase[I, D]) ForEach(receiver abstract.IndexedReceiver[int, D]) {
	var index int
	iterator := list.head

	for iterator != nil {
		if !receiver(index, iterator.Data) {
			break
		}
		iterator = iterator.right
		index++
	}
}

func (list *LinkedListBase[I, D]) First() D {
	if list.head == nil {
		return utils.Zero[D]()
	}
	return list.At(0)
}

func (list *LinkedListBase[I, D]) Last() D {
	if list.head == nil {
		return utils.Zero[D]()
	}
	return list.At(-1)
}

func (list *LinkedListBase[I, D]) Pop(index int) D {
	node := list.findNodeByIndex(index)

	if node == nil {
		return utils.Zero[D]()
	}

	if node.left != nil {
		node.left.right = node.right
	}

	if node.right != nil {
		node.right.left = node.left
	}

	if list.tail == node {
		list.tail = node.left
	}

	if list.head == node {
		list.head = node.right
	}

	list.size--
	return node.Data
}

func (list *LinkedListBase[I, D]) Swap(i, j int) {
	node0 := list.findNodeByIndex(i)

	if node0 == nil {
		return
	}

	node1 := list.findNodeByIndex(j)

	if node1 == nil || node0 == node1 {
		return
	}

	left0, right0 := node0.left, node0.right
	left1, right1 := node1.left, node1.right

	if right0 == node1 {
		node1.left = left0
		node1.right = node0
		node0.left = node1
		node0.right = right1

		if left0 != nil {
			left0.right = node1
		}
		if right1 != nil {
			right1.left = node0
		}
	} else if left0 == node1 {
		node0.left = left1
		node0.right = node1
		node1.left = node0
		node1.right = right0

		if left1 != nil {
			left1.right = node0
		}
		if right0 != nil {
			right0.left = node1
		}
	} else {
		node0.left = left1
		node0.right = right1
		node1.left = left0
		node1.right = right0

		if left0 != nil {
			left0.right = node1
		}
		if right0 != nil {
			right0.left = node1
		}
		if left1 != nil {
			left1.right = node0
		}
		if right1 != nil {
			right1.left = node0
		}
	}

	if list.head == node0 {
		list.head = node1
	} else if list.head == node1 {
		list.head = node0
	}

	if list.tail == node0 {
		list.tail = node1
	} else if list.tail == node1 {
		list.tail = node0
	}

}

func (list *LinkedListBase[I, D]) Move(from, to int) {
	i, _ := list.calcAbsoluteIndex(from)
	j, _ := list.calcAbsoluteIndex(to)

	node0, node1 := list.findNodeByIndex(i), list.findNodeByIndex(j)

	list.move(node0, node1, i < j)
}

func (list *LinkedListBase[I, D]) MoveToFront(node0 *LinkedNode[D]) {
	list.move(node0, list.head, false)
}

func (list *LinkedListBase[I, D]) PopLeft() D {
	if list.head == nil {
		return utils.Zero[D]()
	}

	node := list.head

	if list.head.right != nil {
		list.head.right.left = nil
	}

	list.head = list.head.right

	if list.head == list.tail {
		list.tail = nil
	}

	node.right = nil
	list.size--

	return node.Data
}

func (list *LinkedListBase[I, D]) PopRight() D {
	if list.size == 0 {
		return utils.Zero[D]()
	}

	if list.tail == nil {
		return list.PopLeft()
	}

	node := list.tail

	node.left.right = nil

	if node.left != list.head {
		list.tail = node.left
	} else {
		list.tail = nil
	}

	node.left = nil
	list.size--

	return node.Data
}

func (list *LinkedListBase[I, D]) Shrink(capacity int) {
	if capacity >= list.size {
		return
	}

	if capacity == 0 {
		list.DeleteAll()
		return
	}

	count := list.size - capacity
	iterator := list.tail

	for count != 0 {
		next := iterator.left
		list.Remove(iterator)
		iterator = next
		count--
	}
}

func (list *LinkedListBase[I, D]) move(node0, node1 *LinkedNode[D], leftToRight bool) {
	if node0 == nil || node1 == nil || node0 == node1 {
		return
	}

	left0, right0 := node0.left, node0.right
	left1, right1 := node1.left, node1.right

	if left0 == node1 || right0 == node1 {
		list.swap(node0, node1)
		return
	}

	if list.head == node0 {
		list.head = node0.right
	} else if list.head == node1 {
		list.head = node0
	}

	if list.tail == node0 {
		list.tail = node0.left
	} else if list.tail == node1 {
		list.tail = node0
	}

	if leftToRight {
		node0.left, node0.right = node1, right1
		node1.left, node1.right = left1, node0

		if right1 != nil {
			right1.left = node0
		}
	} else {
		node0.left, node0.right = left1, node1
		node1.left, node1.right = node0, right1

		if left1 != nil {
			left1.right = node0
		}
	}

	if left0 != nil {
		left0.right = nil
		if right0 != nil {
			left0.right = right0
		}
	}

	if right0 != nil {
		right0.left = nil
		if left0 != nil {
			right0.left = left0
		}
	}
}

func (list *LinkedListBase[I, D]) swap(node0, node1 *LinkedNode[D]) {
	if node0 == nil || node1 == nil || node0 == node1 {
		return
	}

	left0, right0 := node0.left, node0.right
	left1, right1 := node1.left, node1.right

	if right0 == node1 {
		node1.left = left0
		node1.right = node0
		node0.left = node1
		node0.right = right1

		if left0 != nil {
			left0.right = node1
		}
		if right1 != nil {
			right1.left = node0
		}
	} else if left0 == node1 {
		node0.left = left1
		node0.right = node1
		node1.left = node0
		node1.right = right0

		if left1 != nil {
			left1.right = node0
		}
		if right0 != nil {
			right0.left = node1
		}
	} else {
		node0.left = left1
		node0.right = right1
		node1.left = left0
		node1.right = right0

		if left0 != nil {
			left0.right = node1
		}
		if right0 != nil {
			right0.left = node1
		}
		if left1 != nil {
			left1.right = node0
		}
		if right1 != nil {
			right1.left = node0
		}
	}

	if list.head == node0 {
		list.head = node1
	} else if list.head == node1 {
		list.head = node0
	}

	if list.tail == node0 {
		list.tail = node1
	} else if list.tail == node1 {
		list.tail = node0
	}

}

func (list *LinkedListBase[I, D]) Sort(comparator abstract.Comparator[D]) {
	if list.head == nil {
		return
	}

	list.head = quickSort(list.head, getTail(list.head), comparator)

	// Rebuild left pointers and update tail after sorting
	list.head.left = nil
	curr := list.head
	for curr.right != nil {
		curr.right.left = curr
		curr = curr.right
	}
	curr.left = nil
	if curr != list.head {
		list.tail = curr
		// Find the node before tail
		temp := list.head
		for temp.right != list.tail {
			temp = temp.right
		}
		list.tail.left = temp
	} else {
		list.tail = nil
	}
}

func quickSort[D any](head, tail *LinkedNode[D], comparator abstract.Comparator[D]) *LinkedNode[D] {
	if head == nil || head == tail {
		return head
	}

	newHead, newEnd := partition(head, tail, comparator)

	// If pivot is not the only element
	if newHead != newEnd {
		// Find node before pivot
		temp := newHead
		for temp.right != newEnd {
			temp = temp.right
		}
		temp.right = nil

		// Recursively sort before pivot
		newHead = quickSort(newHead, temp, comparator)

		// Get tail of left part and connect to pivot
		temp = getTail(newHead)
		if temp != nil {
			temp.right = newEnd
		}
	}

	// Recursively sort after pivot
	if newEnd.right != nil {
		rightTail := getTail(newEnd.right)
		newEnd.right = quickSort(newEnd.right, rightTail, comparator)
	}

	return newHead
}

func partition[D any](head, end *LinkedNode[D], comparator abstract.Comparator[D]) (*LinkedNode[D], *LinkedNode[D]) {
	if head == nil || end == nil {
		return head, end
	}

	pivot := end
	prev, curr := (*LinkedNode[D])(nil), head
	tail := pivot

	for curr != nil && curr != pivot {
		next := curr.right
		if comparator(curr.Data, pivot.Data) < 0 {
			// Keep in left partition
			if prev == nil {
				head = curr
			} else {
				prev.right = curr
			}
			prev = curr
			curr.right = next
		} else {
			// Move to right partition
			if prev != nil {
				prev.right = next
			} else {
				head = next
			}
			curr.right = nil
			tail.right = curr
			tail = curr
		}
		curr = next
	}

	// Connect left partition to pivot
	if prev == nil {
		head = pivot
	} else {
		prev.right = pivot
	}

	return head, pivot
}

func getTail[D any](head *LinkedNode[D]) *LinkedNode[D] {
	if head == nil {
		return nil
	}
	for head.right != nil {
		head = head.right
	}
	return head
}
