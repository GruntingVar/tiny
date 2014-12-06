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
	handles := testNode.getHandles("GET")
	testItem(t, len(handles), 1, "testNode get")

	uidNode.post(oneHandle)
	handles = uidNode.getHandles("POST")
	testItem(t, len(handles), 1, "uidNode post")

	uidNode.put(twoHandles)
	handles = uidNode.getHandles("PUT")
	testItem(t, len(handles), 2, "uidNode put")

	uidNode.delete(oneHandle)
	handles = uidNode.getHandles("DELETE")
	testItem(t, len(handles), 1, "uidNode delete")

	uidNode.patch(oneHandle)
	handles = uidNode.getHandles("PATCH")
	testItem(t, len(handles), 1, "uidNode patch")

	uidNode.head(oneHandle)
	handles = uidNode.getHandles("HEAD")
	testItem(t, len(handles), 1, "uidNode head")

	uidNode.options(oneHandle)
	handles = uidNode.getHandles("OPTIONS")
	testItem(t, len(handles), 1, "uidNode options")

	uidNode.all(oneHandle)
	handles = uidNode.getHandles("ALL")
	testItem(t, len(handles), 1, "uidNode all")
}
