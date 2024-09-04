package core

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"
)

const (
	StatusInvalidContent             = int(90)  // Content is not available, is nil, or is of the wrong type, usually found via unmarshalling
	StatusIOError                    = int(91)  // I/O operation failed
	StatusJsonDecodeError            = int(92)  // Json decoding failed
	StatusJsonEncodeError            = int(93)  // Json encoding failed
	StatusContentEncodingError       = int(94)  // Content encoding error
	StatusContentEncodingInvalidType = int(95)  // Content encoding error
	StatusNotProvided                = int(96)  // No status is available
	StatusRateLimited                = int(97)  // Rate limited
	StatusNotStarted                 = int(98)  // Not started
	StatusHaveContent                = int(99)  // Content is available
	StatusGzipEncodingError          = int(100) // Gzip encoding error
	StatusGzipDecodingError          = int(101) // Gzip decoding error
	StatusTxnBeginError              = int(102) // Transaction processing begin error
	StatusTxnRollbackError           = int(103) // Transaction processing rollback error
	StatusTxnCommitError             = int(104) // Transaction processing commit error
	StatusExecError                  = int(105) // Execution error, as in a database call
	StatusNoChange                   = int(106) // State has not changed

	/*
		StatusOK                 = codes.OK                 // Not an error; returned on success.
		StatusCancelled          = codes.Canceled           // The operation was cancelled, typically by the caller.
		StatusUnknown            = codes.Unknown            // Unknown error. For example, this error may be returned when a Status value received from another address space belongs to an error space that is not known in this address space. Also errors raised by APIs that do not return enough error information may be converted to this error.
	*/

	StatusInvalidArgument  = 3 //codes.InvalidArgument    // The client specified an invalid argument. Note that this differs from FAILED_PRECONDITION. INVALID_ARGUMENT indicates arguments that are problematic regardless of the state of the system (e.g., a malformed file name).
	StatusDeadlineExceeded = 4 //codes.DeadlineExceeded   // The deadline expired before the operation could complete. For operations that change the state of the system, this error may be returned even if the operation has completed successfully. For example, a successful response from a server could have been delayed long

	/*	StatusNotFound           = codes.NotFound           // Some requested entity (e.g., file or directory) was not found. Note to server developers: if a request is denied for an entire class of users, such as gradual feature rollout or undocumented allowlist, NOT_FOUND may be used. If a request is denied for some users within a class of users, such as user-based access control, PERMISSION_DENIED must be used.
		StatusAlreadyExists      = codes.AlreadyExists      // The entity that a client attempted to create (e.g., file or directory) already exists.
		StatusPermissionDenied   = codes.PermissionDenied   // The caller does not have permission to execute the specified operation. PERMISSION_DENIED must not be used for rejections caused by exhausting some startup (use RESOURCE_EXHAUSTED instead for those errors). PERMISSION_DENIED must not be used if the caller can not be identified (use UNAUTHENTICATED instead for those errors). This error code does not imply the request is valid or the requested entity exists or satisfies other pre-conditions.
		StatusResourceExhausted  = codes.ResourceExhausted  // Some startup has been exhausted, perhaps a per-user quota, or perhaps the entire file system is out of space.
		StatusFailedPrecondition = codes.FailedPrecondition // The operation was rejected because the system is not in a state required for the operation's execution. For example, the directory to be deleted is non-empty, an rmdir operation is applied to a non-directory, etc. Service implementors can use the following guidelines to decide between FAILED_PRECONDITION, ABORTED, and UNAVAILABLE: (a) Use UNAVAILABLE if the client can retry just the failing call. (b) Use ABORTED if the client should retry at a higher level (e.g., when a client-specified test-and-set fails, indicating the client should restart a read-modify-write sequence). (c) Use FAILED_PRECONDITION if the client should not retry until the system state has been explicitly fixed. E.g., if an "rmdir" fails because the directory is non-empty, FAILED_PRECONDITION should be returned since the client should not retry unless the files are deleted from the directory.
		StatusAborted            = codes.Aborted            // The operation was aborted, typically due to a concurrency issue such as a sequencer check failure or transaction abort. See the guidelines above for deciding between FAILED_PRECONDITION, ABORTED, and UNAVAILABLE.
		StatusOutOfRange         = codes.OutOfRange         // The operation was attempted past the valid range. E.g., seeking or reading past end-of-file. Unlike INVALID_ARGUMENT, this error indicates a problem that may be fixed if the system state changes. For example, a 32-bit file system will generate INVALID_ARGUMENT if asked to read at an offset that is not in the range [0,2^32-1], but it will generate OUT_OF_RANGE if asked to read from an offset past the current file size. There is a fair bit of overlap between FAILED_PRECONDITION and OUT_OF_RANGE. We recommend using OUT_OF_RANGE (the more specific error) when it applies so that callers who are iterating through a space can easily look for an OUT_OF_RANGE error to detect when they are done.
		StatusUnimplemented      = codes.Unimplemented      // The operation is not implemented or is not supported/enabled in this service.
		StatusInternal           = codes.Internal           // Internal errors. This means that some invariants expected by the underlying system have been broken. This error code is reserved for serious errors.
		StatusUnavailable        = codes.Unavailable        // The service is currently unavailable. This is most likely a transient condition, which can be corrected by retrying with a backoff. Note that it is not always safe to retry non-idempotent operations.
		StatusDataLoss           = codes.DataLoss           // Unrecoverable data loss or corruption.
		StatusUnauthenticated    = codes.Unauthenticated    // The request does not have valid authentication credentials for the operation.
		_maxGRPCCode             = StatusUnauthenticated
	*/
)

