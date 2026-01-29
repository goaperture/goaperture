package responce

type SuccessType struct {
	Success bool `json:"success"`
}

type Responce struct {
	Data     any    `json:"data,omitempty"`
	Redirect string `json:"redirect,omitempty"`
	Error    any    `json:"error,omitempty"`
}

func Success(status bool) Responce {
	return Create(SuccessType{Success: status})
}

func Create(data any) Responce {
	return Responce{Data: data}
}

func Redirect(to string) Responce {
	return Responce{Redirect: to}
}

func Error(error any) Responce {
	return Responce{Error: error}
}
