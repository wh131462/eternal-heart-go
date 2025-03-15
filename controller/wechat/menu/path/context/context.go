package context

import (
	"eh_go/controller/wechat/menu/path/manager"
	"eh_go/controller/wechat/menu/path/node"
	"eh_go/controller/wechat/menu/sessions"
)

// PathContext 菜单的上下文
type PathContext struct {
	UserID     string
	RawCommand string
	Node       *node.PathNode
	Session    *sessions.UserSession
}

func (c *PathContext) GetUserID() string {
	return c.UserID
}

func (c *PathContext) GetRawCommand() string {
	return c.RawCommand
}

func (c *PathContext) GetNode() *node.PathNode {
	return c.Node
}

func (c *PathContext) GetSession() *sessions.UserSession {
	return c.Session
}

func (c *PathContext) GoToMenu(id string) string {
	session := c.GetSession()
	session.UpdateId(id)

	return manager.GetMenuText(id, session.Timestamp)
}

func (c *PathContext) Back() string {
	parent := c.Node.Parent
	if parent == nil || parent.ID == "" {
		return c.GoToHome()
	}
	return c.GoToMenu(parent.ID)
}

func (c *PathContext) GoToHome() string {
	return c.GoToMenu("main")
}
