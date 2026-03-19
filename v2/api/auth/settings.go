package auth

type Settings map[string]map[string]string

func (s *Settings) Get(prefix, key, def string) string {
	if props, ok_prefix := (*s)[prefix]; ok_prefix {
		if val, ok := props[key]; ok {
			return val
		}
	}

	return def
}
