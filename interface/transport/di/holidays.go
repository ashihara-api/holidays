package di

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/ashihara-api/holidays/core/interactor/usecase"
	"github.com/ashihara-api/holidays/interface/datasource/crawler"
	"github.com/ashihara-api/holidays/interface/datasource/mysql"
	"github.com/ashihara-api/holidays/interface/transport/presenter"
)

func Holidays(c *http.Client, wdb, rdb *sql.DB, logger *slog.Logger) presenter.HolidayBloc {
	writer := mysql.HolidaysWriter(wdb, logger)
	reader := mysql.HolidaysReader(rdb, logger)

	crawl := crawler.Holidays(c, logger)

	return presenter.NewHolidayBloc(
		usecase.FindHoliday(reader, logger),
		usecase.ImportHoliday(crawl, writer, logger),
		usecase.SearchHoliday(reader, logger),
	)
}
