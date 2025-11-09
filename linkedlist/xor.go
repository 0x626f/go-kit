package linkedlist

import (
	"github.com/0x626f/go-kit/abstract"
	"github.com/0x626f/go-kit/raw"
	"github.com/0x626f/go-kit/utils"
	"unsafe"
)

type link uintptr

type XORNode[T any] struct {
	xor  link
	Data T
}

func (xor *XORNode[T]) link() link {
	return link(unsafe.Pointer(xor))
}

type XORLinkedList[T any, I int] struct {
	head, tail link
	size       int
}

func NewXORList[T any]() *XORLinkedList[T, int] {
	return &XORLinkedList[T, int]{}
}

func (xor *XORLinkedList[T, I]) xor(link0, link1 link) link {
	return link0 ^ link1
}

//go:nocheckptr
func (xor *XORLinkedList[T, I]) getNode(node link) *XORNode[T] {
	if node == 0 {
		return nil
	}
	return (*XORNode[T])(unsafe.Pointer(node))
}

func (xor *XORLinkedList[T, I]) createNode() *XORNode[T] {
	node := raw.AllocateBlank[XORNode[T]]()
	node.Data = utils.Zero[T]()
	node.xor = 0
	return node
}

func (xor *XORLinkedList[T, I]) removeNode(node *XORNode[T]) {
	raw.DeAllocate(node)
}

func (xor *XORLinkedList[T, I]) Free() {
	if xor.head == 0 || xor.size == 0 {
		return
	}

	var prev link
	iterator := xor.head

	for iterator != 0 {
		node := xor.getNode(iterator)
		prev, iterator = iterator, xor.xor(prev, node.xor)
		raw.DeAllocate(xor.getNode(prev))
	}
	xor.size = 0
}

func (xor *XORLinkedList[T, I]) insert(data T, back bool) *XORNode[T] {
	node := xor.createNode()
	node.Data = data

	nodeLink := node.link()

	if back {
		if xor.head == 0 {
			xor.head = nodeLink
		} else if xor.tail == 0 {
			xor.tail = nodeLink

			node.xor = xor.head
			head := xor.getNode(xor.head)

			head.xor = xor.tail
		} else {
			tail := xor.getNode(xor.tail)

			tail.xor ^= nodeLink
			node.xor = xor.tail

			xor.tail = nodeLink
		}
	} else {
		if xor.head == 0 {
			xor.head = nodeLink
		} else if xor.tail == 0 {
			xor.tail = xor.head

			node.xor = xor.tail
			tail := xor.getNode(xor.head)

			xor.head = nodeLink
			tail.xor = xor.head
		} else {
			head := xor.getNode(xor.head)

			head.xor ^= nodeLink
			node.xor = xor.head

			xor.head = nodeLink
		}
	}

	xor.size++
	return node
}

func (xor *XORLinkedList[T, I]) deleteByIndex(index int) {
	left, candidate, right := xor.findWithNeighbors(index)
	xor.delete(left, candidate, right)
}

func (xor *XORLinkedList[T, I]) delete(left, candidate, right link) {
	if candidate == 0 {
		return
	}

	leftNode, candidateNode, rightNode := xor.getNode(left), xor.getNode(candidate), xor.getNode(right)

	if left != 0 {
		leftNode.xor = xor.xor(xor.xor(leftNode.xor, candidate), right)
	}

	if right != 0 {
		rightNode.xor = xor.xor(xor.xor(rightNode.xor, candidate), left)
	}

	raw.DeAllocate(candidateNode)

	if candidate == xor.tail {
		xor.tail = left
	}

	if candidate == xor.head {
		xor.head = right
	}

	xor.size--
}

func (xor *XORLinkedList[T, I]) calcAbsoluteIndex(index int) (int, bool) {
	if xor.head == 0 {
		return 0, false
	}

	idx := index
	if index < 0 {
		idx = xor.size + index
	}

	if idx < 0 || idx >= xor.size {
		return 0, false
	}

	return idx, true
}

func (xor *XORLinkedList[T, I]) findWithNeighbors(index int) (link, link, link) {

	idx, exists := xor.calcAbsoluteIndex(index)

	if !exists {
		return 0, 0, 0
	}

	if idx == 0 {
		if xor.head == 0 {
			return 0, 0, 0
		}
		node := xor.getNode(xor.head)
		return 0, xor.head, node.xor
	}

	med := xor.size / 2

	var node *XORNode[T]

	if idx < med {

		var left link
		iterator := xor.head

		for idx != 0 {
			node = xor.getNode(iterator)
			left, iterator = iterator, xor.xor(left, node.xor)
			idx--
		}

		node = xor.getNode(iterator)
		return left, iterator, xor.xor(left, node.xor)
	}

	idx = xor.size - 1 - idx

	var right link
	iterator := xor.tail

	for idx != 0 {
		node = xor.getNode(iterator)
		right, iterator = iterator, xor.xor(right, node.xor)
		idx--
	}

	node = xor.getNode(iterator)
	return xor.xor(right, node.xor), iterator, right
}

func (xor *XORLinkedList[T, I]) Size() int {
	return xor.size
}

func (xor *XORLinkedList[T, I]) IsEmpty() bool {
	return xor.size == 0
}

