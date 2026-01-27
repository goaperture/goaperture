package auth

import "strings"

func GetAccessKeyFromUrl(url string) string {
	items := strings.Split(url, "/")

	key := strings.Join(items[2:], "_")
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, "-", "_")

	return key
}
