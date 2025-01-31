package state

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/helmfile/helmfile/pkg/filesystem"
	"github.com/helmfile/helmfile/pkg/runtime"
	"github.com/helmfile/helmfile/pkg/tmpl"
)

func TestExecuteTemplateExpressions(t *testing.T) {
	render := tmpl.NewFileRenderer(filesystem.DefaultFileSystem(), "", map[string]any{
		"Values": map[string]any{
			"foo": map[string]any{
				"releaseName": "foo",
			},
		},
		"Release": map[string]any{
			"Name": "foo",
		},
	})

	v := runtime.GoccyGoYaml
	runtime.GoccyGoYaml = true
	t.Cleanup(func() {
		runtime.GoccyGoYaml = v
	})

	rs := ReleaseSpec{
		Name:      "foo",
		Chart:     "bar",
		Namespace: "baz",
		ValuesTemplate: []any{
			map[string]any{
				"fullnameOverride": "{{ .Values | get (printf \"%s.releaseName\" .Release.Name) .Release.Name }}",
			},
		},
	}
	result, err := rs.ExecuteTemplateExpressions(render)

	require.NoErrorf(t, err, "failed to execute template expressions: %v", err)
	require.Equalf(t, result.ValuesTemplate[0].(map[string]any)["fullnameOverride"], "foo", "failed to execute template expressions")
}
