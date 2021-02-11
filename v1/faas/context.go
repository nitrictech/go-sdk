package faas

// SourceType - enum of the possible sources for a Nitric request.
type SourceType string

const (
	// Request - HTTP Request Source Type
	Request SourceType = "REQUEST"
	// Subscription - Topic Subscription Source Type
	Subscription = "SUBSCRIPTION"
	// Unknown - Unknown Source Types, used when the source can't be determined.
	Unknown = "UNKNOWN"
)

// Each of the source type consts above must be included here.
var sourceTypes []SourceType = []SourceType{
	Request,
	Subscription,
	Unknown,
}

// sourceTypeFromString - converts a string, typically the x-nitric-source-type header, into a SourceType
func sourceTypeFromString(s string) SourceType {
	for _, t := range sourceTypes {
		if s == string(t) {
			return t
		}
	}
	// Default to Unknown if the source type isn't one that has been defined.
	return Unknown
}

func (p SourceType) String() string {
	x := p
	for _, v := range sourceTypes {
		if v == x {
			return string(x)
		}
	}
	return Unknown // This will only happen if manually changed.
}

// NitricContext - Represents the contextual metadata for a Nitric function request.
type NitricContext struct {
	requestID   string
	source      string
	sourceType  SourceType
	payloadType string
}

// GetRequestID - return the request id of the request.
func (c *NitricContext) GetRequestID() string {
	return c.requestID
}

// GetSource - return the source of the request.
func (c *NitricContext) GetSource() string {
	return c.source
}

// GetSourceType - return the source type of the request
func (c *NitricContext) GetSourceType() SourceType {
	return c.sourceType
}

// GetPayloadType - return the payload type of the request payload. Typically a typehint.
func (c *NitricContext) GetPayloadType() string {
	return c.payloadType
}
