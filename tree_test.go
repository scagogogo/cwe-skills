package cwe

import (
	"testing"
)

func TestNewTreeNode(t *testing.T) {
	cwe := &CWE{ID: 79, Name: "XSS"}
	node := NewTreeNode(cwe)

	if node == nil {
		t.Fatal("expected non-nil node")
	}
	if node.CWE != cwe {
		t.Error("expected CWE to be set")
	}
	if len(node.Children) != 0 {
		t.Errorf("expected empty Children, got %d", len(node.Children))
	}
	if node.Parent != nil {
		t.Error("expected nil Parent")
	}
	if node.Depth != 0 {
		t.Errorf("expected Depth 0, got %d", node.Depth)
	}
}

func TestTreeNode_AddChild(t *testing.T) {
	parent := NewTreeNode(&CWE{ID: 1, Name: "Root"})
	child := NewTreeNode(&CWE{ID: 2, Name: "Child"})

	parent.AddChild(child)

	if len(parent.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(parent.Children))
	}
	if parent.Children[0] != child {
		t.Error("child not added correctly")
	}
	if child.Parent != parent {
		t.Error("child's Parent not set")
	}
	if child.Depth != 1 {
		t.Errorf("expected child Depth 1, got %d", child.Depth)
	}

	// Add grandchild
	grandchild := NewTreeNode(&CWE{ID: 3, Name: "Grandchild"})
	child.AddChild(grandchild)

	if grandchild.Depth != 2 {
		t.Errorf("expected grandchild Depth 2, got %d", grandchild.Depth)
	}
	if grandchild.Parent != child {
		t.Error("grandchild's Parent not set")
	}
}

func TestTreeNode_Walk(t *testing.T) {
	root := NewTreeNode(&CWE{ID: 1, Name: "Root"})
	child1 := NewTreeNode(&CWE{ID: 2, Name: "Child1"})
	child2 := NewTreeNode(&CWE{ID: 3, Name: "Child2"})
	root.AddChild(child1)
	root.AddChild(child2)

	t.Run("DFS traversal", func(t *testing.T) {
		var visited []int
		root.Walk(func(node *TreeNode) bool {
			visited = append(visited, node.CWE.ID)
			return true
		})
		if len(visited) != 3 {
			t.Errorf("expected 3 visited, got %d: %v", len(visited), visited)
		}
		if visited[0] != 1 {
			t.Errorf("expected root first, got %d", visited[0])
		}
	})

	t.Run("early termination", func(t *testing.T) {
		var visited []int
		root.Walk(func(node *TreeNode) bool {
			visited = append(visited, node.CWE.ID)
			return node.CWE.ID != 1 // stop after root
		})
		if len(visited) != 1 {
			t.Errorf("expected 1 visited with early termination, got %d", len(visited))
		}
	})

	t.Run("nil node", func(t *testing.T) {
		var node *TreeNode
		called := false
		node.Walk(func(n *TreeNode) bool {
			called = true
			return true
		})
		if called {
			t.Error("expected Walk on nil node to not call fn")
		}
	})
}

func TestTreeNode_WalkBFS(t *testing.T) {
	root := NewTreeNode(&CWE{ID: 1, Name: "Root"})
	child1 := NewTreeNode(&CWE{ID: 2, Name: "Child1"})
	child2 := NewTreeNode(&CWE{ID: 3, Name: "Child2"})
	grandchild := NewTreeNode(&CWE{ID: 4, Name: "Grandchild"})
	root.AddChild(child1)
	root.AddChild(child2)
	child1.AddChild(grandchild)

	t.Run("BFS traversal", func(t *testing.T) {
		var visited []int
		root.WalkBFS(func(node *TreeNode) bool {
			visited = append(visited, node.CWE.ID)
			return true
		})
		// BFS order: 1, 2, 3, 4
		expected := []int{1, 2, 3, 4}
		if len(visited) != len(expected) {
			t.Fatalf("expected %d visited, got %d: %v", len(expected), len(visited), visited)
		}
		for i, id := range visited {
			if id != expected[i] {
				t.Errorf("visited[%d] = %d, expected %d", i, id, expected[i])
			}
		}
	})

	t.Run("early termination", func(t *testing.T) {
		var visited []int
		root.WalkBFS(func(node *TreeNode) bool {
			visited = append(visited, node.CWE.ID)
			return node.CWE.ID != 1 // stop after root
		})
		if len(visited) != 1 {
			t.Errorf("expected 1 visited with early termination, got %d", len(visited))
		}
	})

	t.Run("nil node", func(t *testing.T) {
		var node *TreeNode
		called := false
		node.WalkBFS(func(n *TreeNode) bool {
			called = true
			return true
		})
		if called {
			t.Error("expected WalkBFS on nil node to not call fn")
		}
	})
}

