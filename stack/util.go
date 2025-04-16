package stack

import (
	"runtime"
	"strings"
)

func shortFuncName(f *runtime.Func) string {
	// f.Name() is like one of these:
	// - "github.com/example/proj/package.FuncName"
	// - "github.com/example/proj/package.Receiver.MethodName"
	// - "github.com/example/proj/package.(*PtrReceiver).MethodName"
	longName := f.Name()

	withoutPath := longName[strings.LastIndex(longName, "/")+1:]
	withoutPackage := withoutPath[strings.Index(withoutPath, ".")+1:]

	shortName := withoutPackage
	shortName = strings.Replace(shortName, "(", "", 1)
	shortName = strings.Replace(shortName, "*", "", 1)
	shortName = strings.Replace(shortName, ")", "", 1)

	return shortName
}
