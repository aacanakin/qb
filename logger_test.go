package qb

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestLogger(t *testing.T) {
	engine, err := New("sqlite3", ":memory:")
	metadata := MetaData()
	actors := Table("actors",
		Column("id", BigInt()).NotNull(),
		PrimaryKey("id"),
	)
	metadata.AddTable(actors)
	metadata.CreateAll(engine)
	defer metadata.DropAll(engine)
	logCapture := &TestingLogWriter{t, nil}
	defer logCapture.Flush()
	engine.SetLogger(&DefaultLogger{LQuery | LBindings, log.New(logCapture, "", log.LstdFlags)})
	engine.Logger().SetLogFlags(LQuery)

	_, err = engine.Exec(actors.Insert().Values(map[string]interface{}{"id": 5}))
	assert.Nil(t, err)

	engine.Logger().SetLogFlags(LBindings)
	_, err = engine.Exec(actors.Insert().Values(map[string]interface{}{"id": 10}))
	assert.Nil(t, err)

	assert.Equal(t, engine.Logger().LogFlags(), LQuery|LBindings)
}
