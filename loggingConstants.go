package logging

const (
	LogLevel                  ContextKey = "logLevel"
	RequestID                 ContextKey = "requestId"
	TransactionID             ContextKey = "transactionId"
	CommitInfo                ContextKey = "commitInfo"
	ContextValues             ContextKey = "values"
	invalidArgument                      = "Invalid Argument %s: %v"
	systemError                          = "A system error has occurred"
	marshalingError                      = "Marshaling Error"
	dataNotFound                         = "%v Not Found"
	partialResponse                      = "Partial Response"
	responseError                        = "Response Error: %v"
	requestError                         = "Request Error: %v"
	constructHTTPRequestError            = "Failed to construct the HTTP request"
	httpRequestFailed                    = "Http Request Failed"
	missingConfig                        = "Unable to retrieve the config value"
	configError                          = "Unable to load config :%v\n"
	fullRequest                          = "Full Request Completed"
	childRequest                         = "%v Request Completed"
	maxStackLength                       = 50
	skipLevelThree                       = 3
)
