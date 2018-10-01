package main

import (
	"fmt"
	"reflect"
	"sort"
)

// from "CLRS" (Cormen, Leiserson, Rivest, Stein):
// use a fixed integer t ≥ 2, the minimum degree of the B-tree
// root: 1..2*t-1
// node: t-1..2*t-1
// t = [2, ...] (50..2000)

const KEY_EMPTY int = -1

// const VALUE_EMPTY string = ""

type BTree struct {
	leaf bool    // is a leaf
	n    int     // number of keys
	root *BTree  // root
	c    []BTree // children
	key  []int   // keys
	t    int     // minimum degree
}

func allocateNode(t int) *BTree {
	n := new(BTree)
	n.t = t
	n.key = make([]int, 2*t-1)
	n.c = make([]BTree, 2*t)
	return n
}

func diskRead(x *BTree) {
	fmt.Println("DISK-READ")
}

func diskWrite(x *BTree) {
	fmt.Println("DISK-WRITE")
}

func btreeCreate(b *BTree, t int) {
	x := allocateNode(t)
	x.leaf = true
	x.n = 0
	x.t = t // save minimum degree to a new node x
	diskWrite(x)
	b.root = x
	b.t = t // and save minimum degree to b
}

func compareKeys(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func centerText(s string, w int) string {
	return fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(s))/2, s))
}

func btreeToMap(r BTree, m map[int][][]int) {
	/*
	   Convert B-Tree to a map:

	                |5|
	               /   \
	             |2|   |8|
	            /  \   |  \_
	          /    |    \   \
	       |0|1| |3|4| |6|7| |9|

	   btreeMap = {
	       5: [2, 8],
	       2: [[0, 1], [3, 4]],
	       8: [[6, 7]. [9]],
	   }
	*/
	for i := 0; i < r.n+1; i++ {
		if r.c != nil {
			if i < r.n {
				if r.c[i].n > 0 {
					m[r.key[i]] = append(m[r.key[i]], r.c[i].key[:r.c[i].n])
				}
			} else {
				if r.c[i].n > 0 {
					m[r.key[i-1]] = append(m[r.key[i-1]], r.c[i].key[:r.c[i].n])
				}
			}
			btreeToMap(r.c[i], m)
		}
	}
}

func btreeSubPrint(r BTree, m map[int][]*BTree, t int, level int) []int {
	keys := make([]int, 2*t-1)
	keyIndex := 0
	// children count = keys count + 1 (c0|k0|c1)
	for i := 0; i < r.n+1; i++ {
		if i < r.n {
			keys[keyIndex] = r.key[i]
			keyIndex++
		}
		if r.c != nil {
			btreeSubPrint(r.c[i], m, t, level+1)
		}
	}
	// add only nodes with keys
	if r.n != 0 {
		m[level] = append(m[level], &r)
	}
	return keys
}

func btreePrint(b BTree) {
	/*
	          [T]
	           |
	       [T->root]
	    [c0][k0]...[cx]
	     /           \
	   ...           ...
	*/
	if b.root == nil {
		return
	}
	r := b.root

	m := make(map[int][]*BTree)
	btreeSubPrint(*r, m, b.t, 0)

	// sort keys
	sortedKeys := make([]int, 0, len(m))
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Ints(sortedKeys)

	printEnable := false
	maxLineLen := 0
	for _, k := range sortedKeys {
		lineLen := 0
		line := " "
		if printEnable {
			fmt.Printf("%d: {\n", k)
			fmt.Printf("\t")
		}
		for _, v := range m[k] {
			// find parent nodes
			for pk := 0; pk < k; pk++ {
				for pvi, pv := range m[pk] { // pvs = [c|k|c] [] []
					for pki := 0; pki < pv.n+1; pki++ {
						if compareKeys(pv.c[pki].key, v.key) {
							var fromKeyIndex int
							if pki > m[pk][pvi].n-1 {
								fromKeyIndex = pki - 1
								fmt.Sprintf("%s", fmt.Sprintf("(%v|->)", m[pk][pvi].key[fromKeyIndex]))
							} else {
								fromKeyIndex = pki
								fmt.Sprintf("%s", fmt.Sprintf("(|%v->)", m[pk][pvi].key[fromKeyIndex]))
							}
							break
						}
					}
				}
			}

			if v.n != 0 {
				for ki := 0; ki < v.n; ki++ {
					lineLen += len(fmt.Sprintf("|%d", v.key[ki]))
					line += fmt.Sprintf("|%d", v.key[ki])
				}
				lineLen += len("|")
				line += "|"
			}
			lineLen += len(" ")
			line += " "
		}
		if printEnable {
			fmt.Printf("%s len=%d", line[:len(line)-1], lineLen)
		}
		if maxLineLen < lineLen {
			maxLineLen = lineLen
		}
		if printEnable {
			fmt.Printf("\n")
			fmt.Printf("}\n")
		}
	}
	if printEnable {
		fmt.Printf("len=%d\n", maxLineLen)
	}

	// print out
	for _, k := range sortedKeys {
		line := ""
		for _, v := range m[k] {
			if v.n != 0 {
				for ki := 0; ki < v.n; ki++ {
					// skip if marked deleted
					if v.key[ki] == KEY_EMPTY {
						// pass
					} else {
						line += fmt.Sprintf("|%d", v.key[ki])
					}
				}
				line += "|"
			}
			line += " "
		}
		fmt.Printf("%s\n", centerText(line, maxLineLen))
	}
}

