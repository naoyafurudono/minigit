package minigit

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
)

type Object interface {
	// objectの名前を取得する
	Name() [sha1.Size]byte
	// objectのデータを取得する
	Data() []byte
	// zlibで圧縮したデータを取得する
	Compress() []byte
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
	w := zlib.NewWriter(&r)
	w.Write(b.Data())
	w.Close()
	return r.Bytes()
}
