package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Node struct {
	Value int
	Next  *Node
}

type LinkedList struct {
	FirstNode *Node
}

func NewLinkedList() LinkedList {
	newList := LinkedList{FirstNode: nil}
	return newList
}

func (l *LinkedList) Append(v int) {
	newNode := Node{Value: v, Next: nil}
	if l.FirstNode == nil {
		l.FirstNode = &newNode
	} else {
		lastNode := l.FirstNode
		// fmt.Println("firstNode", l.FirstNode, lastNode)
		// fmt.Println("check", &l.FirstNode, &lastNode)
		for {
			if lastNode.Next == nil {
				lastNode.Next = &newNode
				break
			}
			lastNode = lastNode.Next
		}
	}

}

func (l *LinkedList) delete(v int) bool {
	if l.FirstNode == nil {
		return false
	}
	lastNode := l.FirstNode
	if lastNode.Value == v {
		l.FirstNode = lastNode.Next
		return true
	}
	for {
		nextNode := lastNode.Next
		if nextNode.Value == v {
			lastNode.Next = nextNode.Next
			return true
		}
		if nextNode.Next == nil {
			return false
		}
		lastNode = nextNode
	}

}
func (l *LinkedList) toSlice() []int {
	var s []int
	if l.FirstNode == nil {
		return s
	}
	lastNode := l.FirstNode
	for {
		//fmt.Println(lastNode)
		s = append(s, lastNode.Value)
		if lastNode.Next == nil {
			break
		}
		lastNode = lastNode.Next

	}
	return s

}

type ValuesList struct {
	Values []int `json:"values"`
	Length int   `json:"length"`
}

func (l *LinkedList) toBytes() []byte {

	valuesList := ValuesList{Values: l.toSlice(), Length: len(l.toSlice())}
	b, err := json.Marshal(valuesList)
	if err != nil {
		fmt.Println("error:", err)
	}
	return b
}

func fromBytes(bytes []byte) LinkedList {
	var s ValuesList
	err := json.Unmarshal(bytes, &s)
	if err != nil {
		log.Fatal(err)
	}
	l := NewLinkedList()
	for _, v := range s.Values {
		l.Append(v)
	}
	return l
}

func (l *LinkedList) UnmarshalJSON(b []byte) error {
	var s ValuesList
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	fmt.Println("unmarshal:", s)

	for _, v := range s.Values {
		l.Append(v)
	}

	return nil
}

func (l *LinkedList) MarshalJSON() ([]byte, error) {

	valuesList := ValuesList{Values: l.toSlice(), Length: len(l.toSlice())}
	fmt.Println("vlaueslist: ", valuesList)
	b, err := json.Marshal(valuesList)
	if err != nil {
		fmt.Println("error:", err)
	}
	return b, err
}

func main() {
	l := NewLinkedList()

	for i := 1; i < 16; i++ {
		l.Append(i)
	}
	fmt.Println(l)

	s := l.toSlice()
	fmt.Println(s)

	fmt.Println(l.delete(9))
	fmt.Println(l.delete(112))
	fmt.Println(l.delete(1))

	s = l.toSlice()
	fmt.Println(s)

	bytes := l.toBytes()
	fmt.Println("bytes", string(bytes))

	fromBytes := fromBytes(bytes)
	s = fromBytes.toSlice()
	fmt.Println("reconstructed", s)

	b, err := json.Marshal(&l)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("marshal:", string(b))

	l2 := NewLinkedList()
	err = json.Unmarshal(b, &l2)
	if err != nil {
		fmt.Println("error:", err)
	}

	s2 := l2.toSlice()
	fmt.Println(s2)
}
