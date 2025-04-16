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
	"fmt"
	"os"
	"strings"
	"sync"
)

var mutex sync.RWMutex
var cache = map[string][]string{}

const nbrLinesBefore = 5
const nbrLinesAfter = 5

func readFile(path string) ([]string, bool) {
	mutex.RLock()
	lines, ok := cache[path]
	mutex.RUnlock()

	if ok {
		return lines, true
	}

	if !strings.HasSuffix(path, ".go") {
		return nil, false
	}

	// bearer:disable go_gosec_filesystem_filereadtaint
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	lines = strings.Split(string(b), "\n")

	mutex.Lock()
	cache[path] = lines
	mutex.Unlock()

	return lines, true
}

func getSourceFromFrame(frame stacktraceFrame) []string {
	lines, ok := readFile(frame.rawFile)
	if !ok {
		return []string{}
	}

	if len(lines) < frame.line {
		return []string{}
	}

	current := frame.line - 1
	start := max(0, current-nbrLinesBefore)
	end := min(len(lines)-1, current+nbrLinesAfter)

	var output []string

	for i := start; i <= end; i++ {
		if i < 0 || i >= len(lines) {
			continue
		}

		line := lines[i]
		message := fmt.Sprintf("%d\t%s", i+1, line)
		output = append(output, message)

		if i == current {
			lenWithoutLeadingSpaces := len(strings.TrimLeft(line, " \t"))
			lenLeadingSpaces := len(line) - lenWithoutLeadingSpaces
			nbrTabs := strings.Count(line[0:lenLeadingSpaces], "\t")
			firstCharIndex := lenLeadingSpaces + (8-1)*nbrTabs // 8 chars per tab

			sublinePrefix := strings.Repeat(" ", firstCharIndex)
			subline := strings.Repeat("^", lenWithoutLeadingSpaces)
			output = append(output, "\t"+sublinePrefix+subline)
		}
	}

	return output
}
