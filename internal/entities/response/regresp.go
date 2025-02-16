package response

import "fmt"

type OkResp struct {
	Anwser string `json:"anwser"`
}

// ErrorResp model info
// @Description Error response information.
// @Description with Error description and description what caused this.
// @Description if Error was anwser always be "error".
type ErrorResp struct {
	Anwer string `json:"anwser"`
	Error string `json:"error"`
	Cause string `json:"cause"`
}

func NewOkResp(anwser string) OkResp {
	return OkResp{
		Anwser: anwser,
	}
}
func (r OkResp) JsonString() string {
	okBody := `{
    "anwser": "%s"
}`
	return fmt.Sprintf(okBody, r.Anwser)
}
func NewErrorResp(cause string, err error) ErrorResp {
	return ErrorResp{
		Error: err.Error(),
		Cause: cause,
	}
}
func (r ErrorResp) JsonString() string {
	errorBody := `{
    "anwser": "error",
    "cause": "%s",
    "error": "%s"
}`
	return fmt.Sprintf(errorBody, r.Cause, r.Error)
}
