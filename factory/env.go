package factory

import (
	"github.com/peace0phmind/bud/stream"
	"os"
	"strings"
)

func init() {
	NamedSingleton[map[string]string]("env").SetInitFunc(func() *map[string]string {
		result, _ := stream.ToMap[string, string, string](stream.Of(os.Environ()), func(s string) (string, string, error) {
			kv := strings.SplitN(s, "=", 2)
			if len(kv) == 2 {
				return kv[0], kv[1], nil
			}
			return "", "", stream.ErrContinue
		})

		return &result
	})
}
