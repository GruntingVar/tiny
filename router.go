package tiny

import (
	"strings"
)

type Handler func(*Context)

type methodHandler struct {
	handlers map[string][]Handler
}

func newMethodHandler() methodHandler {
	return methodHandler{make(map[string][]Handler)}
}

func (mh methodHandler) getHandlers(method string) []Handler {
	if mh.handlers[method] != nil {
		return mh.handlers[method]
	} else {
		return mh.handlers["ALL"]
	}
}

func (mh methodHandler) addHandlers(method string, handlers []Handler) {
	mh.handlers[method] = append(mh.handlers[method], handlers...)
}

func (mh methodHandler) post(handlers []Handler) {
	mh.addHandlers("POST", handlers)
}

func (mh methodHandler) get(handlers []Handler) {
	mh.addHandlers("GET", handlers)
}

func (mh methodHandler) put(handlers []Handler) {
	mh.addHandlers("PUT", handlers)
}

func (mh methodHandler) patch(handlers []Handler) {
	mh.addHandlers("PATCH", handlers)
}

func (mh methodHandler) delete(handlers []Handler) {
	mh.addHandlers("DELETE", handlers)
}

func (mh methodHandler) head(handlers []Handler) {
	mh.addHandlers("HEAD", handlers)
}

func (mh methodHandler) options(handlers []Handler) {
	mh.addHandlers("OPTIONS", handlers)
}

func (mh methodHandler) all(handlers []Handler) {
	mh.addHandlers("ALL", handlers)
}

// 路由节点，存储路径名、类型、处理方法和子节点
type routeNode struct {
	name          string       // 节点名称
	kind          string       // 节点类型，可以是path, param
	methodHandler              // 处理该节点的handler
	subNodes      []*routeNode // 子节点
}

// 递归查找
func (rn *routeNode) find(paths []string, data map[string]string) (found bool, node *routeNode) {
	node = &routeNode{}
	found = false
	for _, subNode := range rn.subNodes {

		switch subNode.kind {
		case "path":
			if paths[0] == subNode.name {
				found = true
			}
		case "param":
			found = true
			data[subNode.name] = strings.Replace(paths[0], ":", "", 1)
		}

		if found == true {
			if len(paths) == 1 {
				node = subNode
				return
			} else {
				found, node = subNode.find(paths[1:], data)
				return
			}
		}
	}
	return
}

// 从该节点中寻找url对应的节点，url应该以"/"开头，例如"/test"、"/users/1"，该节点看作是根节点，从子节点中寻找匹配的节点。
// 如果返回的found为true，则node为匹配到的节点，data为匹配过程中获取到的参数
func (rn *routeNode) findUrl(url string) (found bool, node *routeNode, data map[string]string) {
	paths := strings.Split(url, "/")
	data = make(map[string]string)
	found, node = rn.find(paths[1:], data)
	return
}

// 递归创建
func (rn *routeNode) add(paths []string) (node *routeNode) {
	found := false

	var name string
	var kind string
	if strings.HasPrefix(paths[0], ":") {
		name = strings.Replace(paths[0], ":", "", 1)
		kind = "param"
	} else {
		name = paths[0]
		kind = "path"
	}

	for _, subNode := range rn.subNodes {
		if subNode.name == name {
			found = true
			node = subNode
			break
		}
	}
	if found == false {
		node = &routeNode{name, kind, newMethodHandler(), []*routeNode{}}
		rn.subNodes = append(rn.subNodes, node)
	}
	if len(paths) > 1 {
		node = node.add(paths[1:])
	}
	return
}

// 在该节点下添加并返回与url对应的节点，如果节点已存在，则直接返回现有的节点
func (rn *routeNode) addUrl(url string) (node *routeNode) {
	paths := strings.Split(url, "/")
	node = rn.add(paths[1:])
	return
}

// 创建根节点，只是为了让逻辑更清晰的虚拟节点，没有实际的作用
func createRoot() *routeNode {
	return &routeNode{"root", "root", newMethodHandler(), []*routeNode{}}
}
