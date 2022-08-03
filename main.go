package main

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/blugelabs/bluge"
	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/information_schema"
)

type Match struct {
	Count int
	Path  string
}

func regexSearch(dir string, search string) []Match {
	re := regexp.MustCompile(`(?i)(` + search + `)`)
	matches := []Match{}
	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		if matchCount := re.FindAllIndex(b, -1); matchCount != nil {
			matches = append(matches, Match{len(matchCount), path})
		}

		return nil
	}
	filepath.Walk(dir, walk)
	return matches
}

func indexDir(dir string) (*bluge.Reader, error) {
	config := bluge.InMemoryOnlyConfig()
	writer, err := bluge.OpenWriter(config)
	if err != nil {
		return nil, err
	}
	defer writer.Close()

	batch := bluge.NewBatch()

	docsInBatch := 0
	maxBatchSize := 500

	flushIfRequired := func(force bool) error {
		docsInBatch++
		needFlush := force || (maxBatchSize > 0 && docsInBatch >= maxBatchSize)
		if !needFlush {
			return nil
		}
		err := writer.Batch(batch)
		if err != nil {
			return err
		}
		docsInBatch = 0
		batch.Reset()
		return nil
	}

	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		dashboard := struct {
			Title string `json:"title"`
		}{}

		if err := json.Unmarshal(b, &dashboard); err != nil {
			return err
		}

		doc := bluge.NewDocument(path).AddField(bluge.NewTextField("name", dashboard.Title))
		batch.Insert(doc)
		if err := flushIfRequired(false); err != nil {
			return err
		}
		err = writer.Update(doc.ID(), doc)
		if err != nil {
			return err
		}

		return nil
	}
	if err := flushIfRequired(true); err != nil {
		return nil, err
	}

	filepath.Walk(dir, walk)

	return writer.Reader()
}

func blugeSearch(reader *bluge.Reader, search string) ([]Match, error) {
	matches := []Match{}

	query := bluge.NewWildcardQuery("*" + search + "*").SetField("name")
	request := bluge.NewTopNSearch(100, query).
		WithStandardAggregations()
	documentMatchIterator, err := reader.Search(context.Background(), request)
	if err != nil {
		return nil, err
	}
	match, err := documentMatchIterator.Next()
	for err == nil && match != nil {
		// load the identifier for this match
		err = match.VisitStoredFields(func(field string, value []byte) bool {
			if field == "_id" && value != nil {
				matches = append(matches, Match{Count: 1, Path: string(value)})
			}
			return true
		})
		if err != nil {
			return nil, err
		}
		match, err = documentMatchIterator.Next()
	}
	if err != nil {
		return nil, err
	}

	return matches, nil
}

func createDatabase(dir string) (*sqle.Engine, error) {
	const (
		dbName    = "grafana"
		tableName = "dashboards"
	)
	db := memory.NewDatabase(dbName)
	table := memory.NewTable(tableName, sql.NewPrimaryKeySchema(sql.Schema{
		{Name: "name", Type: sql.Text, Nullable: false, Source: tableName},
		{Name: "path", Type: sql.Text, Nullable: false, Source: tableName},
	}), nil)

	db.AddTable(tableName, table)
	ctx := sql.NewEmptyContext()
	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		dashboard := struct {
			Title string `json:"title"`
		}{}

		if err := json.Unmarshal(b, &dashboard); err != nil {
			return err
		}
		table.Insert(ctx, sql.NewRow(dashboard.Title, path))
		return nil
	}

	filepath.Walk(dir, walk)

	return sqle.NewDefault(
		sql.NewDatabaseProvider(
			db,
			information_schema.NewInformationSchemaDatabase(),
		)), nil
}

func sqlSearch(engine *sqle.Engine, search string) ([]Match, error) {
	ctx := sql.NewEmptyContext()
	query := `SELECT name, path FROM grafana.dashboards WHERE name LIKE "%` + search + `%";`
	_, rows, err := engine.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close(ctx)
	}()

	matches := []Match{}

	row, err := rows.Next(ctx)
	for err == nil && row != nil {
		if row[0] != nil && strings.Contains(row[0].(string), search) {
			matches = append(matches, Match{Count: 1, Path: row[1].(string)})
		}
		row, err = rows.Next(ctx)
	}
	if err != nil && err != io.EOF {
		return nil, err
	}
	return matches, nil
}
