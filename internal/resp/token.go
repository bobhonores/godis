package resp

import "strconv"

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Token struct {
	Typ   string
	Str   string
	num   int
	Bulk  string
	Array []Token
}

func (t Token) Marshal() []byte {
	switch t.Typ {
	case "array":
		return t.marshalArray()
	case "bulk":
		return t.marshalBulk()
	case "string":
		return t.marshalString()
	case "null":
		return t.marshalNull()
	case "error":
		return t.marshalError()
	default:
		return []byte{}
	}
}

func (t Token) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, t.Str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (t Token) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(t.Bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, t.Bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (t Token) marshalArray() []byte {
	len := len(t.Array)

	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < len; i++ {
		bytes = append(bytes, t.Array[i].Marshal()...)
	}

	return bytes
}

func (t Token) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, t.Str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (t Token) marshalNull() []byte {
	return []byte("$-1\r\n")
}
