package tiny

import (
	"strings"
)

type Handle func(*Context)

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
	name string
	kind string
	methodHandler
	subNodes []*routeNode
}

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

func (rn *routeNode) findUrl(url string) (found bool, node *routeNode, data map[string]string) {
	paths := strings.Split(url, "/")
	data = make(map[string]string)
	found, node = rn.find(paths[1:], data)
	return
}

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

func (rn *routeNode) addUrl(url string) (node *routeNode) {
	paths := strings.Split(url, "/")
	node = rn.add(paths[1:])
	return
}

func createRoot() *routeNode {
	return &routeNode{"root", "root", newMethodHandler(), []*routeNode{}}
}
