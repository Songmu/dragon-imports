package dragon

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func out(libChan chan lib, w io.Writer) error {
	libs := map[string]lib{}
	ambiguous := map[string]bool{}
	for lib := range libChan {
		full := lib.path
		key := lib.pkg + "." + lib.object
		if exist, ok := libs[key]; ok {
			if exist.path != full {
				ambiguous[key] = true
			}
		} else {
			libs[key] = lib
		}
	}

	stdlib := map[string]map[string]bool{
		"unsafe": map[string]bool{
			"Alignof":       true,
			"ArbitraryType": true,
			"Offsetof":      true,
			"Pointer":       true,
			"Sizeof":        true,
		},
	}
	for key, lib := range libs {
		if ambiguous[key] {
			continue
		}
		objMap, ok := stdlib[lib.path]
		if !ok {
			objMap = make(map[string]bool)
		}
		objMap[lib.object] = true
		stdlib[lib.path] = objMap
	}

	_, err := fmt.Fprintf(w, `// AUTO-GENERATED BY dragon-imports

package imports

var stdlib = %#v
`, stdlib)

	return err
}

func outPath() string {
	for _, src := range srcDirs() {
		outPath := filepath.Join(src, "golang.org/x/tools/imports")
		if _, err := os.Stat(outPath); err == nil {
			return outPath
		}
	}
	return ""
}