func btreeInsert(b *BTree, k int) {
	/*
	          [b]
	           |
	       [b->root]
	    [c0][k0]...[cx]
	     /           \
	   ...           ...
	*/
	fmt.Printf("INSERT key %d\n", k)
	r := b.root
	if r.n == 2*b.t-1 {
		// root is full - allocate new root for splitting
		s := allocateNode(b.t)
		b.root = s
		// new root have some children (i.e. not a leaf)
		s.leaf = false
		s.n = 0
		// s.c[0] = *r // moved after btreeSplitChild
		// split the former root
		btreeSplitChild(s, 0, r)
		s.c[0] = *r // TODO: review this
		btreeInsertNonfull(s, k)
	} else {
		btreeInsertNonfull(r, k)
	}

}

func btreeInsertNonfull(x *BTree, k int) {
	fmt.Printf("INSERT NOT NULL key %d\n", k)
	i := x.n - 1
	if x.leaf {
		// shift keys right to the first key less than k
		for i >= 0 && k < x.key[i] {
			x.key[i+1] = x.key[i]
			i--
		}
		x.key[i+1] = k
		x.n++
		diskWrite(x)
	} else {
		// find index of a child to insert k
		// keys are ordered, find first key less than k
		for i >= 0 && k < x.key[i] {
			i--
		}
		i++
		diskRead(&x.c[i])
		//fmt.Printf("%d|l=%d|n=%d\n", i, len(x.c), x.n)
		if x.c[i].n == 2*x.t-1 {
			// split child node
			btreeSplitChild(x, i, &x.c[i])
			if k > x.key[i] {
				i++
			}
		}
		btreeInsertNonfull(&x.c[i], k)
	}
}

// x - new root, y - old root
func btreeSplitChild(x *BTree, i int, y *BTree) {
	fmt.Printf("SPLIT\n")
	/*
	   T = 2
	   [c0]k1[c1]k2[c3]k3[c4]

	   old:
	   []2[]4[]9[] (leaf = true)
	   insert 5 - split as n = 3 (T * 2 - 1)

	   new:
	       [c0]?[]?[]?[]        --- x, leaf = false
	        /
	   []2[]4[]9[]              --- y, leaf = true

	   // allocate z and copy keys after T from old root (y)

	              []9?[]?[]?[]  --- z, leaf = true too (as y)
	                | copy from old root keys at indices [T..2*T-1)
	                | and copy children at indices [T..2*T)

	   // set used keys number to T - 1 (to cut keys from old root), NB: keys are still there

	   []2[]?[]?[]              --- y, leaf = true

	   // update root

	   // shift all children to the right of i, as c[i + 1] = *z
	       [c0]?[c1]?[]?[]      --- x, leaf = false
	        /     \
	   []2[]?[]?[] \            --- y, leaf = true
	               []9?[]?[]?[] --- z, leaf = true

	   // shift all keys to the right of i
	   // set new root's key[i] (i=0) to old root's middle key[T - 1]
	       [c0]4[c1]?[]?[]      --- x, leaf = false
	        /     \
	   []2[]?[]?[] \            --- y, leaf = true
	               []9?[]?[]?[] --- z, leaf = true
	*/

	// allocate new node for splitting
	z := allocateNode(x.t)
	z.leaf = y.leaf
	z.n = x.t - 1
	// copy keys (from the right)
	for j := 0; j < x.t-1; j++ {
		z.key[j] = y.key[j+x.t]
	}
	// copy children if y is not a leaf
	if !y.leaf {
		for j := 0; j < x.t; j++ {
			z.c[j] = y.c[j+x.t]
		}
	}
	y.n = x.t - 1

	// now update x (new root)
	// shift children to the right of i
	// [a]2[b]3[c]?[], x.n = 2, i = 0
	for j := x.n + 0; j >= i+1; j-- {
		//fmt.Printf(">%-v\n", j)
		x.c[j+1] = x.c[j]
	}
	// [a]2[?]3[b]?[c]
	// and set i + 1 child to z
	x.c[i+1] = *z
	// [a]2[z]3[b]?[c]

	// now shift keys to the right of i
	// [a]2[z]3[b]?[c], x.n = 2, i = 0
	for j := x.n - 1; j >= i; j-- {
		x.key[j+1] = x.key[j]
	}
	// [a]?[z]2[b]3[c]
	x.key[i] = y.key[x.t-1]
	// [a]y.key[z]2[b]3[c]

	x.n++
	diskWrite(y)
	diskWrite(z)
	diskWrite(x)
}

