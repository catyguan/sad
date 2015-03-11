package constv

const (
	STATUS_CONTINUE   = 100
	STATUS_DONE       = 200
	STATUS_ASYNC      = 202
	STATUS_OK         = 204
	STATUS_REDIRECT   = 302
	STATUS_INVALID    = 400
	STATUS_REJECT     = 403
	STATUS_TIMEOUT    = 408
	STATUS_ERROR      = 500
	STATUS_BADGATEWAY = 502
)

const (
	TYPES_NULL   = 0
	TYPES_BOOL   = 1
	TYPES_INT    = 2
	TYPES_LONG   = 3
	TYPES_FLOAT  = 4
	TYPES_DOUBLE = 5
	TYPES_STRING = 6
	TYPES_VAR    = 7
	TYPES_ARRAY  = 8
	TYPES_MAP    = 9
	TYPES_BINARY = 10
)

const (
	KEY_TIMEOUT        = "Timeout"
	KEY_DEADLINE       = "Deadline"
	KEY_APP_ID         = "AppId"
	KEY_REQ_ID         = "ReqId"
	KEY_TRANSACTION_ID = "TransactionId"
	KEY_SESSION_ID     = "SessionId"
	KEY_ASYNC_MODE     = "AsyncMode"
	KEY_ASYNC_ID       = "AsyncId"
	KEY_CALLBACK       = "Callback"
)
