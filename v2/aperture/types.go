package aperture

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Responce struct {
	Data  any    `json:"data"`
	Error *Error `json:"error,omitempty"`
}

type DocOutput struct {
	Url         string   `json:"url"`
	Version     string   `json:"version"`
	Alias       string   `json:"alias"`
	Method      string   `json:"method"`
	Input       any      `json:"inputType"`
	Output      any      `json:"outputType"`
	Pathprops   []string `json:"pathProps"`
	Exceptions  []string `json:"exceptions"`
	Description string   `json:"description"`
}

type DocResult struct {
	Schema  any `json:"schema"`
	Version int `json:"version"`
}