// search a key in BTree
func btreeSearch(t BTree, k int) (*BTree, int) {
	fmt.Printf("SEARCH %d\n", k)
	x := t.root
	if x == nil {
		x = &t
	}
	i := 0
	for i < x.n && k > x.key[i] {
		i++
	}
	// fmt.Printf("%d\n", i)
	if i < x.n && k == x.key[i] {
		return x, i
	} else if x.leaf {
		return nil, 0
	} else {
		diskRead(&x.c[i])
		return btreeSearch(x.c[i], k)
	}
	return x, i
}

func removeKey(x *BTree, i int) {
	// mark as deleted with some tombsone marker (KEY_EMPTY)
	// x.key[i] = KEY_EMPTY // mark deleted
	// // x.n--
	// compact
	for j := i; j < x.n-1; j++ {
		x.key[j] = x.key[j+1]
	}
	x.n--
}

func btreeDelete(t *BTree, k int) {
	fmt.Printf("DELETE %d\n", k)
	x := t.root
	if x == nil {
		x = t
	}

	// determine an index i of a root x.c[i] of the subtree with k
	i := 0
	for i < x.n && k > x.key[i] {
		i++
	}
	found := false
	if i < x.n && k == x.key[i] {
		found = true
	}
	fmt.Printf("found = %t, i = %d\n", found, i)

	if found {
		if x.leaf {
			fmt.Printf("1: DELETE leaf at %d key %d\n", i, k)
			removeKey(x, i)
			return
		} else {
			// 2:
			fmt.Printf("2: DELETE not leaf at %d\n", i)
			// y - child that precedes k in x
			y := &x.c[i]
			// z - child that follows k in x
			z := &x.c[i+1]

			fmt.Printf("y.n = %d, z.n = %d\n", y.n, z.n)

			if y.n >= x.t {
				// 2a: if the child y, that precedes k in x, has at least t keys
				// find the predecessor k' of k in the subtree rooted at y
				k2 := y.key[y.n-1]
				fmt.Printf("2a: p: k' = %d\n", k2)
				// recursively delete k', and replace k by k' in x
				btreeDelete(y, k2)
				x.key[i] = k2
			} else if z.n >= x.t {
				// 2b: if y has fewer than t keys and z has at least t keys
				// find successor k' of k in the subtree rooted at z
				k2 := z.key[0]
				fmt.Printf("2b: s: k' = %d\n", k2)
				// recursively delete k', and replace k by k' in x
				btreeDelete(z, k2)
				x.key[i] = k2
			} else if y.n == x.t-1 && z.n == x.t-1 {
				fmt.Printf("2c: merge\n")
				/*
				   Delete key 4

				        T = 3

				      1|c0|4|c1|7|c2
				        /     \    \  | y.n = len(2,3) -> T - 1 = 2
				      2|3     5|6   8 | z.n = len(5,6) -> T - 1 = 2

				   1) move k into y

				      1|c0|4|c1|7|c2
				        /     \    \
				     2|3|4    5|6   8

				   2) move all from z into y

				      1|c0|4|c1|7|c2
				        /     \    \
				   2|3|4|5|6  5|6   8

				   3) remove k=4

				      1|c0|7|c2
				        /     \
				   2|3|4|5|6   8

				   4) remove k=4 from y

				      1|c0|7|c2
				        /     \
				     2|3|5|6   8
				*/
				merge(x, i)
				btreeDelete(y, k)
			}
		}
	} else {
		// the key k is not in the internal node x
		fmt.Printf("3: DELETE at %d\n", i)
		if x.c == nil {
			fmt.Printf("3: DELETE x.c == nil\n")
			return
		}
		fmt.Printf("3: x.c[i=%d].n = %d\n", i, x.c[i].n)
		// if x.c[i] has T - 1 keys
		if x.c[i].n == x.t-1 {
			fmt.Printf("3: DELETE at %d, x.c[i=%d].n == T - 1 = %d\n", i, i, x.c[i].n)
			// 3a: but sibling has T
			if i > 0 && x.c[i-1].n == x.t {
				fmt.Printf("3a(left): %d vs %d\n", x.c[i-1].n, x.c[i].n)
				/*
				   Delete key 5

				        T = 3

				       c0|4|c1
				       /     \       | x.c[i].n = len(5,6) -> T - 1 = 2
				    1|2|3    5|6     | x.c[i - 1].n = len(1,2,3) -> T = 3

				   1) move key from x down into x.c[i]

				       c0|4|c1
				       /     \
				    1|2|3    4|5|6

				   2) move a key from left sibling (x.c[i - 1]) up into x
				   and remove key at index n - 1 from x.c[i - 1]

				       c0|3|c1
				       /     \
				     1|2     4|5|6

				   3) delete key

				       c0|3|c1
				       /     \
				     1|2     4|6
				*/
				leftSibling := &x.c[i-1]

				// shift all keys right
				for j := x.c[i].n + 0; /*1*/ j > 0; j-- {
					x.c[i].key[j] = x.c[i].key[j-1]
				}
				// move key from x into x.c[i]
				x.c[i].key[0] = x.key[i-1]

				// update pointers to children
				if !leftSibling.leaf {
					for j := x.c[i].n + 1; /*2*/ j > 0; j-- {
						x.c[i].c[j] = x.c[i].c[j-1]
					}
					// TODO: review -- update old x's child (that is root key now) with new pointer
					x.c[i].c[0] = leftSibling.c[leftSibling.n] // was
					// leftSibling.c[leftSibling.n] = x.c[i].c[0]
				}
				x.c[i].n++
				// move a key from left sibling (x.c[i - 1]) up into x
				x.key[i-1] = leftSibling.key[leftSibling.n-1]
				// and remove key (here only reduce n) at index n - 1 from x.c[i - 1]
				// leftSibling.n--
				// TODO: review
				removeKey(leftSibling, leftSibling.n-1)

			} else if i < x.n && x.c[i+1].n == x.t {
				fmt.Printf("3a(right): %d vs %d\n", x.c[i].n, x.c[i+1].n)
				/*
				   Delete key 1

				        T = 3

				       c0|3|c1
				       /     \         | x.c[i].n = len(1,2) -> T - 1 = 2
				    1|2      5|6|7     | x.c[i + 1].n = len(5,6,7) -> T = 3

				   1) move key from x down into x.c[i]

				      c0|3|c1
				      /     \
				   1|2|3    5|6|7

				   2) move key from right sibling (x.c[i + 1]) up into x
				   and remove key at index 0 from x.c[i + 1]

				      c0|5|c1
				      /     \
				   1|2|3    6|7

				   3) delete key

				      c0|5|c1
				      /     \
				    2|3     6|7
				*/
				rightSibling := &x.c[i+1]

				// move key from x down into x.c[i]
				x.c[i].key[x.c[i].n+0 /*1*/] = x.key[i]
				x.c[i].n++
				// move key from right sibling (x.c[i + 1]) up into x
				x.key[i] = rightSibling.key[0]

				// remove key at index 0 from x.c[i + 1]
				removeKey(rightSibling, 0)

				// update pointer to x.c[i]
				if !rightSibling.leaf {
					// TODO: review -- update old x's child (that is root key now) with new pointer
					x.c[i].c[x.c[i].n+1] = rightSibling.c[0] // was
					// rightSibling.c[0] = x.c[i].c[x.c[i].n + 1]
					// shift pointers
					for j := 0; j < rightSibling.n+1; j++ {
						rightSibling.c[j] = rightSibling.c[j+1]
					}
				}
			} else {
				if i > 0 {
					fmt.Printf("3b: %d vs %d\n", x.c[i-1].n, x.c[i].n)
				} else {
					fmt.Printf("3b: %d vs %d\n", x.c[i].n, x.c[i+1].n)
				}
				/*
				   Delete key 4

				        T = 3

				       c0|3|c1
				       /     \        | x.c[i].n = len(1,2) -> T - 1 = 2
				     1|2     4|5      | x.c[i + 1].n = len(4,5) -> T - 1 = 2

				   1) i = 1 (left sibling)
				   set c0.n = 2 * T - 1 = 5

				       c0|3|c1
				       /     \
				   1|2|?|?|?  4|5

				   set c0.key[T=3] to x.key[0]

				       c0|3|c1
				       /     \
				   1|2|3|?|?  4|5

				   copy all keys from c1 (T - 1 = 2 keys)

				       c0|3|c1
				       /     \
				   1|2|3|4|5  4|5
				*/

				// merge
				if i > 0 {
					// left sibling
					merge(x, i-1)
					// update i of child (x.c) to go into on the next iteration
					i = i - 1
				} else if i < x.n {
					// right sibling
					merge(x, i)
				}
			}
		}
		btreeDelete(&x.c[i], k)
	}
}

