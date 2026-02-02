package generate

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"
)

type FileRoute struct {
	Import  string
	Package string
	Route   string
	Test    string
	Type    string
	Url     string
}

func getRoutesFrom(path string, routes *[]FileRoute, readOnlyVersions bool) {
	content, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Не удалось прочитать ", path)
		return
	}

	for _, entry := range content {
		if readOnlyVersions && !entry.IsDir() {
			continue
		}

		if readOnlyVersions {
			match, _ := regexp.MatchString(`^v\d*$`, entry.Name())

			if !match {
				log.Fatalf("первая папка (%s) должна быть версией быть в формате 'v[D]'", entry.Name())
			}

		}

		full := path + "/" + entry.Name()
		if entry.IsDir() {
			getRoutesFrom(full, routes, false)
			continue
		}

		if !strings.HasSuffix(full, ".go") {
			continue
		}

		name, _ := strings.CutSuffix(entry.Name(), ".go")

		exp := strings.Split(path, "/")
		routePackage := exp[len(exp)-1]
		routeMethod := strings.ToUpper(name[:1]) + name[1:]
		// routeTest := routeMethod + "Test"
		routeInputType := routeMethod + "Input"
		routeUrl := getPrettyPath(path + "/" + name)

		*routes = append(*routes, FileRoute{
			Import:  path,
			Route:   routeMethod,
			Package: routePackage,
			// Test:    routeTest,
			Type: routeInputType,
			Url:  routeUrl,
		})
	}
}

func getPrettyPath(path string) string {
	var builder strings.Builder
	pt := false

	for _, ch := range path {
		if unicode.IsUpper(ch) {
			if !pt {
				builder.WriteRune('-')
			}
			builder.WriteRune(unicode.ToLower(ch))

			continue
		}

		pt = ch == '/'

		builder.WriteRune(ch)
	}

	return builder.String()
}