func TestTreeNode_Find(t *testing.T) {
	root := NewTreeNode(&CWE{ID: 1, Name: "Root"})
	child := NewTreeNode(&CWE{ID: 2, Name: "Child"})
	grandchild := NewTreeNode(&CWE{ID: 3, Name: "Grandchild"})
	root.AddChild(child)
	child.AddChild(grandchild)

	tests := []struct {
		name   string
		id     int
		found  bool
	}{
		{"found at root", 1, true},
		{"found in child", 2, true},
		{"found in grandchild", 3, true},
		{"not found", 999, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := root.Find(tt.id)
			if (result != nil) != tt.found {
				t.Errorf("Find(%d) found=%v, want %v", tt.id, result != nil, tt.found)
			}
		})
	}

	t.Run("nil node", func(t *testing.T) {
		var node *TreeNode
		result := node.Find(1)
		if result != nil {
			t.Error("expected nil for nil node")
		}
	})
}

func TestTreeNode_Path(t *testing.T) {
	root := NewTreeNode(&CWE{ID: 1, Name: "Root"})
	child := NewTreeNode(&CWE{ID: 2, Name: "Child"})
	grandchild := NewTreeNode(&CWE{ID: 3, Name: "Grandchild"})
	root.AddChild(child)
	child.AddChild(grandchild)

	t.Run("from root", func(t *testing.T) {
		path := root.Path()
		if len(path) != 1 {
			t.Fatalf("expected path length 1, got %d", len(path))
		}
		if path[0].CWE.ID != 1 {
			t.Errorf("expected path[0] = 1, got %d", path[0].CWE.ID)
		}
	})

	t.Run("from leaf", func(t *testing.T) {
		path := grandchild.Path()
		if len(path) != 3 {
			t.Fatalf("expected path length 3, got %d", len(path))
		}
		expectedIDs := []int{1, 2, 3}
		for i, node := range path {
			if node.CWE.ID != expectedIDs[i] {
				t.Errorf("path[%d].ID = %d, expected %d", i, node.CWE.ID, expectedIDs[i])
			}
		}
	})

	t.Run("nil node", func(t *testing.T) {
		var node *TreeNode
		path := node.Path()
		if path != nil {
			t.Errorf("expected nil path for nil node, got %v", path)
		}
	})
}

func TestTreeNode_LeafNodes(t *testing.T) {
	root := NewTreeNode(&CWE{ID: 1, Name: "Root"})
	child1 := NewTreeNode(&CWE{ID: 2, Name: "Child1"})
	child2 := NewTreeNode(&CWE{ID: 3, Name: "Child2"})
	root.AddChild(child1)
	root.AddChild(child2)

	t.Run("with leaves", func(t *testing.T) {
		leaves := root.LeafNodes()
		if len(leaves) != 2 {
			t.Errorf("expected 2 leaves, got %d", len(leaves))
		}
	})

	t.Run("no leaves (single node)", func(t *testing.T) {
		single := NewTreeNode(&CWE{ID: 1, Name: "Single"})
		leaves := single.LeafNodes()
		if len(leaves) != 1 {
			t.Errorf("expected 1 leaf (itself), got %d", len(leaves))
		}
	})

	t.Run("nil node", func(t *testing.T) {
		var node *TreeNode
		leaves := node.LeafNodes()
		if leaves != nil {
			t.Errorf("expected nil for nil node, got %v", leaves)
		}
	})
}

