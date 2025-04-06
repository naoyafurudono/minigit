package minigit

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
)

type Object interface {
	// objectの名前を取得する
	Name() [sha1.Size]byte
	// objectのデータを取得する
	Data() []byte
	// zlibで圧縮したデータを取得する
	Compress() []byte
	// しかるべき形で永続化する
	Store() error
}

type blob struct {
	content []byte
	root    string
}

var _ Object = &blob{}

func NewBlob(content []byte, root string) *blob {
	return &blob{content: content, root: root}
}

func ReadBlob(root string, name string) (*blob, error) {
	data, err := ReadObject(root, name)
	if err != nil {
		return nil, err
	}
	content, err := parse(data)
	return NewBlob(content, root), nil
}

func parse(data []byte) ([]byte, error) {
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

func (b *blob) Name() [sha1.Size]byte {
	return sha1.Sum(b.Data())
}

func (b *blob) Data() []byte {
	// "blob <length of content>\0<content>"をASCIIエンコードしたバイト列がblobの表現
	l := len(b.content)
	header := append([]byte(fmt.Sprintf("blob %d", l)), []byte{0}...)
	return append(header, b.content...)
}

func (b *blob) Compress() []byte {
	var r bytes.Buffer
	w, err := zlib.NewWriterLevel(&r, 1)
	if err != nil {
		// 指定するレベルがまずいときにだけエラーになる。テストで担保するのでpanicで良い。
		panic(err)
	}
	w.Write(b.Data())
	w.Close()
	return r.Bytes()
}

func (b *blob) Store() error {
	n := b.Name()
	d := path.Join(b.root, ".git", "objects", fmt.Sprintf("%x", n[:1]))
	if err := os.MkdirAll(d, 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	p := path.Join(d, fmt.Sprintf("%x", n[1:]))
	if err := os.WriteFile(p, b.Compress(), 0644); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func ReadObject(root string, name string) ([]byte, error) {
	// オブジェクトの特定
	h, err := hex.DecodeString(name)
	if err != nil {
		return nil, err
	}
	p := path.Join(root, ".git", "objects", hex.EncodeToString(h[:1]), hex.EncodeToString(h[1:]))

	// オブジェクトの取得
	f, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	r, err := zlib.NewReader(bytes.NewReader(f))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	data, err := io.ReadAll(r)
	return data, err
}
