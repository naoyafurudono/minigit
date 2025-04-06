package minigit

// TODO: tree objectを実装する
// tree objectのapiはどんなものが良いだろうか
// New, toSlice, Encodeくらいか

type ndkind int

// treeかblob
type node interface {
}

type tree struct {
	// name to node
	children map[string]node
}

func NewTree( /* object のスライスを受け取ると良さそうか */ ) *tree {
	return nil
}

func (t *tree) toSlice() /* object のスライスを変えす？ */ {
	return
}

func (t *tree) Encode() []byte {
	return nil
}
