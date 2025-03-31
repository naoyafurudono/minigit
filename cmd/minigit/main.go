package main

import (
	"encoding/hex"
	"log/slog"
	"os"

	"github.com/naoyafurudono/minigit"
)

func main() {
	{
		h := slog.NewJSONHandler(os.Stdout, nil)
		slog.SetDefault(slog.New(h))
	}
	h := hex.EncodeToString
	root, ok := os.LookupEnv("ROOT")
	if !ok {
		slog.Error("ROOT is not set")
		os.Exit(1)
	}

	src := "hello\n"
	b := minigit.NewBlob([]byte(src), root)
	name := b.Name()
	slog.Info("result", "src", src, "name", h(name[:]), "data", h(b.Data()), "compress", h(b.Compress()))
	if err := b.Store(); err != nil {
		slog.Error("failed to store", "err", err)
		os.Exit(1)
	}
}