var okStatus = func() *Status {
	s := new(Status)
	s.Code = http.StatusOK
	return s
}()

var badRequestStatus = func() *Status {
	s := new(Status)
	s.Code = http.StatusBadRequest
	return s
}()

var notFoundStatus = func() *Status {
	s := new(Status)
	s.Code = http.StatusNotFound
	return s
}()

type Status struct {
	Code     int   `json:"code"`
	Err      error `json:"err"`
	Handled  bool  `json:"handled"`
	Duration time.Duration
	Content  any
	trace    []string `json:"location"`
}

func StatusOK() *Status {
	return okStatus
}

func StatusBadRequest() *Status {
	return badRequestStatus
}

func StatusNotFound() *Status {
	return notFoundStatus
}

func NewStatus(code int) *Status {
	s := new(Status)
	s.Code = code
	return s
}

func NewStatusError(code int, err error) *Status {
	s := new(Status)
	s.Code = code
	s.Err = err
	s.addTrace(getLocation(2))
	return s
}

func NewStatusDuration(code int, duration time.Duration) *Status {
	s := NewStatus(code)
	s.Duration = duration
	return s
}

func (s *Status) OK() bool {
	return s.Code == http.StatusOK
}

func (s *Status) Timeout() bool {
	return s.Code == http.StatusGatewayTimeout || s.Code == StatusDeadlineExceeded
}

func (s *Status) NotFound() bool {
	return s.Code == http.StatusNotFound
}

func (s *Status) StatusCode() int {
	return s.Code
}

func (s *Status) HttpCode() int {
	return HttpCode(s.Code)
}

func (s *Status) Trace() []string {
	return s.trace
}

func (s *Status) AddLocation() *Status {
	if s.OK() {
		return s
	}
	s.addTrace(getLocation(2))
	return s
}

func (s *Status) addParentLocation() *Status {
	if s.OK() {
		return s
	}
	s.addTrace(getLocation(3))
	return s
}

func (s *Status) addTrace(loc string) *Status {
	if len(s.trace) > 0 && loc == s.trace[len(s.trace)-1] {
		return s
	}
	s.trace = append(s.trace, loc)
	return s
}

func (s *Status) String() string {
	if s.Err != nil {
		return fmt.Sprintf("%v [%v]", HttpStatus(s.Code), s.Err)
	} else {
		return fmt.Sprintf("%v", HttpStatus(s.Code))
	}
}

func getLocation(skip int) string {
	if pc, _, _, ok := runtime.Caller(skip); ok {
		if details := runtime.FuncForPC(pc); details != nil {
			return runtimeNameToUri(details.Name())
		}
	}
	return ""
}

func runtimeNameToUri(uri string) string {
	//uri := details.Name()
	i := strings.Index(uri, githubDotCom)
	if i == -1 {
		return uri
	}
	uri = strings.Replace(uri, githubDotCom, githubHost, len(githubHost))
	// Check for generic function
	i = strings.LastIndex(uri, "[")
	if i != -1 {
		i = strings.LastIndex(uri[:i-1], ".")
	} else {
		i = strings.LastIndex(uri, ".")
	}
	if i == -1 {
		return uri
	}
	uri = uri[:i] + ":" + uri[i+1:]
	return uri
}

