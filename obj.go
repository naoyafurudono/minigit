package minigit

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"os"
	"path"
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
}

var _ Object = &blob{}

func NewBlob(content []byte) *blob {
	return &blob{content: content}
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
	root, ok := os.LookupEnv("ROOT")
	if !ok {
		return fmt.Errorf("ROOT is not set")
	}
	d := path.Join(root, ".git", "objects", fmt.Sprintf("%x", n[:1]))
	if err := os.MkdirAll(d, 0755); err != nil {
		return err
	}
	p := path.Join(d, fmt.Sprintf("%x", n[1:]))
	if err := os.WriteFile(p, b.Compress(), 0660); err != nil {
		return err
	}
	return nil
}
