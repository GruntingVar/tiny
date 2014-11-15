package tiny

import (
	"testing"
)

func createTree(kind string, name string) *routeTree {
	return &routeTree{
		kind:     kind,
		name:     name,
		subTrees: []*routeTree{},
	}
}

func createTestTree() *routeTree {
	root := createTree("root", "root")
	blogs := createTree("path", "blogs")
	blogId := createTree("param", "id")
	stars := createTree("path", "stars")
	users := createTree("path", "users")
	userId := createTree("param", "id")
	public := createTree("path", "public")
	dir := createTree("dir", "dir")

	root.subTrees = append(root.subTrees, blogs)
	root.subTrees = append(root.subTrees, users)
	root.subTrees = append(root.subTrees, public)

	blogs.subTrees = append(blogs.subTrees, blogId)
	blogId.subTrees = append(blogId.subTrees, stars)

	users.subTrees = append(users.subTrees, userId)

	public.subTrees = append(public.subTrees, dir)

	return root
}

func Test_Find(t *testing.T) {
	root := createTestTree()

	tree, params, found := root.find("/public/")
	if found == false {
		t.Error("匹配 /public/ 失败")
	} else {
		if tree.name != "public" {
			t.Error("匹配 /public/ 时匹配到错误的树")
		}
	}

	tree, params, found = root.find("/blogs")
	if found == false {
		t.Error("匹配 /blogs 失败")
	} else {
		if tree.name != "blogs" {
			t.Error("匹配 /blogs 时匹配到错误的树")
		}
	}

	tree, params, found = root.find("/blogs/123abc/stars")
	if found == false {
		t.Error("匹配 /blogs/123abc/stars 失败")
	} else {
		if tree.name != "stars" {
			t.Error("匹配 /blogs/123abc/stars 时匹配到错误的树")
		}
		if params["id"] != "123abc" {
			t.Error("匹配 /blogs/123abc/stars 时没能取到正确的参数")
		}
	}

	tree, params, found = root.find("/users/abc123")
	if found == false {
		t.Error("匹配 /users/abc123 失败")
	} else {
		if tree.name != "id" {
			t.Error("匹配 /users/abc123 时匹配到错误的树")
		}
		if params["id"] != "abc123" {
			t.Error("匹配 /users/abc123 时没能取到正确的参数")
		}
	}

	tree, params, found = root.find("/noexists")
	if found == true {
		t.Error("不应该匹配到 /noexists")
	}

}
