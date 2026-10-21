package auth

import "strings"

func GetAccessKeyFromUrl(url string, accessPrefix string) string {
	items := strings.Split(url, "/")

	if accessPrefix != "" {
		items = append([]string{accessPrefix}, items[2:]...)
	} else {
		items = items[2:]
	}

	key := strings.Join(items, "_")
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, "-", "_")

	return key
}