func TestTreeNode_MaxDepth(t *testing.T) {
	root := NewTreeNode(&CWE{ID: 1, Name: "Root"})
	child := NewTreeNode(&CWE{ID: 2, Name: "Child"})
	grandchild := NewTreeNode(&CWE{ID: 3, Name: "Grandchild"})
	root.AddChild(child)
	child.AddChild(grandchild)

	if root.MaxDepth() != 2 {
		t.Errorf("expected MaxDepth 2, got %d", root.MaxDepth())
	}

	t.Run("nil node", func(t *testing.T) {
		var node *TreeNode
		if node.MaxDepth() != 0 {
			t.Errorf("expected 0 for nil node, got %d", node.MaxDepth())
		}
	})
}

func TestTreeNode_Count(t *testing.T) {
	root := NewTreeNode(&CWE{ID: 1, Name: "Root"})
	child1 := NewTreeNode(&CWE{ID: 2, Name: "Child1"})
	child2 := NewTreeNode(&CWE{ID: 3, Name: "Child2"})
	root.AddChild(child1)
	root.AddChild(child2)

	if root.Count() != 3 {
		t.Errorf("expected Count 3, got %d", root.Count())
	}

	t.Run("nil node", func(t *testing.T) {
		var node *TreeNode
		if node.Count() != 0 {
			t.Errorf("expected 0 for nil node, got %d", node.Count())
		}
	})
}

func TestTreeNode_IsLeaf(t *testing.T) {
	root := NewTreeNode(&CWE{ID: 1, Name: "Root"})
	child := NewTreeNode(&CWE{ID: 2, Name: "Child"})

	if !root.IsLeaf() {
		t.Error("expected root to be leaf before adding children")
	}

	root.AddChild(child)
	if root.IsLeaf() {
		t.Error("expected root to not be leaf after adding children")
	}
	if !child.IsLeaf() {
		t.Error("expected child to be leaf")
	}
}

func TestTreeNode_IsRoot(t *testing.T) {
	root := NewTreeNode(&CWE{ID: 1, Name: "Root"})
	child := NewTreeNode(&CWE{ID: 2, Name: "Child"})

	if !root.IsRoot() {
		t.Error("expected root to be root")
	}

	root.AddChild(child)
	if child.IsRoot() {
		t.Error("expected child to not be root")
	}
}

func TestTreeNode_String(t *testing.T) {
	t.Run("normal node", func(t *testing.T) {
		node := NewTreeNode(&CWE{ID: 79, Name: "XSS"})
		s := node.String()
		if s == "" {
			t.Error("expected non-empty string")
		}
	})

	t.Run("nil node", func(t *testing.T) {
		var node *TreeNode
		s := node.String()
		if s != "TreeNode<nil>" {
			t.Errorf("expected 'TreeNode<nil>', got %q", s)
		}
	})

	t.Run("nil CWE", func(t *testing.T) {
		node := &TreeNode{}
		s := node.String()
		if s != "TreeNode<nil>" {
			t.Errorf("expected 'TreeNode<nil>', got %q", s)
		}
	})
}

