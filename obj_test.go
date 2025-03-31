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
	runGit(t, []string{"add", "-A"})
	runGit(t, []string{"commit", "-m", "initial commit"})
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
			hash := runGit(t, []string{"hash-object", f})
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
}

// gitを実行してその出力をテストのログにプロキシする.
// gitが0以外のステータスコードで終了した場合はテストを失敗させる.
func runGit(t *testing.T, args []string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Fatal(err)
	} else {
		r := string(out)
		return strings.Trim(r, "\n \t")
	}
	return ""
}
