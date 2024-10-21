package lib

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
)

func RenderChart(base *charts.Line) string {
	snippet := base.RenderSnippet()
	return fmt.Sprintf(`
	%s
	%s
`, snippet.Element, snippet.Script)
}
