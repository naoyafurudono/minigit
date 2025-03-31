package main

import (
	"encoding/hex"
	"log/slog"
	"os"

	minigit "github.com/naoyafurudono/mini-git"
)

func main() {
	{
		h := slog.NewJSONHandler(os.Stdout, nil)
		slog.SetDefault(slog.New(h))
	}
	h := hex.EncodeToString

	src := "hello\n"
	b := minigit.NewBlob([]byte(src))
	name := b.Name()
	slog.Info("result", "src", src, "name", h(name[:]), "data", h(b.Data()), "compress", h(b.Compress()))
	if err := b.Store(); err != nil {
		slog.Error("failed to store", "err", err)
		os.Exit(1)
	}
}
