package markdown

import (
	"fmt"
	"testing"
)

func TestFrontmatter(t *testing.T) {
	str := "\n  --- \ntitle: Go FAQ\nmenu:\n  docs:\n    parent: 'extras'\n    weight: 20\n\n--- \n "
	s, e := Frontmatter([]byte(str))
	fmt.Println("'", str[s:e], "'")


	str = "---\ndate: 2020-06-02\nmanagedBy: xy\n---\n\n\nBei u"
	s, e = Frontmatter([]byte(str))
	fmt.Println("'", str[s:e], "'")
}
