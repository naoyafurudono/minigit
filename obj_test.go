package minigit_test

import (
	"encoding/hex"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/naoyafurudono/minigit"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestBlob(t *testing.T) {
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.Chdir(orig)
	})
	datarepo := path.Join(orig, "testdata", "repo")

	temprepo := t.TempDir()
	if err := os.CopyFS(temprepo, os.DirFS(datarepo)); err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(temprepo); err != nil {
		t.Fatal(err)
	}

	runGit(t, []string{"init"})
	entries, err := os.ReadDir(temprepo)
	if err != nil {
		t.Fatal(err)
	}
	fs := lo.FilterMap(entries, func(e fs.DirEntry, i int) (string, bool) {
		if e.IsDir() {
			return "", false
		}
		return path.Join(temprepo, e.Name()), true
	})

	t.Run("Name", func(t *testing.T) {
		for _, f := range fs {
			hash := strings.Trim(
				runGit(t, []string{"hash-object", f}),
				"\n")
			t.Log(f, hash)
			content, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}

			n := minigit.NewBlob([]byte(content), temprepo).Name()
			h := hex.EncodeToString(n[:])
			if hash != h {
				t.Fatalf("expected %s, but got %s", hash, h)
			}
		}
	})

	t.Run("Store - Read", func(t *testing.T) {
		for _, f := range fs {
			content, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			b := minigit.NewBlob([]byte(content), temprepo)
			if err := b.Store(); err != nil {
				t.Fatal(err)
			}
			n := b.Name()
			h := hex.EncodeToString(n[:])
	
			// gitがobjectファイルを正しく読めることを検証
			res := runGit(t, []string{"cat-file", "-p", h})
			if res != string(content) {
				t.Fatalf("expected %s, but got %s", string(content), res)
			}

			a, err := minigit.ReadBlob(temprepo, h)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, b.Name(), a.Name(), "Name()")
			assert.Equal(t, b.Data(), a.Data(), "Data()")
		}
	})
}

// gitを実行して空白などを除去した上で出力を文字列として返す.
// gitが0以外のステータスコードで終了した場合はテストを失敗させる.
func runGit(t *testing.T, args []string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Fatal(err)
	} else {
		return string(out)

	}
	return ""
}
