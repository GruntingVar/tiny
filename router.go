package tiny

import (
	"strings"
)

type Handle func(Context)

type routeTree struct {
	kind       string // "root", "path", "param", "dir"
	name       string
	preHandles []Handle
	handles    []Handle
	endHandles []Handle
	subTrees   []*routeTree
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
