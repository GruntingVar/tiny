package tiny

import (
	"strings"
)

type Handle func(*Context)

type matchData map[string]interface{}

type pathMatcher interface {
	match(string, matchData) bool
}

func newPathMatcher(path string) pathMatcher {
	switch {
	case path == "":
		return &dirNode{}
	case path == "**":
		return &anyNode{}
	case strings.HasPrefix(path, ":"):
		name := strings.Replace(path, ":", "", 1)
		return &paramNode{name: name}
	default:
		return &pathNode{name: path}
	}
}

type rootNode struct {
}

func (n *rootNode) match(path string, data matchData) bool {
	return true
}

type pathNode struct {
	name string
}

func (n *pathNode) match(path string, data matchData) bool {
	return n.name == path
}

type paramNode struct {
	name string
}

func (n *paramNode) match(path string, data matchData) bool {
	data[n.name] = path
	return true
}

type dirNode struct {
}

func (n *dirNode) match(path string, data matchData) bool {
	return path == ""
}

type anyNode struct {
}

func (n *anyNode) match(path string, data matchData) bool {
	data["endMatch"] = true
	return true
}

type methodHandler struct {
	handles map[string][]Handle
}

func newMethodHandler() methodHandler {
	return methodHandler{make(map[string][]Handle)}
}

func (mh methodHandler) getHandles(method string) []Handle {
	if mh.handles[method] != nil {
		return mh.handles[method]
	} else {
		return mh.handles["ALL"]
	}
}

func (mh methodHandler) addHandles(method string, handles []Handle) {
	mh.handles[method] = append(mh.handles[method], handles...)
}

func (mh methodHandler) Post(handles []Handle) {
	mh.addHandles("POST", handles)
}

func (mh methodHandler) Get(handles []Handle) {
	mh.addHandles("GET", handles)
}

func (mh methodHandler) Put(handles []Handle) {
	mh.addHandles("PUT", handles)
}

func (mh methodHandler) Patch(handles []Handle) {
	mh.addHandles("PATCH", handles)
}

func (mh methodHandler) Delete(handles []Handle) {
	mh.addHandles("DELETE", handles)
}

func (mh methodHandler) Head(handles []Handle) {
	mh.addHandles("HEAD", handles)
}

func (mh methodHandler) Options(handles []Handle) {
	mh.addHandles("OPTIONS", handles)
}

func (mh methodHandler) All(handles []Handle) {
	mh.addHandles("ALL", handles)
}

type routeNode struct {
	pathMatcher
	methodHandler
	subNodes []*routeNode
}

func (rn *routeNode) find(paths []string, data matchData) (found bool, node *routeNode) {
	node = &routeNode{}
	found = false
	for _, subNode := range rn.subNodes {
		found = subNode.match(paths[0], data)
		if found == true {
			if len(paths) == 1 {
				node = subNode
				return
			} else {
				if data["endMatch"] == true {
					delete(data, "endMatch")
					node = subNode
					return
				} else {
					found, node = subNode.find(paths[1:], data)
					return
				}
			}
		}
	}
	return
}

func (rn *routeNode) findUrl(url string) (found bool, node *routeNode, data matchData) {
	paths := strings.Split(url, "/")
	data = matchData{}
	found, node = rn.find(paths[1:], data)
	return
}

func (rn *routeNode) add(paths []string) (node *routeNode) {
	data := matchData{}
	found := false
	for _, subNode := range rn.subNodes {
		found = subNode.match(paths[0], data)
		if found == true {
			node = subNode
			break
		}
	}
	if found == false {
		matcher := newPathMatcher(paths[0])
		node = createNode(matcher)
		rn.subNodes = append(rn.subNodes, node)
	}
	if len(paths) > 1 {
		node = node.add(paths[1:])
	}
	return
}

func (rn *routeNode) addUrl(url string) (node *routeNode) {
	paths := strings.Split(url, "/")
	node = rn.add(paths[1:])
	return
}

func createRoot() *routeNode {
	return &routeNode{
		pathMatcher:   &rootNode{},
		methodHandler: newMethodHandler(),
		subNodes:      []*routeNode{},
	}
}

func createNode(matcher pathMatcher) *routeNode {
	return &routeNode{
		pathMatcher:   matcher,
		methodHandler: newMethodHandler(),
		subNodes:      []*routeNode{},
	}
}
