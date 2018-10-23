package errors

import "strconv"

func (c Code) String() string {
	switch c {
	case OK:
		return "ok"
	case Canceled:
		return "canceled"
	case Unknown:
		return "unknown"
	case InvalidArgument:
		return "invalid_argument"
	case DeadlineExceeded:
		return "deadline_exceeded"
	case NotFound:
		return "not_found"
	case AlreadyExists:
		return "already_exists"
	case PermissionDenied:
		return "permission_denied"
	case ResourceExhausted:
		return "resource_exhausted"
	case FailedPrecondition:
		return "failed_precondition"
	case Aborted:
		return "aborted"
	case OutOfRange:
		return "out_of_range"
	case Unimplemented:
		return "unimplemented"
	case Internal:
		return "internal"
	case Unavailable:
		return "unavailable"
	case DataLoss:
		return "data_loss"
	case Unauthenticated:
		return "unauthenticated"
	default:
		return "code_" + strconv.FormatInt(int64(c), 10)
	}
}