// Package linkedlist provides a generic doubly-linked list implementation with O(1)
// insertion/deletion operations and efficient bidirectional traversal.
//
// The linked list supports:
//   - Insertion at front or back in O(1) time
//   - Deletion by index or predicate
//   - Access by index with optimized traversal
//   - Functional operations (filter, find, forEach, etc.)
//   - In-place sorting using quicksort
//   - Node manipulation for cache implementations
//
// The list uses bidirectional links, allowing efficient traversal from either end
// and enabling optimizations like accessing elements closer to the back by traversing
// from the tail.
package linkedlist

import (
	"github.com/0x626f/go-kit/abstract"
	"github.com/0x626f/go-kit/utils"
)

// LinkedList is a generic doubly-linked list that can store elements of any type.
// It provides constant-time insertion and deletion at both ends, and efficient
// indexed access through optimized bidirectional traversal.
//
// Type parameters:
//   - D: The type of data stored in the list
//
// Example:
//
//	list := linkedlist.NewLinkedList[int]()
//	list.Push(1)
//	list.Push(2)
//	list.PushFront(0)
//	// List now contains: 0, 1, 2
type LinkedList[D any] struct {
	LinkedListBase[int, D]
}

// LinkedListBase is the underlying implementation of the doubly-linked list.
// It maintains references to both head and tail for O(1) operations at both ends.
//
// Type parameters:
//   - I: Index type (always int in practice)
//   - D: The type of data stored in nodes
type LinkedListBase[I int, D any] struct {
	// head points to the first node in the list, or nil if empty
	head, tail *LinkedNode[D]
	// size tracks the current number of elements in the list
	size int
}

// LinkedNode represents a single node in the doubly-linked list.
// Each node stores data and maintains bidirectional links to adjacent nodes.
//
// Type parameters:
//   - D: The type of data stored in this node
type LinkedNode[D any] struct {
	// left points to the previous node (nil if this is the head)
	left, right *LinkedNode[D]
	// Data holds the value stored in this node
	Data D
}

// NewLinkedList creates and initializes a new empty linked list.
//
// Type parameters:
//   - D: The type of data to store in the list
//
// Returns:
//   - A pointer to the newly created LinkedList
//
// Example:
//
//	list := linkedlist.NewLinkedList[string]()
//	list.Push("hello")
//	list.Push("world")
func NewLinkedList[D any]() *LinkedList[D] {
	return &LinkedList[D]{}
}

// insert is an internal method that adds a new node to the list.
// When back is true, inserts at the tail; when false, inserts at the head.
//
// Parameters:
//   - data: The data to store in the new node
//   - back: If true, insert at tail; if false, insert at head
//
// Returns:
//   - A pointer to the newly created node
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

// deleteByIndex is an internal method that removes a node at the specified index.
//
// Parameters:
//   - index: The index of the node to remove
func (list *LinkedListBase[I, D]) deleteByIndex(index int) {
	list.Remove(list.findNodeByIndex(index))
}

// Remove removes a specific node from the list.
// This is used internally by cache implementations that maintain references to nodes.
//
// Parameters:
//   - node: The node to remove from the list
//
// Time complexity: O(1)
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

// calcAbsoluteIndex converts a potentially negative index to an absolute position.
// Negative indices count from the end (-1 is the last element, -2 is second-to-last, etc.).
//
// Parameters:
//   - index: The index to convert (can be negative)
//
// Returns:
//   - The absolute index and true if valid, or the original index and false if out of bounds
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

// findNodeByIndex finds a node at the specified index using optimized bidirectional traversal.
// If the index is closer to the head, traverses from head; if closer to tail, traverses from tail.
//
// Parameters:
//   - index: The index to find (supports negative indices)
//
// Returns:
//   - The node at the specified index, or nil if index is out of bounds
//
// Time complexity: O(n/2) on average due to bidirectional optimization
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

// Size returns the number of elements in the list.
//
// Returns:
//   - The current size of the list
//
// Time complexity: O(1)
func (list *LinkedListBase[I, D]) Size() int {
	return list.size
}

// IsEmpty checks whether the list contains any elements.
//
// Returns:
//   - true if the list is empty, false otherwise
//
// Time complexity: O(1)
func (list *LinkedListBase[I, D]) IsEmpty() bool {
	return list.size == 0
}

