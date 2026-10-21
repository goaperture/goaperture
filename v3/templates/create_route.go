package templates

import (
	"os"
	"path/filepath"
)

func CreateRoute(path, name, description string, sequre bool) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	os.MkdirAll(absPath, 0777)
	pkg := filepath.Base(absPath)

	fileContent := GetRouteCode(pkg, name, description, sequre)

	file := filepath.Join(absPath, name+".go")

	os.WriteFile(file, []byte(fileContent), 0777)
}
