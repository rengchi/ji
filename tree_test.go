package ji

import (
	"testing"
)

func TestTreeOperations(t *testing.T) {
	// 树的示例数据
	data := []TreeNode[int]{
		{ID: 1, PID: 0, Name: "Root", Data: 1},
		{ID: 2, PID: 1, Name: "Child1", Data: 2},
		{ID: 3, PID: 1, Name: "Child2", Data: 3},
		{ID: 4, PID: 2, Name: "Child1.1", Data: 4},
		{ID: 5, PID: 3, Name: "Child2.1", Data: 5},
	}

	// 创建树
	tree := NewTree(data)

	// 测试获取根节点
	rootNodes := tree.GetRootNodes()
	if len(rootNodes) != 1 || rootNodes[0].ID != 1 {
		t.Fatalf("期望有 1 个根节点，ID 为 1，实际为 %v", rootNodes)
	}

	// 测试获取子节点 ID
	subCategoryIDs := tree.GetSubCategoryIDs(1, true)
	if len(subCategoryIDs) != 5 {
		t.Fatalf("期望有 5 个子节点 ID，实际为 %v", subCategoryIDs)
	}

	// 测试获取节点层级
	level := tree.TreeLevel(4)
	if level != 2 {
		t.Fatalf("期望节点 4 的层级为 2，实际为 %v", level)
	}

	// 测试判断父子关系
	if !tree.IsParent(2, 1) {
		t.Fatalf("期望节点 2 是节点 1 的子节点")
	}

	// 测试并发查找子树
	if !tree.IsInSubTreeConcurrent(1, 4) {
		t.Fatalf("期望节点 4 在节点 1 的子树中")
	}
}

// 性能基准测试 IsInSubTreeConcurrent
func BenchmarkIsInSubTreeConcurrent(b *testing.B) {
	// 创建测试数据
	data := []TreeNode[int]{
		{ID: 1, PID: 0, Name: "Root", Data: 1},
		{ID: 2, PID: 1, Name: "Child1", Data: 2},
		{ID: 3, PID: 1, Name: "Child2", Data: 3},
		{ID: 4, PID: 2, Name: "Child1.1", Data: 4},
		{ID: 5, PID: 3, Name: "Child2.1", Data: 5},
	}
	tree := NewTree(data)

	// 重置计时器
	b.ResetTimer()
	// 性能测试
	for i := 0; i < b.N; i++ {
		tree.IsInSubTreeConcurrent(1, 4)
	}
}
