package core

type ResponseEntity struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func ResponseOK() ResponseEntity {
	return ResponseEntity{
		Code:    200,
		Data:    nil,
		Message: "",
	}
}

func ResponseError(code int, message string) ResponseEntity {
	return ResponseEntity{
		Code:    code,
		Data:    nil,
		Message: message,
	}
}