// At retrieves the element at the specified index.
// Supports negative indices (-1 for last element, -2 for second-to-last, etc.).
//
// Parameters:
//   - index: The index to access (can be negative)
//
// Returns:
//   - The element at the index, or a zero value if index is out of bounds
//
// Time complexity: O(n/2) average due to bidirectional traversal optimization
//
// Example:
//
//	list := linkedlist.NewLinkedList[int]()
//	list.PushAll(10, 20, 30)
//	fmt.Println(list.At(0))   // Output: 10
//	fmt.Println(list.At(-1))  // Output: 30
func (list *LinkedListBase[I, D]) At(index int) D {
	node := list.findNodeByIndex(index)

	if node == nil {
		return utils.Zero[D]()
	}

	return node.Data
}

// Get is an alias for At. Retrieves the element at the specified index.
//
// Parameters:
//   - index: The index to access (can be negative)
//
// Returns:
//   - The element at the index, or a zero value if index is out of bounds
func (list *LinkedListBase[I, D]) Get(index int) D {
	return list.At(index)
}

// Push appends an element to the end of the list.
//
// Parameters:
//   - data: The element to append
//
// Time complexity: O(1)
func (list *LinkedListBase[I, D]) Push(data D) {
	_ = list.insert(data, true)
}

// PushFront inserts an element at the beginning of the list.
//
// Parameters:
//   - data: The element to insert at the front
//
// Time complexity: O(1)
func (list *LinkedListBase[I, D]) PushFront(data D) {
	_ = list.insert(data, false)
}

// PushAll appends multiple elements to the end of the list in order.
//
// Parameters:
//   - data: Variable number of elements to append
//
// Time complexity: O(k) where k is the number of elements to add
//
// Example:
//
//	list := linkedlist.NewLinkedList[int]()
//	list.PushAll(1, 2, 3, 4, 5)
//	// List now contains: 1, 2, 3, 4, 5
func (list *LinkedListBase[I, D]) PushAll(data ...D) {
	for _, value := range data {
		_ = list.insert(value, true)
	}
}

// Insert adds an element to the end of the list and returns the created node.
// This is used by cache implementations that need to maintain node references.
//
// Parameters:
//   - data: The element to insert
//
// Returns:
//   - A pointer to the newly created node
//
// Time complexity: O(1)
func (list *LinkedListBase[I, D]) Insert(data D) *LinkedNode[D] {
	return list.insert(data, true)
}

// InsertFront adds an element to the beginning of the list and returns the created node.
// This is used by cache implementations that need to maintain node references.
//
// Parameters:
//   - data: The element to insert
//
// Returns:
//   - A pointer to the newly created node
//
// Time complexity: O(1)
func (list *LinkedListBase[I, D]) InsertFront(data D) *LinkedNode[D] {
	return list.insert(data, false)
}

// IndexOf finds the index of the first element matching the predicate.
//
// Parameters:
//   - predicate: A function that returns true for the desired element
//
// Returns:
//   - The index of the first matching element and true if found
//   - 0 and false if no element matches
//
// Time complexity: O(n)
//
// Example:
//
//	list := linkedlist.NewLinkedList[int]()
//	list.PushAll(10, 20, 30)
//	idx, found := list.IndexOf(func(x int) bool { return x == 20 })
//	// idx = 1, found = true
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

// Join appends all elements from another collection to this list.
// This modifies the current list in place.
//
// Parameters:
//   - collection: The collection whose elements to append
//
// Time complexity: O(k) where k is the size of the collection to join
//
// Example:
//
//	list1 := linkedlist.NewLinkedList[int]()
//	list1.PushAll(1, 2, 3)
//	list2 := linkedlist.NewLinkedList[int]()
//	list2.PushAll(4, 5, 6)
//	list1.Join(list2)
//	// list1 now contains: 1, 2, 3, 4, 5, 6
func (list *LinkedListBase[I, D]) Join(collection abstract.Collection[int, D]) {
	collection.ForEach(func(index int, data D) bool {
		_ = list.insert(data, true)
		return true
	})
}

// Merge creates a new list containing all elements from this list and another collection.
// The original lists are not modified.
//
// Parameters:
//   - collection: The collection to merge with
//
// Returns:
//   - A new list containing elements from both collections
//
// Time complexity: O(n + k) where n is this list's size and k is the collection's size
//
// Example:
//
//	list1 := linkedlist.NewLinkedList[int]()
//	list1.PushAll(1, 2, 3)
//	list2 := linkedlist.NewLinkedList[int]()
//	list2.PushAll(4, 5, 6)
//	merged := list1.Merge(list2)
//	// merged contains: 1, 2, 3, 4, 5, 6
//	// list1 and list2 are unchanged
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

