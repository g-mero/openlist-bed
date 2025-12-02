package R

import (
	"openlist-bed/pkg/utils"

	"github.com/duke-git/lancet/v2/slice"
)

type RespError struct {
	Msg        string `json:"msg"`
	Code       string `json:"code"`
	StatusCode int    `json:"status_code"`
}

func (e RespError) Error() string {
	return e.Msg
}

func NewRespError(msg string, code string, statusCode ...int) *RespError {
	status, err := utils.GetFirstElementInSlice(statusCode)
	if err != nil {
		status = 500
	}

	return &RespError{msg, code, status}
}

func RespErrWithDetail(err *RespError, detail string) *RespError {
	return &RespError{detail, err.Code, err.StatusCode}
}

func ErrInvalidParam(params ...string) *RespError {
	if len(params) == 0 {
		return ErrParamInvalid
	}
	return RespErrWithDetail(ErrParamInvalid, "Invalid param: "+slice.Join(params, ", "))
}

func ErrWithTelegram(detail string) *RespError {
	return RespErrWithDetail(ErrTelegram, detail)
}

func ErrWithDatabase(detail string) *RespError {
	return RespErrWithDetail(ErrDatabase, detail)
}

var (
	ErrParseBody          = NewRespError("parse body error.", "REQUEST_BODY_PARSE_ERROR")
	ErrParseStruct        = NewRespError("parse struct error.", "RESPONSE_DATA_PARSE_ERROR")
	ErrWrongPassword      = NewRespError("Wrong password.", "USER_INVALID_PWD")
	ErrTokenInvalid       = NewRespError("Invalid token.", "USER_TOKEN_INVALID", 401)
	ErrApiKeyInvalid      = NewRespError("Invalid API key.", "API_KEY_INVALID", 401)
	ErrTokenDeprecated    = NewRespError("Token deprecated.", "USER_TOKEN_DEPRECATED", 401)
	ErrPermissionRequired = NewRespError("Permission required.", "USER_PERMISSION_REQUIRED", 403)
	ErrGenerateToken      = NewRespError("Generate token failed.", "TOKEN_GENERATE_FAILED")
	ErrGenerateUUID       = NewRespError("Generate uuid failed.", "UUID_GENERATE_FAILED")
	ErrParamInvalid       = NewRespError("Param invalid.", "PARAM_INVALID")
	ErrTimeout            = NewRespError("Request timeout.", "REQUEST_TIMEOUT", 408)
	ErrDatabase           = NewRespError("Database error.", "DATABASE_ERROR")
	ErrNotFound           = NewRespError("Not found.", "404_NOT_FOUND", 404)

	/* telegram related */

	ErrTelegram = NewRespError("Telegram related error.", "TELEGRAM_ERROR")
)
