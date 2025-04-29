package types

// MaskedString is a type that fully masks a string when String is called.
//
//nolint:recvcheck
type MaskedString string

// Decode implements the envconfig.Decoder interface.
func (ms *MaskedString) Decode(value string) error {
	*ms = MaskedString(value)

	return nil
}

// String fully masks the underlying string value.
func (ms MaskedString) String() string {
	return "****"
}