// Delete removes the element at the specified index.
// Supports negative indices (-1 for last element, etc.).
//
// Parameters:
//   - index: The index of the element to remove
//
// Time complexity: O(n/2) average due to bidirectional traversal
func (list *LinkedListBase[I, D]) Delete(index int) {
	list.deleteByIndex(index)
}

// DeleteBy removes all elements matching the predicate.
//
// Parameters:
//   - predicate: A function that returns true for elements to delete
//
// Time complexity: O(n)
//
// Example:
//
//	list := linkedlist.NewLinkedList[int]()
//	list.PushAll(1, 2, 3, 4, 5)
//	list.DeleteBy(func(x int) bool { return x%2 == 0 })
//	// List now contains: 1, 3, 5
func (list *LinkedListBase[I, D]) DeleteBy(predicate abstract.Predicate[D]) {
	iterator := list.head

	for iterator != nil {
		if predicate(iterator.Data) {
			list.Remove(iterator)
		}
		iterator = iterator.right
	}
}

// DeleteAll removes all elements from the list and clears all node links.
// The list becomes empty after this operation.
//
// Time complexity: O(n)
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

// Some checks if at least one element matches the predicate.
//
// Parameters:
//   - predicate: A function that returns true for matching elements
//
// Returns:
//   - true if any element matches, false otherwise
//
// Time complexity: O(n) in worst case, but returns early on first match
//
// Example:
//
//	list := linkedlist.NewLinkedList[int]()
//	list.PushAll(1, 2, 3)
//	hasEven := list.Some(func(x int) bool { return x%2 == 0 })
//	// hasEven = true (because 2 is even)
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

// Find returns the first element matching the predicate.
//
// Parameters:
//   - predicate: A function that returns true for the desired element
//
// Returns:
//   - The first matching element and true if found
//   - A zero value and false if no element matches
//
// Time complexity: O(n) in worst case, but returns early on first match
//
// Example:
//
//	list := linkedlist.NewLinkedList[int]()
//	list.PushAll(10, 20, 30)
//	val, found := list.Find(func(x int) bool { return x > 15 })
//	// val = 20, found = true
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

// Filter creates a new list containing only elements matching the predicate.
// The original list is not modified.
//
// Parameters:
//   - predicate: A function that returns true for elements to include
//
// Returns:
//   - A new list containing only matching elements
//
// Time complexity: O(n)
//
// Example:
//
//	list := linkedlist.NewLinkedList[int]()
//	list.PushAll(1, 2, 3, 4, 5)
//	evens := list.Filter(func(x int) bool { return x%2 == 0 })
//	// evens contains: 2, 4
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

// ForEach iterates over all elements in the list, calling the receiver function for each.
// If the receiver returns false, iteration stops early.
//
// Parameters:
//   - receiver: A function called for each element with its index and value
//
// Time complexity: O(n)
//
// Example:
//
//	list := linkedlist.NewLinkedList[int]()
//	list.PushAll(1, 2, 3)
//	list.ForEach(func(i int, val int) bool {
//	    fmt.Printf("Index %d: %d\n", i, val)
//	    return true // continue iteration
//	})
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

// First returns the first element in the list.
//
// Returns:
//   - The first element, or a zero value if the list is empty
//
// Time complexity: O(1)
func (list *LinkedListBase[I, D]) First() D {
	if list.head == nil {
		return utils.Zero[D]()
	}
	return list.At(0)
}

// Last returns the last element in the list.
//
// Returns:
//   - The last element, or a zero value if the list is empty
//
// Time complexity: O(1)
func (list *LinkedListBase[I, D]) Last() D {
	if list.head == nil {
		return utils.Zero[D]()
	}
	return list.At(-1)
}

// Pop removes and returns the element at the specified index.
// Supports negative indices (-1 for last element, etc.).
//
// Parameters:
//   - index: The index of the element to remove and return
//
// Returns:
//   - The removed element, or a zero value if index is out of bounds
//
// Time complexity: O(n/2) average due to bidirectional traversal
//
// Example:
//
//	list := linkedlist.NewLinkedList[int]()
//	list.PushAll(10, 20, 30)
//	val := list.Pop(1)
//	// val = 20, list now contains: 10, 30
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

