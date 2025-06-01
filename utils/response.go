package utils

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

func CreateResp(Success bool, Message string, data ...string) Response {
	res := Response{
		Success: Success,
		Message: Message,
	}
	if len(data) > 0 {
		res.Data = data[0]
	}
	return res
}
