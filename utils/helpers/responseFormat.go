package helpers

const (
	MessageBadRequest           = "Bad Request"
	MessageUnauthorized         = "Unauthorized"
	MessageForbidden            = "Forbidden"
	MessageNotFound             = "Not Found"
	MessageMethodNotAllowed     = "Method Not Allowed"
	MessageRequestTimeout       = "Request Timeout"
	MessageConflict             = "Conflict"
	MessageGone                 = "Gone"
	MessageLengthRequired       = "Length Required"
	MessagePreconditionFailed   = "Precondition Failed"
	MessagePayloadTooLarge      = "Payload Too Large"
	MessageUnsupportedMediaType = "Unsupported Media Type"
	MessageUnprocessableEntity  = "Unprocessable Entity"
	MessageInternalServerError  = "Internal Server Error"
	MessageNotImplemented       = "Not Implemented"
	MessageBadGateway           = "Bad Gateway"
	MessageServiceUnavailable   = "Service Unavailable"
	MessageGatewayTimeout       = "Gateway Timeout"
	MessageSuccess              = "Success"
	MessageError                = "Error"
)

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}
