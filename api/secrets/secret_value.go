package secrets

// SecretValue - Interface
type SecretValue interface {
	// SecretVersionRef - Returns a reference to the version of this value
	Ref() SecretVersionRef
	// AsBytes - Returns the secret value as bytes
	AsBytes() []byte
	// AsString - Returns the secret value as a string
	AsString() string
}

type secretValueImpl struct {
	version SecretVersionRef
	val     []byte
}

func (s *secretValueImpl) Ref() SecretVersionRef {
	return s.version
}

func (s *secretValueImpl) AsBytes() []byte {
	return s.val
}

func (s *secretValueImpl) AsString() string {
	return string(s.AsBytes())
}
