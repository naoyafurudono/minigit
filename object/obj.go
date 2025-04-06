package object

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
)

// gitのオブジェクトストレージへの読み書きを提供する.
type Object struct {
	data []byte
}

// dataを保持する Object を作成する.
func NewObject(data []byte) Object {
	return Object{data}
}

// オブジェクトの名前. dataの内容に対して一意であると期待できる.
func (o Object) Name() [sha1.Size]byte {
	return sha1.Sum(o.data)
}

// オブジェクトのデータ.
func (o Object) Data() []byte { return o.data }

// Object o を root に保存する.
func (o Object) Store(root string) error {
	n := o.Name()
	d := path.Join(root, ".git", "objects", fmt.Sprintf("%x", n[:1]))
	if err := os.MkdirAll(d, 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	p := path.Join(d, fmt.Sprintf("%x", n[1:]))
	if err := os.WriteFile(p, o.compress(), 0644); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

// rootに保存された名前 name を持つオブジェクトを読み出す.
func ReadObject(root string, name [sha1.Size]byte) (*Object, error) {
	h := name[:]
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
	if err != nil {
		return nil, err
	}
	o := NewObject(data)
	n := o.Name()
	if hex.EncodeToString(n[:]) != hex.EncodeToString(h){
		return nil, errors.New("fatal: the name of object is invalid")
	}
	return &o, nil
}

func (r Object) compress() []byte {
	var buf bytes.Buffer
	w, err := zlib.NewWriterLevel(&buf, 1)
	if err != nil {
		// 指定するレベルがまずいときにだけエラーになる。テストで担保するのでpanicで良い。
		panic(err)
	}
	w.Write(r.data)
	w.Close()
	return buf.Bytes()
}
