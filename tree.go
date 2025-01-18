package ji

import "sync"

// TreeNode 表示一个树节点
// 泛型 T 用于支持任意类型的数据
type TreeNode[T any] struct {
	ID       int            // 当前节点 ID
	PID      int            // 父节点 ID
	Name     string         // 节点名称
	Data     T              // 存储节点数据
	Children []*TreeNode[T] // 子节点列表
}

// BuildTree 构建树形数据结构
func BuildTree[T any](data []TreeNode[T], pid int) []*TreeNode[T] {
	// 使用 map 预先构建每个节点的子节点列表，避免重复查找
	childrenMap := make(map[int][]*TreeNode[T])

	for i := range data {
		childrenMap[data[i].PID] = append(childrenMap[data[i].PID], &data[i])
	}

	// 使用递归构建树
	var build func(pid int) []*TreeNode[T]
	build = func(pid int) []*TreeNode[T] {
		var result []*TreeNode[T]
		for _, node := range childrenMap[pid] {
			node.Children = build(node.ID)
			result = append(result, node)
		}
		return result
	}

	return build(pid)
}

// GetSubCategoryIDs 获取指定节点的所有子节点 ID
func GetSubCategoryIDs[T any](node *TreeNode[T], ids *[]int) {
	if node == nil {
		return
	}
	*ids = append(*ids, node.ID)
	for _, child := range node.Children {
		GetSubCategoryIDs(child, ids)
	}
}

// TreeLevel 获取指定节点的层级
func TreeLevel[T any](id int, nodes []TreeNode[T]) int {
	// 预先构建节点的父子关系索引
	parentMap := make(map[int]int)
	for _, node := range nodes {
		parentMap[node.ID] = node.PID
	}

	// 遍历父节点，查找层级
	level := 0
	currentID := id
	for currentID != 0 {
		currentID = parentMap[currentID]
		level++
	}
	return level
}

// IsParent 判断 id 节点是否为 parentID 节点的子节点
func IsParent[T any](id int, parentID int, nodes []TreeNode[T]) bool {
	// 预先构建树
	subTree := BuildTree(nodes, parentID)
	return IsInSubTreeConcurrent(id, subTree)
}

// IsInSubTreeConcurrent 使用并发判断某个 ID 是否在子树中
func IsInSubTreeConcurrent[T any](id int, subTree []*TreeNode[T]) bool {
	var wg sync.WaitGroup
	var found bool
	var mu sync.Mutex

	// 使用队列（非递归方式）避免递归过深
	queue := subTree
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		wg.Add(1)
		go func(node *TreeNode[T]) {
			defer wg.Done()

			mu.Lock()
			if found {
				mu.Unlock()
				return
			}
			mu.Unlock()

			if node.ID == id {
				mu.Lock()
				found = true
				mu.Unlock()
				return
			}

			// 处理子节点
			if len(node.Children) > 0 {
				mu.Lock()
				queue = append(queue, node.Children...)
				mu.Unlock()
			}
		}(node)
	}

	wg.Wait()
	return found
}