// HttpCode - conversion of a code to HTTP status code
func HttpCode(code int) int {
	// Catch all valid http status codes
	if code >= http.StatusContinue {
		return code
	}
	// map known
	switch code {
	case StatusInvalidArgument:
		return http.StatusInternalServerError
	case StatusDeadlineExceeded:
		return http.StatusGatewayTimeout
	}
	// all others
	return http.StatusInternalServerError
}

// HttpStatus - string representation of status code
func HttpStatus(code int) string {
	switch code {
	// Mapped
	case StatusInvalidContent:
		return "Invalid Content"
	case StatusIOError:
		return "I/O Failure"
	case StatusJsonEncodeError:
		return "Json Encode Failure"
	case StatusJsonDecodeError:
		return "Json Decode Failure"
	case StatusContentEncodingError:
		return "Content Decoding Failure"
	case StatusContentEncodingInvalidType:
		return "Invalid Content Encoding Type"
	case StatusGzipDecodingError:
		return "gzip Decoding Failure"
	case StatusGzipEncodingError:
		return "gzip Encoding Failure"
	case StatusNotProvided:
		return "Not Provided"
	case StatusRateLimited:
		return "Rate Limited"
	case StatusNotStarted:
		return "Not Started"
	case StatusDeadlineExceeded:
		return "Deadline Exceeded"
	case StatusInvalidArgument:
		return "Invalid Argument"
	case StatusHaveContent:
		return "Content Available"

		//case StatusUnavailable:
		//	return "Invalid Argument"

		//Http
	case http.StatusOK:
		return "OK"
	case http.StatusAccepted:
		return "Accepted"
	case http.StatusBadRequest:
		return "Bad Request"
	case http.StatusTeapot:
		return "I'm A Teapot"
	case http.StatusGatewayTimeout:
		return "Timeout"
	case http.StatusNotFound:
		return "Not Found"
	case http.StatusMethodNotAllowed:
		return "Method Not Allowed"
	case http.StatusForbidden:
		return "Permission Denied"
	case http.StatusInternalServerError:
		return "Internal Error"
	case http.StatusServiceUnavailable:
		return "Service Unavailable"
	case http.StatusUnauthorized:
		return "Unauthorized"

		// Unmapped
		/*
			case StatusCancelled:
				return "The operation was cancelled, typically by the caller"
			case StatusUnknown:
				return "Unknown error" // For example, this error may be returned when a Status value received from another address space belongs to an error space that is not known in this address space. Also errors raised by APIs that do not return enough error information may be converted to this error."
			case StatusAlreadyExists:
				return "The entity that a client attempted to create already exists"
			case StatusResourceExhausted:
				return "Some startup has been exhausted" //perhaps a per-user quota, or perhaps the entire file system is out of space."
			case StatusFailedPrecondition:
				return "The operation was rejected because the system is not in a state required for the operation's execution" //For example, the directory to be deleted is non-empty, an rmdir operation is applied to a non-directory, etc. Service implementors can use the following guidelines to decide between FAILED_PRECONDITION, ABORTED, and UNAVAILABLE: (a) Use UNAVAILABLE if the client can retry just the failing call. (b) Use ABORTED if the client should retry at a higher level (e.g., when a client-specified test-and-set fails, indicating the client should restart a read-modify-write sequence). (c) Use FAILED_PRECONDITION if the client should not retry until the system state has been explicitly fixed. E.g., if an "rmdir" fails because the directory is non-empty, FAILED_PRECONDITION should be returned since the client should not retry unless the files are deleted from the directory.
			case StatusAborted:
				return "The operation was aborted" // typically due to a concurrency issue such as a sequencer check failure or transaction abort. See the guidelines above for deciding between FAILED_PRECONDITION, ABORTED, and UNAVAILABLE."
			case StatusOutOfRange:
				return "The operation was attempted past the valid range" // E.g., seeking or reading past end-of-file. Unlike INVALID_ARGUMENT, this error indicates a problem that may be fixed if the system state changes. For example, a 32-bit file system will generate INVALID_ARGUMENT if asked to read at an offset that is not in the range [0,2^32-1], but it will generate OUT_OF_RANGE if asked to read from an offset past the current file size. There is a fair bit of overlap between FAILED_PRECONDITION and OUT_OF_RANGE. We recommend using OUT_OF_RANGE (the more specific error) when it applies so that callers who are iterating through a space can easily look for an OUT_OF_RANGE error to detect when they are done."
			case StatusUnimplemented:
				return "The operation is not implemented or is not supported/enabled in this service"
			case StatusDataLoss:
				return "Unrecoverable data loss or corruption"
		*/
	}
	return fmt.Sprintf("error: code not mapped: %v", code)

}
