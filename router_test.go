package tiny

import (
	"testing"
)

func Test_Router_AddAndFind(t *testing.T) {
	root := createRoot()

	node1 := root.addUrl("/path1")
	node2 := root.addUrl("/users/:id")
	node3 := root.addUrl("/users/:id/blogs")
	node4 := root.addUrl("/users/:id/blogs/:blogId")
	node5 := root.addUrl("/blogs/:id")
	node6 := root.addUrl("/")

	found, node, data := root.findUrl("/path1")
	testItem(t, found, true, "find /path1")
	testItem(t, node, node1, "find /path1")

	found, node, data = root.findUrl("/users/123")
	testItem(t, found, true, "find /users/123")
	testItem(t, data["id"], "123", "find /users/123")
	testItem(t, node, node2, "find /users/123")

	found, node, data = root.findUrl("/users/123/blogs")
	testItem(t, found, true, "find /users/123/blogs")
	testItem(t, data["id"], "123", "find /users/123/blogs")
	testItem(t, node, node3, "find /users/123/blogs")

	found, node, data = root.findUrl("/users/123/blogs/000")
	testItem(t, found, true, "find /users/123/blogs/000")
	testItem(t, data["id"], "123", "find /users/123/blogs/000")
	testItem(t, data["blogId"], "000", "find /users/123/blogs/000")
	testItem(t, node, node4, "find /users/123/blogs/000")

	found, node, data = root.findUrl("/blogs/123")
	testItem(t, found, true, "find /blogs/123")
	testItem(t, data["id"], "123", "find /blogs/123")
	testItem(t, node, node5, "find /blogs/123")

	found, _, data = root.findUrl("/nopath")
	testItem(t, found, false, "find /nopath")

	found, node, data = root.findUrl("/")
	testItem(t, found, true, "find /")
	testItem(t, node, node6, "find /")
}

func Test_Router_Method(t *testing.T) {
	handle := func(ctx *Context) {}
	oneHandle := []Handler{handle}
	twoHandles := []Handler{handle, handle}
	root := createRoot()

	testNode := root.addUrl("/test")
	uidNode := root.addUrl("/users/:id")

	testNode.get(oneHandle)
	handlers := testNode.getHandlers("GET")
	testItem(t, len(handlers), 1, "testNode get")

	uidNode.post(oneHandle)
	handlers = uidNode.getHandlers("POST")
	testItem(t, len(handlers), 1, "uidNode post")

	uidNode.put(twoHandles)
	handlers = uidNode.getHandlers("PUT")
	testItem(t, len(handlers), 2, "uidNode put")

	uidNode.delete(oneHandle)
	handlers = uidNode.getHandlers("DELETE")
	testItem(t, len(handlers), 1, "uidNode delete")

	uidNode.patch(oneHandle)
	handlers = uidNode.getHandlers("PATCH")
	testItem(t, len(handlers), 1, "uidNode patch")

	uidNode.head(oneHandle)
	handlers = uidNode.getHandlers("HEAD")
	testItem(t, len(handlers), 1, "uidNode head")

	uidNode.options(oneHandle)
	handlers = uidNode.getHandlers("OPTIONS")
	testItem(t, len(handlers), 1, "uidNode options")

	uidNode.all(oneHandle)
	handlers = uidNode.getHandlers("ALL")
	testItem(t, len(handlers), 1, "uidNode all")
}
