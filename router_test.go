package tiny

import (
	"testing"
)

func createTestTree() *routeTree {
	root := newTree("root", "root")

	root.addRoute("/public/")
	root.addRoute("/blogs/:id/stars")
	root.addRoute("/users/:id")

	return root
}

func Test_AddRoute(t *testing.T) {
	root := newTree("root", "root")

	tree := root.addRoute("/public/")
	if tree.name != "public" && tree.kind != "path" {
		t.Error("添加 /public/ 时返回错误的树")
	}

	tree = root.addRoute("/blogs/:id/stars")
	if tree.name != "stars" && tree.kind != "path" {
		t.Error("添加 /blogs/:id/stars 时返回错误的树")
	}

	tree = root.addRoute("/blogs/:id")
	if tree.name != "id" && tree.kind != "param" {
		t.Error("添加 /blogs/:id 时返回错误的树")
	}

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

	tree, params, found = root.find("/public/test.png")
	if found == false {
		t.Error("匹配 /public/test.png 失败")
	} else {
		if tree.name != "public" {
			t.Error("匹配 /public/test.png 时匹配到错误的树")
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
