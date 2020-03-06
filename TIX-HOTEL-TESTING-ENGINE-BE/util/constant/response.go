package constant

const (
	SUCCESS                    = "SUCCESS"
	SYSTEM_ERROR               = "SYSTEM_ERROR"
	DUPLICATE_DATA             = "DUPLICATE_DATA"
	DATA_NOT_EXIST             = "DATA_NOT_EXIST"
	BIND_ERROR                 = "BIND_ERROR"
	RUNTIME_ERROR              = "RUNTIME_ERROR"
	DATE_NOT_VALID             = "DATE_NOT_VALID"
	VENDOR_SHUTDOWN            = "VENDOR_SHUTDOWN"
	METHOD_ARGUMENTS_NOT_VALID = "METHOD_ARGUMENTS_NOT_VALID"
	TOO_MANY_REQUEST           = "TOO_MANY_REQUEST"
	UNAUTHORIZE                = "UNAUTHORIZE"
)

var ResponseMessage = make(map[string]string, 0)

func init() {
	ResponseMessage[SUCCESS] = "success"
	ResponseMessage[SYSTEM_ERROR] = "system_error"
	ResponseMessage[DUPLICATE_DATA] = "duplicate_data"
	ResponseMessage[DATA_NOT_EXIST] = "data_not_exist"
	ResponseMessage[BIND_ERROR] = "bind_error"
	ResponseMessage[RUNTIME_ERROR] = "runtime_error"
	ResponseMessage[DATE_NOT_VALID] = "date_not_valid"
	ResponseMessage[VENDOR_SHUTDOWN] = "Vendor is Shutdown"
	ResponseMessage[METHOD_ARGUMENTS_NOT_VALID] = "method_arguments_not_valid"
	ResponseMessage[TOO_MANY_REQUEST] = "too_many_request"
	ResponseMessage[UNAUTHORIZE] = "unauthorize"
}
