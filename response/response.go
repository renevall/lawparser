package response

const (
	StatusSuccess = "success"
	StatusError   = "error"
	StatusFail    = "fail"
)

//Error contains a JSEND Error json response Struct
type Error struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

//Response contains a JSEND json response Struct
type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type jsend interface {
	Wrap()
}

//Wrap is used to wrap the response object around JSEND standard
func (e *Error) Wrap(s string, err string) {
	e.Status = s
	e.Error = err
}

//Wrap is used to wrap the response object around JSEND standard
func (r *Response) Wrap(s string, d interface{}) {
	r.Status = s
	r.Data = d
}
