package responce

type SuccessType struct {
	Success bool `json:"success"`
}

type Responce struct {
	Data     any    `json:"data"`
	Redirect string `json:"redirect,omitempty"`
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
