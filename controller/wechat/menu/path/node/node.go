package node

import (
	"eh_go/controller/wechat/menu/sessions"
	"errors"
	"fmt"
)

// Context 定义处理器所需的上下文接口
type Context interface {
	GetUserID() string
	GetRawCommand() string
	GetNode() *PathNode
	GetSession() *sessions.UserSession
	GoToMenu(id string) string
	Back() string
	GoToHome() string
}

// OnNodeHandler 当前节点的单独处理器
type OnNodeHandler func(context Context) (string, error)

// PathNode 路径节点
type PathNode struct {
	ID          string //
	Name        string // 菜单名称
	Content     string // 内容
	HandlerFunc OnNodeHandler

	Parent   *PathNode   // 父节点指针
	Children []*PathNode // 子节点集合
}

// NewRoot 创建根节点
func NewRoot(id, name string, content string) *PathNode {
	return &PathNode{
		ID:      id,
		Name:    name,
		Content: content,
	}
}

// AddHandler 添加节点函数处理器
func (n *PathNode) AddHandler(callback OnNodeHandler) {
	n.HandlerFunc = callback
}

// AddChild 添加子节点
func (n *PathNode) AddChild(child *PathNode) {
	child.Parent = n
	n.Children = append(n.Children, child)
}

// Path 查找节点路径（从当前节点到根）
func (n *PathNode) Path() []*PathNode {
	var path []*PathNode
	for current := n; current != nil; current = current.Parent {
		path = append([]*PathNode{current}, path...)
	}
	return path
}

// FindByID 根据ID查找节点
func (n *PathNode) FindByID(id string) *PathNode {
	if n.ID == id {
		return n
	}
	for _, child := range n.Children {
		if found := child.FindByID(id); found != nil {
			return found
		}
	}
	return nil
}

// DFS 深度优先遍历 (DFS)
func (n *PathNode) DFS(fn func(*PathNode)) {
	fn(n)
	for _, child := range n.Children {
		child.DFS(fn)
	}
}

// BFS 广度优先遍历 (BFS)
func (n *PathNode) BFS(fn func(*PathNode)) {
	queue := []*PathNode{n}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		fn(current)
		queue = append(queue, current.Children...)
	}
}

// 判断当前节点是否是目标节点的祖先
func (n *PathNode) isAncestor(target *PathNode) bool {
	if target == nil || n == nil {
		return false
	}

	// 快速路径检查
	if n == target {
		return true // 自己是自己的祖先（根据需求可选）
	}

	// 迭代实现（避免递归栈溢出）
	current := target.Parent
	for current != nil {
		if current == n {
			return true
		}
		current = current.Parent
	}
	return false
}

// 变种方法：判断是否是直接父节点
func (n *PathNode) isParent(child *PathNode) bool {
	return child != nil && child.Parent == n
}

// 从当前节点移除指定子节点（安全版）
func (n *PathNode) removeChild(child *PathNode) error {
	if child == nil || child.Parent != n {
		return fmt.Errorf("node %s is not a child of %s", child.ID, n.ID)
	}

	// 优化查找算法（避免遍历所有子节点）
	found := -1
	for i := 0; i < len(n.Children); i++ {
		if n.Children[i] == child {
			found = i
			break
		}
	}

	if found == -1 {
		return fmt.Errorf("child not found")
	}

	// 高效删除（保持顺序）
	n.Children = append(n.Children[:found], n.Children[found+1:]...)

	// 清空父指针
	child.Parent = nil

	// 维护深度（可选）
	child.updateDepth(-1)
	return nil
}

// 递归更新子树深度
func (n *PathNode) updateDepth(delta int) {
	for _, child := range n.Children {
		child.updateDepth(delta)
	}
}

// MoveNode 移动子树
func MoveNode(source, target *PathNode) error {
	// 检查循环依赖
	if source.isAncestor(target) {
		return errors.New("cannot create circular dependency")
	}

	// 从原父节点移除
	if source.Parent != nil {
		err := source.Parent.removeChild(source)
		if err != nil {
			return err
		}
	}

	// 添加到新父节点
	target.AddChild(source)
	return nil
}

// RemoveChild 删除子树
func (n *PathNode) RemoveChild(child *PathNode) {
	for i, c := range n.Children {
		if c == child {
			n.Children = append(n.Children[:i], n.Children[i+1:]...)
			child.Parent = nil
			return
		}
	}
}

// GetChild 获取对应的子节点
func (n *PathNode) GetChild(idx int) *PathNode {
	if idx < 0 || idx >= len(n.Children) {
		return nil
	}
	return n.Children[idx]
}
