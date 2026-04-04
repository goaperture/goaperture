package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

func ShowTree() {
	if !isApiRoutes() {
		fmt.Println("ERROR - not found api/routes")
		return
	}

	var paths []string
	walk("api", 0, &paths)
}

func walk(root string, margin int, paths *[]string) {
	entrys, err := os.ReadDir(root)

	if err != nil {
		return
	}

	for _, entry := range entrys {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		full := filepath.Join(root, name)

		isVersion := folderIsVersion(name)

		var view string

		if isVersion {
			view = fmt.Sprintf("%s%s %s", strings.Repeat(" ", margin), "└", entry.Name())
		} else if margin != 0 {
			// id := len(*paths)
			*paths = append(*paths, full)

			view = fmt.Sprintf("%s%s %s", strings.Repeat(" ", margin), "└", entry.Name())
		} else {
			view = entry.Name()
		}

		fmt.Println(view)
		count := utf8.RuneCountInString(view) - 1
		walk(full, count, paths)
	}
}

func isApiRoutes() bool {
	info, err := os.Stat("api/routes")
	if err == nil {
		return info.IsDir()
	}

	return false
}

func folderIsVersion(name string) bool {
	matched, _ := regexp.MatchString(`^v\d+$`, name)
	return matched
}
