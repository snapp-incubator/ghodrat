package logger

type FieldType uint8

const (
	AnyType FieldType = iota
	BoolType
	IntType
	Float64Type
	StringType
	ErrorType
)

type Field struct {
	Key   string
	Value interface{}
	Type  FieldType
}

// Any constructs a field with the given key and value.
func Any(key string, val interface{}) Field {
	return Field{Key: key, Value: val, Type: AnyType}
}

// Bool constructs a field with the given key and value.
func Bool(key string, val bool) Field {
	return Field{Key: key, Value: val, Type: BoolType}
}

// Int constructs a field with the given key and value.
func Int(key string, val int) Field {
	return Field{Key: key, Value: val, Type: IntType}
}

// Float constructs a field with the given key and value.
func Float64(key string, val float64) Field {
	return Field{Key: key, Value: val, Type: Float64Type}
}

// String constructs a field with the given key and value.
func String(key string, val string) Field {
	return Field{Key: key, Value: val, Type: StringType}
}

// Error constructs a field with the given key and value.
func Error(val error) Field {
	return Field{Key: "error", Value: val, Type: ErrorType}
}
