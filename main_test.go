package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkSQL(b *testing.B) {
	b.Run("regex 100", regexScenario(100))
	b.Run("regex 1000", regexScenario(1000))
	b.Run("regex 10000", regexScenario(10000))

	b.Run("sql search 100", sqlScenario(100))
	b.Run("sql search 1000", sqlScenario(1000))
	b.Run("sql search 10000", sqlScenario(10000))

	b.Run("bluge search 100", blugeSearchScenario(100))
	b.Run("bluge search 1000", blugeSearchScenario(1000))
	b.Run("bluge search 10000", blugeSearchScenario(10000))
}

func BenchmarkIndex(b *testing.B) {
	b.Run("bluge index 100", blugeIndexScenario(100))
	b.Run("bluge index 1000", blugeIndexScenario(1000))
	b.Run("bluge index 10000", blugeIndexScenario(10000))

	b.Run("sql db create 100", sqlDatabaseCreationScenario(100))
	b.Run("sql db create 1000", sqlDatabaseCreationScenario(1000))
	b.Run("sql db create 10000", sqlDatabaseCreationScenario(10000))
}

func Test(t *testing.T) {
	dir, err := generateFiles(100)
	require.NoError(t, err)

	reader, err := indexDir(dir)
	require.NoError(t, err)
	matches, err := blugeSearch(reader, "99")
	require.NoError(t, err)

	require.Len(t, matches, 1)

	reader.Close()
	err = os.RemoveAll(dir)
	require.NoError(t, err)
}

func TestRegexSearch(t *testing.T) {
	dir, err := generateFiles(100)
	require.NoError(t, err)

	matches := regexSearch(dir, "99")

	require.Len(t, matches, 1)

	err = os.RemoveAll(dir)
	require.NoError(t, err)
}

func TestSQLSearch(t *testing.T) {
	dir, err := generateFiles(100)
	require.NoError(t, err)

	engine, err := createDatabase(dir)
	require.NoError(t, err)

	matches, err := sqlSearch(engine, "99")
	require.NoError(t, err)

	require.Len(t, matches, 1)

	err = os.RemoveAll(dir)
	require.NoError(t, err)
}

func sqlDatabaseCreationScenario(count int) func(b *testing.B) {
	return func(b *testing.B) {
		dir, err := generateFiles(count)
		if err != nil {
			b.Error(err)
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = createDatabase(dir)
		}

		err = os.RemoveAll(dir)
		if err != nil {
			b.Error(err)
		}
	}
}

func sqlScenario(count int) func(b *testing.B) {
	return func(b *testing.B) {
		dir, err := generateFiles(count)
		if err != nil {
			b.Error(err)
		}
		engine, err := createDatabase(dir)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			require.NoError(b, err)
			_, err = sqlSearch(engine, "99")
		}

		err = os.RemoveAll(dir)
		if err != nil {
			b.Error(err)
		}
	}
}

func regexScenario(count int) func(b *testing.B) {
	return func(b *testing.B) {
		dir, err := generateFiles(count)
		if err != nil {
			b.Error(err)
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = regexSearch(dir, "99")
		}

		err = os.RemoveAll(dir)
		if err != nil {
			b.Error(err)
		}
	}
}

func blugeSearchScenario(count int) func(b *testing.B) {
	return func(b *testing.B) {
		dir, err := generateFiles(count)
		if err != nil {
			b.Error(err)
		}

		b.ResetTimer()
		b.StopTimer()
		reader, err := indexDir(dir)
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			_, _ = blugeSearch(reader, "99")
		}

		reader.Close()
		err = os.RemoveAll(dir)
		if err != nil {
			b.Error(err)
		}
	}
}