func (xor *XORLinkedList[T, I]) At(index int) T {
	_, ref, _ := xor.findWithNeighbors(index)

	if ref == 0 {
		return utils.Zero[T]()
	}

	node := xor.getNode(ref)
	return node.Data
}

func (xor *XORLinkedList[T, I]) Get(depth int) T {
	return xor.At(depth)
}

func (xor *XORLinkedList[T, I]) Push(data T) {
	_ = xor.insert(data, true)
}

func (xor *XORLinkedList[T, I]) PushAll(data ...T) {
	for _, value := range data {
		_ = xor.insert(value, true)
	}
}

func (xor *XORLinkedList[T, I]) Join(collection abstract.Collection[T, int]) {
	collection.ForEach(func(index int, data T) bool {
		_ = xor.insert(data, true)
		return true
	})
}

func (xor *XORLinkedList[T, I]) Merge(collection abstract.Collection[T, int]) abstract.Collection[T, int] {
	list := NewXORList[T]()

	var index int
	var prev link
	iterator := xor.head

	for iterator != 0 {
		node := xor.getNode(iterator)

		list.Push(node.Data)

		prev, iterator = iterator, xor.xor(prev, node.xor)
		index++
	}

	collection.ForEach(func(index int, data T) bool {
		list.Push(data)
		return true
	})

	return list
}

func (xor *XORLinkedList[T, I]) Delete(index int) {
	xor.deleteByIndex(index)
}

func (xor *XORLinkedList[T, I]) DeleteBy(predicate abstract.Predicate[T]) {
	var prev, next link
	iterator := xor.head

	for iterator != 0 {
		node := xor.getNode(iterator)
		next = xor.xor(prev, node.xor)

		if predicate(node.Data) {
			xor.delete(prev, iterator, next)
			iterator = next
		} else {
			prev, iterator = iterator, next
		}

	}

}

func (xor *XORLinkedList[T, I]) DeleteAll() {
	xor.Free()
}

func (xor *XORLinkedList[T, I]) Some(predicate abstract.Predicate[T]) bool {
	var prev link
	iterator := xor.head

	for iterator != 0 {
		node := xor.getNode(iterator)

		if predicate(node.Data) {
			return true
		}

		prev, iterator = iterator, xor.xor(prev, node.xor)
	}

	return false
}

func (xor *XORLinkedList[T, I]) Find(predicate abstract.Predicate[T]) (T, bool) {
	var prev link
	iterator := xor.head

	for iterator != 0 {
		node := xor.getNode(iterator)

		if predicate(node.Data) {
			return node.Data, true
		}

		prev, iterator = iterator, xor.xor(prev, node.xor)
	}

	return utils.Zero[T](), false
}

func (xor *XORLinkedList[T, I]) Filter(predicate abstract.Predicate[T]) abstract.Collection[T, int] {
	list := NewXORList[T]()

	var prev link
	iterator := xor.head

	for iterator != 0 {
		node := xor.getNode(iterator)

		if predicate(node.Data) {
			list.Push(node.Data)
		}

		prev, iterator = iterator, xor.xor(prev, node.xor)
	}

	return list
}

func (xor *XORLinkedList[T, I]) ForEach(receiver abstract.IndexedReceiver[int, T]) {
	var index int
	var prev link
	iterator := xor.head

	for iterator != 0 {
		node := xor.getNode(iterator)

		if !receiver(index, node.Data) {
			break
		}
		prev, iterator = iterator, xor.xor(prev, node.xor)
		index++
	}
}

func (xor *XORLinkedList[T, I]) First() T {
	if xor.size == 0 {
		return utils.Zero[T]()
	}
	return xor.At(0)
}

func (xor *XORLinkedList[T, I]) Last() T {
	if xor.size == 0 {
		return utils.Zero[T]()
	}
	return xor.At(-1)
}

func (xor *XORLinkedList[T, I]) PopLeft() T {
	if xor.head == 0 {
		return utils.Zero[T]()
	}

	headNode := xor.getNode(xor.head)

	if headNode.xor != 0 {
		nextNode := xor.getNode(headNode.xor)
		nextNode.xor ^= xor.head
	}

	xor.head = headNode.xor

	if xor.head == xor.tail {
		xor.tail = 0
	}

	data := headNode.Data
	raw.DeAllocate(headNode)

	xor.size--

	return data
}

func (xor *XORLinkedList[T, I]) PopRight() T {
	if xor.size == 0 {
		return utils.Zero[T]()
	}

	if xor.tail == 0 {
		return xor.PopLeft()
	}

	tailNode := xor.getNode(xor.tail)

	newTailNode := xor.getNode(tailNode.xor)
	newTailNode.xor ^= xor.tail

	if tailNode.xor != xor.head {
		xor.tail = tailNode.xor
	} else {
		xor.tail = 0
	}

	data := tailNode.Data
	raw.DeAllocate(tailNode)

	xor.size--

	return data
}

func (xor *XORLinkedList[T, I]) Shrink(capacity int) {
	if capacity >= xor.size {
		return
	}

	if capacity == 0 {
		xor.DeleteAll()
		return
	}

	count := xor.size - capacity
	iterator := xor.tail

	for count != 0 {
		node := xor.getNode(iterator)
		left := xor.xor(node.xor, 0)

		xor.delete(left, iterator, 0)

		iterator = left
		count--
	}
}
