package aperture

type Types struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

func (t *Types) GetWithDefault(alias string) Types {

	if t.Input == "" {
		t.Input = alias + "Input"
	}

	if t.Output == "" {
		t.Output = alias + "Output"
	}

	return *t
}
