package cache

// trying to create double linked list to store and
// remove the node from the Node

// and using map to access the linked list so

type Node struct{
  key string 
  value []byte
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

func (lru *LRUCache) get(key string) ([]byte, bool){
  node, exists:= lru.Nodes[key]
  if !exists {
    return nil, false
  }
  
  lru.MoveFront(node)
  return node.value, true
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

func (lru *LRUCache) AddNode(node *Node){
  node.next = lru.Head
  node.prev = nil
  lru.Head = node
}

func CreateNode(key string, value []byte)*Node{
  return &Node{
    key: key,
    value: value,
    next: nil,
    prev: nil,
  }
}
