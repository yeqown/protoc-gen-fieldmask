package module

import pgs "github.com/lyft/protoc-gen-star"

type pkgMessageCache struct {
	messages   map[string]pgs.Message
	parsedFile map[string]struct{}
}

func newCache(initCapacity int) *pkgMessageCache {
	if initCapacity <= 0 {
		initCapacity = 4
	}

	return &pkgMessageCache{
		messages:   make(map[string]pgs.Message, initCapacity*2),
		parsedFile: make(map[string]struct{}, initCapacity),
	}
}

func (c *pkgMessageCache) cached(fqn string) (pgs.Message, bool) {
	if c == nil {
		return nil, false
	}

	message, ok := c.messages[fqn]
	return message, ok
}

// cache adds the message to the cache. the key is the fully qualified name of the message.
// .eg. '.foo.bar.Baz' for a message named Baz in a package named foo.bar.
func (c *pkgMessageCache) cache(fqn string, message pgs.Message) {
	if c == nil {
		return
	}

	c.messages[fqn] = message
}

// isFileParsed judge file has been visited, the key is the full name of the file,
// such as "common/common.proto".
func (c *pkgMessageCache) isFileParsed(inputPath string) bool {
	if c == nil {
		return false
	}

	_, ok := c.parsedFile[inputPath]
	return ok
}

func (c *pkgMessageCache) markFileParsed(inputPath string) {
	if c == nil {
		return
	}

	c.parsedFile[inputPath] = struct{}{}
}
