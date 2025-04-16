/*
MIT License

Copyright (c) 2023 Samuel Berthe

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package stack

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RemoveGoPath makes a path relative to one of the src directories in the
// $GOPATH environment variable. If $GOPATH is empty or the input path is not
// contained within any of the src directories in $GOPATH, the original path
// is returned. If the input path is contained within multiple of the src
// directories in $GOPATH, it is made relative to the longest one of them.
func RemoveGoPath(path string) string {
	dirs := filepath.SplitList(os.Getenv("GOPATH"))
	// sort in decreasing order by length so the longest matching prefix is
	// removed
	sort.Stable(longestFirst(dirs))
	for _, dir := range dirs {
		srcDir := filepath.Join(dir, "src")
		rel, err := filepath.Rel(srcDir, path)
		if err == nil && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return rel
		}
	}
	return path
}

type longestFirst []string

func (l longestFirst) Len() int           { return len(l) }
func (l longestFirst) Less(i, j int) bool { return len(l[i]) > len(l[j]) }
func (l longestFirst) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
