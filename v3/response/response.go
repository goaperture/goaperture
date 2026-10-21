package response

type SuccessType struct {
	Success bool `json:"success"`
}

type Pagination struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

type Response struct {
	Data       any        `json:"data,omitempty"`
	Redirect   string     `json:"redirect,omitempty"`
	Error      any        `json:"error,omitempty"`
	Pagination Pagination `json:"pagination,omitempty"`
}

func Success(status bool) Response {
	return Create(SuccessType{Success: status})
}

func Create(data any) Response {
	return Response{Data: data}
}

func Redirect(to string) Response {
	return Response{Redirect: to}
}

func Error(error any) Response {
	return Response{Error: error}
}