// merge siblings (deleting a key - case 3b)
func merge(x *BTree, i int) {
	fmt.Printf("3b: merge (i=%d)\n", i)
	c0 := &x.c[i]
	c1 := &x.c[i+1]
	// update num of keys in c0
	c0.n = 2*x.t - 1
	// copy key from the root
	c0.key[x.t-1] = x.key[i] // TODO: x.t - 1 or x.t ?????
	// copy all keys from c1
	for i := 0; i < x.t-1; i++ {
		c0.key[x.t+i] = c1.key[i]
	}
	// if not a leaf - copy pointers to children
	if !c0.leaf {
		for i := 0; i < x.t; i++ {
			c0.c[x.t+i] = c1.c[i]
		}
	}

	// reduce root keys number
	x.n--
	// shift root keys and pointers to children left
	for j := i; j < x.n; j++ { // TODO: < or <= ?????
		x.key[j] = x.key[j+1]
		x.c[j+1] = x.c[j+2]
	}

	// set links to empty nodes
	c1 = &BTree{}
	if x.n == 0 {
		fmt.Printf("x.n == 0 -- clear\n")
		x = &BTree{}
	}
}

func test0() {
	fmt.Println("-----")
	fmt.Println("Test0")
	fmt.Println("-----")
	const T int = 2
	var btree BTree
	btree = BTree{
		leaf: false,
		n:    1,
		root: &btree,
		c: []BTree{
			BTree{leaf: true, n: 1, root: &BTree{}, key: []int{1}, t: T},
			BTree{leaf: true, n: 1, root: &BTree{}, key: []int{3}, t: T},
		},
		key: []int{2},
		t:   T,
	}
	// btreePrint(btree)
	btreeMap := make(map[int][][]int)
	btreeToMap(*btree.root, btreeMap)
	fmt.Printf("btreeMap=%v\n", btreeMap)
	/*
	      |2|
	     /   \
	   |1|   |3|
	*/
	test := map[int][][]int{
		2: [][]int{{1}, {3}},
	}
	if !reflect.DeepEqual(btreeMap, test) {
		fmt.Printf("ERROR: btreeMap is not as expected\n")
		return
	}
}

