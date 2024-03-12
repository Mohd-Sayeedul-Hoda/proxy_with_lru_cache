package cache

import(
  "net/http"
)

// trying to create double linked list to store and
// remove the node from the Node

// and using map to access the linked list so

type Node struct{
  Key string 
  Header http.Header
  Body []byte
  prev *Node
  next *Node
}


type LRUCache struct{
  Nodes (map[string]*Node)
  Head *Node
  Tail *Node
  Capacity uint64
  CurrentSize uint64
}

func (lru *LRUCache) Get(Key string) (*Node, bool){
  node, exists:= lru.Nodes[Key]
  if !exists {
    return nil, false
  }
  
  lru.MoveFront(node)
  return node, true
}


func (lru *LRUCache) MoveFront(node *Node){
  if node == lru.Head{
    return
  }else if node == lru.Tail{
    lru.Tail = node.prev
    node.prev.next = nil
  }else{
    node.prev.next = node.next
    node.next.prev = node.prev
  }
  node.prev = nil
  node.next = lru.Head
  lru.Head = node
}

func (lru *LRUCache) Put(key string, header http.Header, bodyBytes []byte){
  node := CreateNode(key, header, bodyBytes)
  lru.AddNode(node)
}

func (lru *LRUCache) AddNode(node *Node){
  if lru.CurrentSize >= lru.Capacity{
    lru.RemoveLast()
  }
  node.next = lru.Head
  node.prev = nil
  lru.Head = node
  lru.CurrentSize++
  lru.Nodes[node.Key] = node
}

func (lru *LRUCache) RemoveLast(){
  lru.Tail.prev.next = nil
  delete(lru.Nodes, lru.Tail.Key)
  lru.Tail = lru.Tail.prev
  lru.CurrentSize--
}

func CreateNode(key string, header http.Header, bodyBytes []byte)*Node{
  return &Node{
    Key: key,
    Header: header,
    Body: bodyBytes,
    next: nil,
    prev: nil,
  }
}