func TestBuildTree(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 1, Name: "Root", Relationships: []Relationship{
		{Nature: RelationshipParentOf, CWEID: 2},
	}})
	r.Register(&CWE{ID: 2, Name: "Child", Relationships: []Relationship{
		{Nature: RelationshipChildOf, CWEID: 1},
	}})
	r.BuildIndexes()

	t.Run("valid root", func(t *testing.T) {
		tree := BuildTree(r, 1)
		if tree == nil {
			t.Fatal("expected non-nil tree")
		}
		if tree.CWE.ID != 1 {
			t.Errorf("expected root ID 1, got %d", tree.CWE.ID)
		}
		if len(tree.Children) != 1 {
			t.Errorf("expected 1 child, got %d", len(tree.Children))
		}
		if tree.Children[0].CWE.ID != 2 {
			t.Errorf("expected child ID 2, got %d", tree.Children[0].CWE.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		tree := BuildTree(r, 999)
		if tree != nil {
			t.Error("expected nil for non-existent root")
		}
	})

	t.Run("nil registry", func(t *testing.T) {
		tree := BuildTree(nil, 1)
		if tree != nil {
			t.Error("expected nil for nil registry")
		}
	})
}

func TestBuildForest(t *testing.T) {
	t.Run("with pillars", func(t *testing.T) {
		r := NewRegistry()
		r.Register(&CWE{ID: 1, Name: "Pillar1", Abstraction: AbstractionPillar, Relationships: []Relationship{
			{Nature: RelationshipParentOf, CWEID: 2},
		}})
		r.Register(&CWE{ID: 2, Name: "Child1", Abstraction: AbstractionBase, Relationships: []Relationship{
			{Nature: RelationshipChildOf, CWEID: 1},
		}})
		r.BuildIndexes()

		forest := BuildForest(r)
		if len(forest) != 1 {
			t.Fatalf("expected 1 tree in forest, got %d", len(forest))
		}
		if forest[0].CWE.ID != 1 {
			t.Errorf("expected root ID 1, got %d", forest[0].CWE.ID)
		}
	})

	t.Run("empty registry", func(t *testing.T) {
		r := NewRegistry()
		forest := BuildForest(r)
		if forest != nil {
			t.Errorf("expected nil for empty registry, got %v", forest)
		}
	})

	t.Run("nil registry", func(t *testing.T) {
		forest := BuildForest(nil)
		if forest != nil {
			t.Error("expected nil for nil registry")
		}
	})
}

func TestBuildViewTree(t *testing.T) {
	t.Run("valid view", func(t *testing.T) {
		r := NewRegistry()
		r.Register(&CWE{ID: 79, Name: "XSS", Relationships: []Relationship{}})
		r.RegisterView(&View{ID: 1000, Name: "Research Concepts", Members: []ViewMember{
			{CWEID: 79, ViewID: 1000, Direct: true},
		}})
		r.BuildIndexes()

		tree := BuildViewTree(r, 1000)
		if tree == nil {
			t.Fatal("expected non-nil tree")
		}
		if tree.CWE.ID != 1000 {
			t.Errorf("expected root ID 1000, got %d", tree.CWE.ID)
		}
		if len(tree.Children) != 1 {
			t.Errorf("expected 1 child, got %d", len(tree.Children))
		}
		if tree.Children[0].CWE.ID != 79 {
			t.Errorf("expected child ID 79, got %d", tree.Children[0].CWE.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		r := NewRegistry()
		tree := BuildViewTree(r, 9999)
		if tree != nil {
			t.Error("expected nil for non-existent view")
		}
	})

	t.Run("nil registry", func(t *testing.T) {
		tree := BuildViewTree(nil, 1000)
		if tree != nil {
			t.Error("expected nil for nil registry")
		}
	})

	t.Run("view with duplicate member ID", func(t *testing.T) {
		r := NewRegistry()
		r.Register(&CWE{ID: 79, Name: "XSS", Relationships: []Relationship{}})
		r.RegisterView(&View{ID: 1000, Name: "View1", Members: []ViewMember{
			{CWEID: 1000, ViewID: 1000, Direct: true}, // Same as view ID - should be skipped
			{CWEID: 79, ViewID: 1000, Direct: true},
		}})
		r.BuildIndexes()

		tree := BuildViewTree(r, 1000)
		if tree == nil {
			t.Fatal("expected non-nil tree")
		}
		// Only CWE-79 should be a child (1000 is skipped because it's the view ID itself)
		if len(tree.Children) != 1 {
			t.Errorf("expected 1 child, got %d", len(tree.Children))
		}
	})
}