// Swap exchanges the positions of two elements at the specified indices.
// If either index is invalid or the indices are the same, no swap occurs.
//
// Parameters:
//   - i, j: The indices of elements to swap
//
// Time complexity: O(n)
//
// Example:
//
//	list := linkedlist.NewLinkedList[int]()
//	list.PushAll(10, 20, 30)
//	list.Swap(0, 2)
//	// List now contains: 30, 20, 10
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

// Move relocates an element from one position to another in the list.
//
// Parameters:
//   - from: The index of the element to move
//   - to: The destination index
//
// Time complexity: O(n)
//
// Example:
//
//	list := linkedlist.NewLinkedList[int]()
//	list.PushAll(10, 20, 30, 40)
//	list.Move(3, 0)
//	// List now contains: 40, 10, 20, 30
func (list *LinkedListBase[I, D]) Move(from, to int) {
	i, _ := list.calcAbsoluteIndex(from)
	j, _ := list.calcAbsoluteIndex(to)

	node0, node1 := list.findNodeByIndex(i), list.findNodeByIndex(j)

	list.move(node0, node1, i < j)
}

// MoveToFront moves a specific node to the front of the list.
// This is used by cache implementations like LRU cache.
//
// Parameters:
//   - node0: The node to move to the front
//
// Time complexity: O(1)
func (list *LinkedListBase[I, D]) MoveToFront(node0 *LinkedNode[D]) {
	list.move(node0, list.head, false)
}

// PopLeft removes and returns the first element from the list.
//
// Returns:
//   - The first element, or a zero value if the list is empty
//
// Time complexity: O(1)
//
// Example:
//
//	list := linkedlist.NewLinkedList[int]()
//	list.PushAll(10, 20, 30)
//	val := list.PopLeft()
//	// val = 10, list now contains: 20, 30
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

// PopRight removes and returns the last element from the list.
//
// Returns:
//   - The last element, or a zero value if the list is empty
//
// Time complexity: O(1)
//
// Example:
//
//	list := linkedlist.NewLinkedList[int]()
//	list.PushAll(10, 20, 30)
//	val := list.PopRight()
//	// val = 30, list now contains: 10, 20
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

// Shrink reduces the list size to the specified capacity by removing elements from the end.
// If capacity is 0, all elements are removed.
// If capacity is greater than or equal to the current size, no elements are removed.
//
// Parameters:
//   - capacity: The maximum number of elements to keep
//
// Time complexity: O(k) where k is the number of elements to remove
//
// Example:
//
//	list := linkedlist.NewLinkedList[int]()
//	list.PushAll(1, 2, 3, 4, 5)
//	list.Shrink(3)
//	// List now contains: 1, 2, 3
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

// move is an internal method that relocates a node to a new position in the list.
//
// Parameters:
//   - node0: The node to move
//   - node1: The destination node
//   - leftToRight: Direction of movement
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

// swap is an internal method that exchanges two nodes in the list.
//
// Parameters:
//   - node0, node1: The nodes to swap
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

// Sort sorts the list in place using quicksort algorithm.
// The list is reordered according to the provided comparator function.
//
// Parameters:
//   - comparator: A function that compares two elements.
//     Should return:
//   - negative value if first element should come before second
//   - zero if elements are equal
//   - positive value if first element should come after second
//
// Time complexity: O(n log n) average, O(n²) worst case
//
// Example:
//
//	list := linkedlist.NewLinkedList[int]()
//	list.PushAll(30, 10, 20)
//	list.Sort(func(a, b int) int { return a - b })
//	// List now contains: 10, 20, 30
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

// quickSort is an internal function implementing the quicksort algorithm for linked lists.
// It partitions the list and recursively sorts the partitions.
//
// Parameters:
//   - head, tail: The range of nodes to sort
//   - comparator: The comparison function
//
// Returns:
//   - The new head of the sorted range
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

// partition is an internal function that partitions the list around a pivot element.
// Elements less than the pivot are placed before it, others after it.
//
// Parameters:
//   - head, end: The range to partition (end is the pivot)
//   - comparator: The comparison function
//
// Returns:
//   - The new head of the partitioned range and the pivot node
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

// getTail is an internal helper function that finds the last node in a chain.
//
// Parameters:
//   - head: The starting node
//
// Returns:
//   - The last node in the chain, or nil if head is nil
func getTail[D any](head *LinkedNode[D]) *LinkedNode[D] {
	if head == nil {
		return nil
	}
	for head.right != nil {
		head = head.right
	}
	return head
}