func test1() {
	fmt.Println("-----")
	fmt.Println("Test1")
	fmt.Println("-----")
	// TEST for T = 2
	// https://www.cs.utexas.edu/users/djimenez/utsa/cs3343/lecture17.html
	// create BTree from data
	const T int = 2
	btree := BTree{}
	btreeCreate(&btree, T)
	keys := []int{
		5, 9, 3, 7, 1, 2, 8, 6, 0, 4,
	}
	for _, k := range keys {
		btreeInsert(&btree, k)
	}
	r := btree.root
	btreeMap := make(map[int][][]int)
	btreeToMap(*r, btreeMap)
	fmt.Printf("btreeMap=%v\n", btreeMap)
	/*
	            |5|
	           /   \
	          /     \
	        |2|      |8|
	       /   \     /  \
	   |0|1| |3|4| |6|7| |9|
	*/
	test := map[int][][]int{
		5: [][]int{{2}, {8}},
		2: [][]int{{0, 1}, {3, 4}},
		8: [][]int{{6, 7}, {9}},
	}
	if !reflect.DeepEqual(btreeMap, test) {
		fmt.Printf("ERROR: btreeMap is not as expected\n")
		return
	}
}

func test2() {
	fmt.Println("-----")
	fmt.Println("Test2")
	fmt.Println("-----")
	// TEST for T = 3
	// create BTree from data
	const T int = 3
	btree := BTree{}
	btreeCreate(&btree, T)
	keys := []int{
		10, 20, 30, 40, 1, 3, 4, 5, 12, 14,
		21, 22, 32, 33, 34, 35, 36, 41, 42,
		6, 37, 23, 15, // 999,
	}
	for _, k := range keys {
		btreeInsert(&btree, k)
		btreePrint(btree)
	}

	fmt.Println("--------------------------")

	btreeDelete(&btree, 21) // case 1
	btreePrint(btree)

	btreeDelete(&btree, 20) // case 2a
	btreePrint(btree)

	btreeInsert(&btree, 24) // case 1
	btreePrint(btree)

	btreeDelete(&btree, 15) // case 2b
	btreePrint(btree)

	btreeDelete(&btree, 10) // case 2c
	btreePrint(btree)

	btreeDelete(&btree, 5) // case 3b
	btreePrint(btree)

	btreeDelete(&btree, 3) // case 3a
	btreePrint(btree)

	btreeDelete(&btree, 35) // case 1
	btreePrint(btree)

	btreeDelete(&btree, 33) // case 3b
	btreePrint(btree)

	btreeDelete(&btree, 40) // case 3b
	btreePrint(btree)

	fmt.Println("--------------------------")
}

