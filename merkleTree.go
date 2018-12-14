package main

import "crypto/sha256"


//默克尔树节点
type MerkleTree struct{
	RootNode *MerkleNode
}

//默克尔根节点
type MerkleNode struct{
	Left *MerkleNode
	Right *MerkleNode
	Data []byte
}
//生成默克尔树中的节点，如果是叶子节点，则Left，right为nil ，如果为非叶子节点，根据Left，right生成当前节点的hash
func NewMerkleNode(left,right *MerkleNode,data []byte) *MerkleNode{
	mnode := MerkleNode{}

	if left ==nil && right==nil{
		mnode.Data = data
	}else{
		prevhashes := append(left.Data,right.Data...)
		firsthash:= sha256.Sum256(prevhashes)
		hash:=sha256.Sum256(firsthash[:])
		mnode.Data = hash[:]
	}

	mnode.Left = left
	mnode.Right = right

	return &mnode
}

//构建默克尔树
func NewMerkleTree(data [][]byte) *MerkleTree{
	var nodes []MerkleNode
	//构建叶子节点。
	for _,datum := range data{
		node:= NewMerkleNode(nil,nil,datum)
		nodes =  append(nodes,*node)
	}
	//j代表的是某一层的第一个元素
	j:=0
	//第一层循环代表 nSize代表某一层的个数，每循环一次减半
	for nSize :=len(data);nSize >1;nSize = (nSize+1)/2{
		//第二条循环i+=2代表两两拼接。 i2是为了当个数是基数的时候，拷贝最后的元素。
		for i:=0 ; i<nSize ;i+=2{
			i2 := min(i+1,nSize-1)

			node := NewMerkleNode(&nodes[j+i],&nodes[j+i2],nil)
			nodes = append(nodes,*node)
		}
		//j代表的是某一层的第一个元素
		j+=nSize

	}

	mTree := MerkleTree{&(nodes[len(nodes)-1])}
	return &mTree
}
