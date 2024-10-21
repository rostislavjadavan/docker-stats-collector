package webserver

import (
	"dsc/database"
	"dsc/ui"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func CreateHomepageHandler(db *database.Database) func(c echo.Context) error {
	return func(c echo.Context) error {
		images, err := db.GetImages()
		if err != nil {
			log.Error().Err(err).Msg("Unable to get images")
			return err
		}

		var components []templ.Component
		for _, image := range images {
			memoryStats, err := db.GetTopMemoryStats(image)
			if err != nil {
				log.Error().Err(err).Msg("Unable to get top memory stats")
				return err
			}

			cpuStats, err := db.GetTopCpuStats(image)
			if err != nil {
				log.Error().Err(err).Msg("Unable to get top cpu stats")
				return err
			}

			stats := ui.DailyTopStats(memoryStats, cpuStats)
			components = append(components, ui.ContainerBox(image, stats))
		}

		return RenderView(c, ui.Layout("docker stats", components))
	}
}

//var (
//	itemCntLine = 6
//	fruits      = []string{"Apple", "Banana", "Peach ", "Lemon", "Pear", "Cherry"}
//)
//
//func generateLineItems() []opts.LineData {
//	items := make([]opts.LineData, 0)
//	for i := 0; i < itemCntLine; i++ {
//		items = append(items, opts.LineData{Value: rand.Intn(300)})
//	}
//	return items
//}
//
//func lineSmooth() *charts.Line {
//	line := charts.NewLine()
//	line.SetGlobalOptions(
//		charts.WithTitleOpts(opts.Title{
//			Title: "smooth style",
//		}),
//	)
//
//	line.SetXAxis(fruits).AddSeries("Category A", generateLineItems()).
//		SetSeriesOptions(charts.WithLineChartOpts(
//			opts.LineChart{
//				Smooth: opts.Bool(true),
//			}),
//		)
//	return line
//}

//func CreateHomepageHandler(db *database.Database) func(c echo.Context) error {
//	return func(c echo.Context) error {
//		line := lineSmooth()
//
//		content := fmt.Sprintf(`
//<html>
//<head>
//	<meta charset="utf-8">
//	<title>DSC</title>
//	<script type="text/javascript" src="/public/echarts.min.js"></script>
//</head>
//<body>
//	%s
//</body>
//</html>
//`, lib.RenderChart(line))
//
//		return c.HTML(200, content)
//	}
// }
