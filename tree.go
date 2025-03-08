package ji

import (
	"container/list"
	"context"
	"sort"
	"sync"
)

// TreeNode 表示一个树节点
type TreeNode struct {
	ID       uint        `json:"id"`        // 当前节点 ID
	ParentID uint        `json:"parent_id"` // 父节点 ID
	Name     string      `json:"name"`      // 节点名称
	Selected bool        `json:"selected"`  // 选中状态
	Sorted   int         `json:"sorted"`    // 排序
	Children []*TreeNode `json:"children"`  // 子节点列表
}

// Tree 结构，封装高效的树操作
type Tree struct {
	nodes map[uint]*TreeNode // 快速查找节点
	mu    sync.RWMutex       // 读写锁，确保并发安全
}

// NewTree 创建一个新的树结构
func NewTree(data []TreeNode) (*Tree, error) {
	// 按照 Sorted 字段进行排序，升序
	sort.Slice(data, func(i, j int) bool {
		return data[i].Sorted < data[j].Sorted
	})

	tree := &Tree{nodes: make(map[uint]*TreeNode)}

	// 预先构建节点映射
	for i := range data {
		tree.nodes[data[i].ID] = &data[i]
	}

	// 构建树结构
	for _, node := range tree.nodes {
		if parent, found := tree.nodes[node.ParentID]; found {
			parent.Children = append(parent.Children, node)
		}
	}
	return tree, nil
}

// GetRootNodes 获取所有的根节点（ParentID = 0）
func (tree *Tree) GetRootNodes() []*TreeNode {
	tree.mu.RLock()
	defer tree.mu.RUnlock()

	var rootNodes []*TreeNode
	for _, node := range tree.nodes {
		if node.ParentID == 0 {
			rootNodes = append(rootNodes, node)
		}
	}
	return rootNodes
}

// GetSubCategoryIDs 获取指定节点的所有子节点 ID（广度优先遍历）
func (tree *Tree) GetSubCategoryIDs(nodeID uint, includeSelf bool) []uint {
	tree.mu.RLock()
	defer tree.mu.RUnlock()

	startNode, found := tree.nodes[nodeID]
	if !found {
		return nil
	}

	subCategoryIDs := make([]uint, 0)
	if includeSelf {
		subCategoryIDs = append(subCategoryIDs, startNode.ID)
	}

	queue := list.New()
	queue.PushBack(startNode)

	// 广度优先遍历
	for queue.Len() > 0 {
		element := queue.Front() // 获取队头元素
		currentNode := element.Value.(*TreeNode)
		queue.Remove(element) // 出队

		// 遍历子节点
		for _, child := range currentNode.Children {
			subCategoryIDs = append(subCategoryIDs, child.ID)
			queue.PushBack(child) // 子节点入队
		}
	}

	return subCategoryIDs
}

// TreeLevel 获取某个节点的层级（O(1) 查找）
func (tree *Tree) TreeLevel(nodeID uint) int {
	tree.mu.RLock()
	defer tree.mu.RUnlock()

	level := 0
	currentNode, found := tree.nodes[nodeID]
	if !found {
		return -1 // 如果节点未找到，返回 -1
	}
	// 从当前节点开始，沿着父节点向上遍历，直到根节点
	for currentNode.ParentID != 0 {
		level++
		currentNode = tree.nodes[currentNode.ParentID]
	}
	return level
}

// IsParent 判断 parentID 是否是 nodeID 的直接父节点（O(1) 查询）
func (tree *Tree) IsParent(nodeID, parentID uint) bool {
	tree.mu.RLock()
	defer tree.mu.RUnlock()

	node, found := tree.nodes[nodeID]
	if !found {
		return false
	}
	return node.ParentID == parentID
}

// IsInSubTreeConcurrent 并发查找目标 ID 是否在子树中
func (tree *Tree) IsInSubTreeConcurrent(rootID, targetID uint) bool {
	tree.mu.RLock()
	rootNode, found := tree.nodes[rootID]
	tree.mu.RUnlock()
	if !found {
		return false
	}

	// 创建取消上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 结果通道
	resultChan := make(chan bool, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	// 并发查找子树
	go func() {
		defer wg.Done()
		tree.searchSubtree(ctx, rootNode, targetID, resultChan)
	}()

	// 等待查找完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return <-resultChan
}

// searchSubtree 递归并发搜索子树
func (tree *Tree) searchSubtree(ctx context.Context, node *TreeNode, targetID uint, resultChan chan<- bool) {
	select {
	case <-ctx.Done():
		return
	default:
		// 如果找到了目标节点
		if node.ID == targetID {
			select {
			case resultChan <- true:
			default:
			}
			return
		}

		var wg sync.WaitGroup
		// 并发搜索子节点
		for _, child := range node.Children {
			wg.Add(1)
			go func(childNode *TreeNode) {
				defer wg.Done()
				tree.searchSubtree(ctx, childNode, targetID, resultChan)
			}(child)
		}
		wg.Wait()
	}
}
