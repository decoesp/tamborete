package resp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

var (
	ErrInvalidSyntax = errors.New("resp: invalid syntax")
)

type Parser struct {
	reader *bufio.Reader
}

func NewParser(r *bufio.Reader) *Parser {
	return &Parser{reader: r}
}

func (p *Parser) Parse() (interface{}, error) {
	line, err := p.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	if len(line) < 3 {
		return nil, ErrInvalidSyntax
	}

	switch line[0] {
	case '+':
		return string(line[1 : len(line)-2]), nil
	case '-':
		return errors.New(string(line[1 : len(line)-2])), nil
	case ':':
		n, err := strconv.Atoi(string(line[1 : len(line)-2]))
		if err != nil {
			return nil, err
		}
		return n, nil
	case '$':
		length, err := strconv.Atoi(string(line[1 : len(line)-2]))
		if err != nil {
			return nil, err
		}
		if length < 0 {
			return nil, nil
		}

		data := make([]byte, length)
		_, err = io.ReadFull(p.reader, data)
		if err != nil {
			return nil, err
		}

		_, _ = p.reader.Discard(2)
		return string(data), nil
	case '*':
		count, err := strconv.Atoi(string(line[1 : len(line)-2]))
		if err != nil {
			return nil, err
		}
		var items []interface{}
		for i := 0; i < count; i++ {
			item, err := p.Parse()
			if err != nil {
				return nil, err
			}
			items = append(items, item)
		}
		return items, nil
	default:
		return nil, ErrInvalidSyntax
	}
}

func Serialize(value interface{}) []byte {
	var buf bytes.Buffer

	switch v := value.(type) {
	case string:
		buf.WriteString(fmt.Sprintf("+%s\r\n", v))
	case error:
		buf.WriteString(fmt.Sprintf("-%s\r\n", v.Error()))
	case int:
		buf.WriteString(fmt.Sprintf(":%d\r\n", v))
	case []byte:
		buf.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v))
	case nil:
		buf.WriteString("$-1\r\n")
	case []interface{}:
		buf.WriteString(fmt.Sprintf("*%d\r\n", len(v)))
		for _, item := range v {
			buf.Write(Serialize(item))
		}
	default:
		buf.WriteString("-ERR unknown type\r\n")
	}

	return buf.Bytes()
}
