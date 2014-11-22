package tiny

import (
	"testing"
)

func Test_AddAndFind(t *testing.T) {
	root := createRoot()

	node1 := root.addUrl("/path1")
	node2 := root.addUrl("/users/:id")
	node3 := root.addUrl("/users/:id/blogs")
	node4 := root.addUrl("/users/:id/blogs/:blogId")
	node5 := root.addUrl("/blogs/:id")
	node6 := root.addUrl("/")

	found, node, data := root.findUrl("/path1")
	itemExpect(t, found, true, "find /path1")
	itemExpect(t, node, node1, "find /path1")

	found, node, data = root.findUrl("/users/123")
	itemExpect(t, found, true, "find /users/123")
	itemExpect(t, data["id"], "123", "find /users/123")
	itemExpect(t, node, node2, "find /users/123")

	found, node, data = root.findUrl("/users/123/blogs")
	itemExpect(t, found, true, "find /users/123/blogs")
	itemExpect(t, data["id"], "123", "find /users/123/blogs")
	itemExpect(t, node, node3, "find /users/123/blogs")

	found, node, data = root.findUrl("/users/123/blogs/000")
	itemExpect(t, found, true, "find /users/123/blogs/000")
	itemExpect(t, data["id"], "123", "find /users/123/blogs/000")
	itemExpect(t, data["blogId"], "000", "find /users/123/blogs/000")
	itemExpect(t, node, node4, "find /users/123/blogs/000")

	found, node, data = root.findUrl("/blogs/123")
	itemExpect(t, found, true, "find /blogs/123")
	itemExpect(t, data["id"], "123", "find /blogs/123")
	itemExpect(t, node, node5, "find /blogs/123")

	found, _, data = root.findUrl("/nopath")
	itemExpect(t, found, false, "find /nopath")

	found, node, data = root.findUrl("/")
	itemExpect(t, found, true, "find /")
	itemExpect(t, node, node6, "find /")
}

func Test_Method(t *testing.T) {
	handle := func(ctx *Context) {}
	oneHandle := []Handle{handle}
	twoHandles := []Handle{handle, handle}
	root := createRoot()

	testNode := root.addUrl("/test")
	uidNode := root.addUrl("/users/:id")

	testNode.Get(oneHandle)
	handles := testNode.getHandles("GET")
	itemExpect(t, len(handles), 1, "testNode Get")

	uidNode.Post(oneHandle)
	handles = uidNode.getHandles("POST")
	itemExpect(t, len(handles), 1, "uidNode Post")

	uidNode.Put(twoHandles)
	handles = uidNode.getHandles("PUT")
	itemExpect(t, len(handles), 2, "uidNode Put")

	uidNode.Delete(oneHandle)
	handles = uidNode.getHandles("DELETE")
	itemExpect(t, len(handles), 1, "uidNode Delete")

	uidNode.Patch(oneHandle)
	handles = uidNode.getHandles("PATCH")
	itemExpect(t, len(handles), 1, "uidNode Patch")

	uidNode.Head(oneHandle)
	handles = uidNode.getHandles("HEAD")
	itemExpect(t, len(handles), 1, "uidNode Head")

	uidNode.Options(oneHandle)
	handles = uidNode.getHandles("OPTIONS")
	itemExpect(t, len(handles), 1, "uidNode Options")

	uidNode.All(oneHandle)
	handles = uidNode.getHandles("ALL")
	itemExpect(t, len(handles), 1, "uidNode All")
}