func test3() {
	fmt.Println("-----")
	fmt.Println("Test3")
	fmt.Println("-----")

	const T int = 3
	btree := BTree{}
	btreeCreate(&btree, T)

	keys := []int{
		16, 3, 7, 13, 1, 2, 4, 5, 6,
		10, 11, 12, 14, 15, 20, 23,
		17, 18, 19, 21, 22, 24, 26,
	}

	for _, i := range keys {
		btreeInsert(&btree, i)
		btreePrint(btree)
	}

	btreeDelete(&btree, 5)
	btreePrint(btree)

	btreeDelete(&btree, 16)
	btreePrint(btree)

	btreeMap := make(map[int][][]int)
	btreeToMap(*btree.root, btreeMap)
	fmt.Printf("btreeMap=%v\n", btreeMap)
	/*
	                            |12|
	                        |3|7| |18|21|
	   |1|2| |4|6| |10|11| |13|14|15|17| |19|20| |22|23|24|26|
	*/
	test := map[int][][]int{
		12: [][]int{{3, 7}, {18, 21}},
		3:  [][]int{{1, 2}},
		7:  [][]int{{4, 6}, {10, 11}},
		18: [][]int{{13, 14, 15, 17}},
		21: [][]int{{19, 20}, {22, 23, 24, 26}},
	}
	if !reflect.DeepEqual(btreeMap, test) {
		fmt.Printf("ERROR: btreeMap is not as expected\n")
		return
	}

	var node *BTree
	var index int

	node, index = btreeSearch(btree, 7)
	fmt.Println(node.key[index])
	if node.key[index] != 7 {
		fmt.Printf("ERROR: btreeSearch(btree, 7) != 7\n")
		return
	}
	// check prev key
	if node.key[index-1] != 3 {
		fmt.Printf("ERROR: node.key[index - 1] != 3\n")
		return
	}

	node, index = btreeSearch(btree, 4)
	fmt.Println(node.key[index])
	if node.key[index] != 4 {
		fmt.Printf("ERROR: btreeSearch(btree, 4) != 4\n")
		return
	}
	// check next key
	if node.key[index+1] != 6 {
		fmt.Printf("ERROR: node.key[index + 1] != 6\n")
		return
	}
}

func tests() {
	test0()
	test1()
	test2()
	test3()
}

func main() {
	fmt.Printf("=========================================\n")
	tests()
	fmt.Printf("=========================================\n")
}
