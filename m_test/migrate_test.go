package m_test

import (
	"github.com/d3v-friends/go-pure/fnMatch"
	"github.com/d3v-friends/go-pure/fnPanic"
	"github.com/d3v-friends/mango/m_migrate"
	"testing"
)

type Index struct {
	Key  map[string]int
	Name string
}

func TestMigrate(test *testing.T) {
	var tool = NewTestTool(true)

	test.Run("migrate", func(t *testing.T) {
		var ctx = tool.Context()
		fnPanic.On(m_migrate.Migrate(
			ctx,
			tool.DB,
			&DocTest{},
		))

		var cursor = fnPanic.Get(tool.DB.Collection(docTestNm).Indexes().List(ctx))

		var nameLs = []string{
			"_id_",
			"inTx_1",
		}

		var count = 0
		for cursor.Next(ctx) {
			var idx = &Index{}
			fnPanic.On(cursor.Decode(idx))

			if !fnMatch.Contain(nameLs, idx.Name) {
				t.Fatal("not index")
			}

			count += 1
		}

		if len(nameLs) != count {
			t.Fatal("not same index count")
		}
	})
}
