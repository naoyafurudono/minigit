package minigit

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"github.com/naoyafurudono/minigit/object"
)

type blob struct {
	content []byte
	root    string
}

func (b *blob) Encode() []byte {
	// "blob <length of content>\0<content>"をASCIIエンコードしたバイト列がblobの表現
	l := len(b.content)
	header := append(fmt.Appendf(nil, "blob %d", l), []byte{0}...)
	return append(header, b.content...)
}

func NewBlob(content []byte) *blob {
	return &blob{content: content}
}

// rootに永続化されている name オブジェクトを読み込む.
func ReadBlob(root string, name object.Name) (*blob, error) {
	r, err := object.ReadObject(root, name)
	if err != nil {
		return nil, err
	}
	content, err := parseBlob(r.Data())
	return NewBlob(content), nil
}

// data をパースしてそのコンテンツを返す.
func parseBlob(data []byte) ([]byte, error) {
	bs := bytes.Split(data, []byte{0})
	if len(bs) != 2 {
		return nil, fmt.Errorf("null char must be 1, %#v", data)
	}
	header := bytes.Split(bs[0], []byte{' '})
	content := bs[1]
	if len(header) != 2 {
		return nil, errors.New("header must be 2 fieleds")
	}
	if string(header[0]) != "blob" {
		return nil, errors.New("blob only supported")
	}
	size, err := strconv.Atoi(string(header[1]))
	if err != nil {
		return nil, fmt.Errorf("header `%s` must contain the content length: %w", string(bs[0]), err)
	}

	if len(content) != size {
		return nil, fmt.Errorf("invalid length, header: %d, content: %s", size, string(content))
	}
	return content, nil
}
