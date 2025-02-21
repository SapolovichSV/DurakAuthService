package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

type OkResp struct {
	Anwser string `json:"anwser"`
}

// ErrorResp model info
// @Description Error response information.
// @Description with Error description and description what caused this.
// @Description if Error was anwser always be "error".
type ErrorResp struct {
	Anwer string `json:"anwser"`
	//map[Can't read db] = Internal error missing userID
	mapErrorCause map[string]string
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

// TODO refactor naming
func NewErrorResp(mapErrorCause map[string]error) ErrorResp {
	mapErrString := make(map[string]string, len(mapErrorCause))
	for k, v := range mapErrorCause {
		mapErrString[k] = v.Error()
	}
	return ErrorResp{
		mapErrorCause: mapErrString,
	}
}
func (r ErrorResp) JsonString() string {
	// That's just example for you
	// 	errorBodyTemplate := `{
	//     "anwser": "error",
	//     "can't parse": "internal error 12312",
	//     "can't read": "db error"
	// }`
	if len(r.mapErrorCause) == 0 {
		return `{
     "anwser": "error"
}`
	} else {
		var builder strings.Builder
		//In write description says "write always return len(p),nil", so i skip error check
		builder.WriteString("{\n")
		builder.WriteString(fmt.Sprintf("\t \"anwser\": \"%s\",\n", "error"))
		errs := make([]string, 0, len(r.mapErrorCause))
		for k := range r.mapErrorCause {
			errs = append(errs, k)
		}
		builder.WriteString(
			fmt.Sprintf(
				"\t \"%s\": \"%s\"",
				errs[0],
				r.mapErrorCause[errs[0]],
			),
		)
		for i := 1; i < len(errs); i++ {
			builder.WriteString(",\n")
			builder.WriteString(
				fmt.Sprintf(
					"\t \"%s\": \"%s\"",
					errs[i],
					r.mapErrorCause[errs[i]],
				),
			)
		}
		builder.WriteString("\n}")
		return builder.String()
	}
}
func BeatifyValidationErrors(errs validator.ValidationErrors) map[string]error {
	result := make(map[string]error, len(errs))
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			result[fmt.Sprintf("field %s invalid", err.Field())] = fmt.Errorf("field %s is required", err.Field())
		case "email":
			result[fmt.Sprintf("field %s invalid", err.Field())] = fmt.Errorf("field %s is not email", err.Field())
		default:
			result[fmt.Sprintf("field %s invalid", err.Field())] = fmt.Errorf("field %s is not valid", err.Field())
		}
	}
	return result
}