func blugeIndexScenario(count int) func(b *testing.B) {
	return func(b *testing.B) {
		dir, err := generateFiles(count)
		if err != nil {
			b.Error(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			reader, err := indexDir(dir)
			if err != nil {
				b.Error(err)
			}
			reader.Close()
		}

		err = os.RemoveAll(dir)
		if err != nil {
			b.Error(err)
		}
	}
}

func generateFiles(count int) (string, error) {
	dir, err := os.MkdirTemp("", "search-test-*")
	if err != nil {
		return "", err
	}

	for i := 0; i < count; i++ {
		path := filepath.Join(dir, fmt.Sprintf("%d.json", i))
		content := fmt.Sprintf(
			`{"annotations":{"list":[{"builtIn":1,"datasource":{"type":"datasource","uid":"grafana"},"enable":true,"hide":true,"iconColor":"rgba(0, 211, 255, 1)","name":"Annotations & Alerts","target":{"limit":100,"matchAny":false,"tags":[],"type":"dashboard"},"type":"dashboard"}]},"editable":true,"fiscalYearStartMonth":0,"graphTooltip":0,"id":25,"links":[],"liveNow":false,"panels":[{"datasource":{"type":"grafana-iot-twinmaker-datasource","uid":"ecqJU43nk"},"fieldConfig":{"defaults":{"color":{"mode":"thresholds"},"custom":{"align":"auto","displayMode":"auto","inspect":false},"mappings":[],"thresholds":{"mode":"absolute","steps":[{"color":"green","value":null}]}},"overrides":[]},"gridPos":{"h":7,"w":11,"x":0,"y":0},"id":2,"options":{"footer":{"fields":"","reducer":["sum"],"show":false},"showHeader":true},"pluginVersion":"9.1.0-pre","targets":[{"datasource":{"type":"grafana-iot-twinmaker-datasource","uid":"ecqJU43nk"},"queryType":"GetAlarms","refId":"A"}],"title":"Alarm List","transformations":[{"id":"twinmaker-register-links","options":{"addSelectionField":true,"title":"Selected Alarm","vars":[{"fieldName":"entityId","name":"${sel_entity}"},{"fieldName":"alarmName","name":"${sel_comp}"}]}}],"type":"table"},{"datasource":{"type":"grafana-iot-twinmaker-datasource","uid":"ecqJU43nk"},"fieldConfig":{"defaults":{"color":{"mode":"thresholds"},"custom":{"fillOpacity":70,"lineWidth":0,"spanNulls":false},"mappings":[{"options":{"ACKNOWLEDGED":{"color":"blue","index":2,"text":"ACKNOWLEDGED"},"ACTIVE":{"color":"red","index":0,"text":"ACTIVE"},"NORMAL":{"color":"green","index":1,"text":"NORMAL"},"SNOOZE_DISABLED":{"color":"yellow","index":3,"text":"SNOOZE_DISABLED"}},"type":"value"}],"thresholds":{"mode":"absolute","steps":[{"color":"green","value":null}]}},"overrides":[]},"gridPos":{"h":7,"w":13,"x":11,"y":0},"id":4,"options":{"alignValue":"left","legend":{"displayMode":"list","placement":"bottom","showLegend":true},"mergeValues":true,"rowHeight":0.9,"showValue":"auto","tooltip":{"mode":"single","sort":"none"}},"targets":[{"componentName":"${sel_comp}","datasource":{"type":"grafana-iot-twinmaker-datasource","uid":"ecqJU43nk"},"entityId":"${sel_entity}","properties":["alarm_status"],"queryType":"EntityHistory","refId":"A"}],"title":"Selected Alarm History","type":"state-timeline"},{"datasource":{"type":"grafana-iot-twinmaker-datasource","uid":"ecqJU43nk"},"gridPos":{"h":12,"w":24,"x":0,"y":7},"id":6,"options":{"customSelCompVarName":"${sel_comp}","customSelEntityVarName":"${sel_entity}","datasource":"AWS IoT TwinMaker-1","sceneId":""},"targets":[{"componentTypeId":"","datasource":{"type":"grafana-iot-twinmaker-datasource","uid":"ecqJU43nk"},"filter":[{"name":"alarm_status","op":"=","value":"ACTIVE"}],"order":"DESCENDING","properties":["alarm_status"],"queryType":"ComponentHistory","refId":"A"}],"title":"Scene Viewer","type":"grafana-iot-twinmaker-sceneviewer-panel"}],"schemaVersion":37,"style":"dark","tags":[],"templating":{"list":[{"current":{"selected":false,"text":"","value":""},"hide":2,"name":"sel_entity","options":[{"selected":true,"text":"","value":""}],"query":"","skipUrlSync":false,"type":"textbox"},{"current":{"selected":false,"text":"","value":""},"hide":2,"name":"sel_comp","options":[{"selected":true,"text":"","value":""}],"query":"","skipUrlSync":false,"type":"textbox"}]},"time":{"from":"now-6h","to":"now"},"timepicker":{},"timezone":"","title":"%d","uid":"alarm","version":3,"weekStart":""}`,
			i,
		)
		if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
			return "", err
		}
	}

	return dir, nil
}
