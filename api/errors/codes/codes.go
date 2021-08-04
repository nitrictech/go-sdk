package codes

type Code int

const (
	OK                 Code = 0
	Cancelled          Code = 1
	Unknown            Code = 2
	InvalidArgument    Code = 3
	DeadlineExceeded   Code = 4
	NotFound           Code = 5
	AlreadyExists      Code = 6
	PermissionDenied   Code = 7
	ResourceExhausted  Code = 8
	FailedPrecondition Code = 9
	Aborted            Code = 10
	OutOfRange         Code = 11
	Unimplemented      Code = 12
	Internal           Code = 13
	Unavailable        Code = 14
	DataLoss           Code = 15
	Unauthenticated    Code = 16
)

func (c Code) String() string {
	switch c {
	case OK:
		return "OK"
	case Cancelled:
		return "Cancelled"
	case Unknown:
		return "Unknown"
	case InvalidArgument:
		return "Invalid Argument"
	case DeadlineExceeded:
		return "Deadline Exceeded"
	case AlreadyExists:
		return "Already Exists"
	case PermissionDenied:
		return "Permission Denied"
	case ResourceExhausted:
		return "Resource Exhausted"
	case FailedPrecondition:
		return "Failed Precondition"
	case Aborted:
		return "Aborted"
	case OutOfRange:
		return "Out of Range"
	case Unimplemented:
		return "Unimplemented"
	case Internal:
		return "Internal"
	case Unavailable:
		return "Unavailable"
	case DataLoss:
		return "Data Loss"
	case Unauthenticated:
		return "Unauthenticated"
	default:
		return "Unknown"
	}
}
