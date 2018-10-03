package btree

import (
	"fmt"
	"reflect"
	"testing"
)

func Test0(t *testing.T) {
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

func Test1(t *testing.T) {
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

func Test2(t *testing.T) {
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

	// Cases from "CLRS" (Cormen, Leiserson, Rivest, Stein)
	// See Chapter 18.3 Deleting a key from a B-tree

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

func Test3(t *testing.T) {
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
