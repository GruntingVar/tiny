package tiny

import (
	"strings"
)

type Handle func(Context)

type routeTree struct {
	kind     string // "root", "path", "param", "dir"
	name     string
	handles  []Handle
	subTrees []*routeTree
}

func newTree(kind string, name string) *routeTree {
	return &routeTree{
		kind:     kind,
		name:     name,
		handles:  []Handle{},
		subTrees: []*routeTree{},
	}
}

// 深度优先递归查找
func doFind(rt *routeTree, paths []string, params map[string]string) (tree *routeTree, resParams map[string]string, found bool) {
	tree = &routeTree{}
	resParams = params
	found = false
	pathsLen := len(paths)
	if pathsLen > 0 {
		switch rt.kind {
		case "root":
			tree = rt
			found = true
		case "path":
			if rt.name == paths[0] {
				tree = rt
				found = true
			} else {
				return
			}
		case "param":
			resParams[rt.name] = paths[0]
			tree = rt
			found = true
		default:
			return
		}

		if pathsLen > 1 {
			// pathsLen > 1
			var subParams map[string]string
			for _, subTree := range rt.subTrees {
				if pathsLen == 2 && subTree.kind == "dir" {
					tree = rt
					found = true
					return
				}
				tree, subParams, found = doFind(subTree, paths[1:], resParams)
				if found == true {
					resParams = subParams
					return
				}

			}

			// 如果匹配子树失败，重新将found设为false
			found = false
		}
	}
	return
}

// 搜索路径并返回叶子节点
func (rt *routeTree) find(path string) (tree *routeTree, params map[string]string, found bool) {
	paths := strings.Split(path, "/")
	params = make(map[string]string)
	tree, params, found = doFind(rt, paths, params)
	return
}

// 如果不存在路径，则添加。返回叶子节点
func (rt *routeTree) addNode(paths []string) (tree *routeTree) {
	var newrt *routeTree
	if strings.HasPrefix(paths[0], ":") {
		exists := false
		for _, subTree := range rt.subTrees {
			if subTree.kind == "param" && subTree.name == paths[0] {
				newrt = subTree
				exists = true
				break
			}
		}
		if exists == false {
			name := strings.Replace(paths[0], ":", "", 1)
			newrt = newTree("param", name)
			rt.subTrees = append(rt.subTrees, newrt)
		}
	} else if paths[0] == "" {
		exists := false
		for _, subTree := range rt.subTrees {
			if subTree.kind == "dir" {
				newrt = rt
				exists = true
				break
			}
		}
		if exists == false {
			newrt = rt
			rt.subTrees = append(rt.subTrees, newTree("dir", "dir"))
		}
	} else {
		exists := false
		for _, subTree := range rt.subTrees {
			if subTree.kind == "path" && subTree.name == paths[0] {
				newrt = subTree
				exists = true
				break
			}
		}
		if exists == false {
			newrt = newTree("path", paths[0])
			rt.subTrees = append(rt.subTrees, newrt)
		}
	}
	if len(paths) > 1 {
		tree = newrt.addNode(paths[1:])
	} else {
		tree = newrt
	}
	return
}

// 添加路由，path必须以"/"开头，如"/blog"
func (rt *routeTree) addRoute(path string) (tree *routeTree) {
	paths := strings.Split(path, "/")
	tree = rt.addNode(paths[1:])
	return
}
