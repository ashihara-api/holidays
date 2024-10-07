package usecase

import (
	"context"
	"log/slog"

	"github.com/ashihara-api/core/domain/errors"
	"github.com/usk81/toolkit/pointer"

	"github.com/ashihara-api/holidays/core/domain/entity"
	"github.com/ashihara-api/holidays/core/domain/repository"
	"github.com/ashihara-api/holidays/core/domain/usecase"
)

type (
	searchHolidayImpl struct {
		searcher repository.Searcher
		logger   *slog.Logger
	}
)

func SearchHoliday(
	searcher repository.Searcher,
	logger *slog.Logger,
) usecase.HolidaySearcher {
	return &searchHolidayImpl{
		searcher: searcher,
		logger:   logger,
	}
}

func (u *searchHolidayImpl) Execute(
	ctx context.Context,
	input usecase.HolidaySearcherInput,
) (
	output *usecase.HolidaySearcherOutput,
	err error,
) {
	options := repository.SearchOption{
		Limit:  input.Limit,
		Offset: input.Offset,
	}
	if input.DateRange != nil {
		dr := &repository.DateRange{}
		if input.DateRange.From == "" {
			dr.Gte = pointer.To(input.DateRange.From)
		}
		if input.DateRange.To == "" {
			dr.Lte = pointer.To(input.DateRange.To)
		}
		if dr.Gte != nil || dr.Lte != nil {
			options.Range = dr
		}
	}
	rs, err := u.searcher.Search(ctx, options)
	if err != nil {
		u.logger.ErrorContext(ctx,
			"SearchHoliday.Execute",
			slog.String("action", "searcher.Search"),
			slog.Any("error", err),
		)
		return nil, errors.NewCause(err, errors.CaseBackendError)
	}
	if len(rs) == 0 {
		return &usecase.HolidaySearcherOutput{
			Holidays: []*entity.Holiday{},
			Lisence:  "CC-BY",
		}, nil
	}
	for i, r := range rs {
		r = dropFields(r, input.Fields)
		rs[i] = r
	}
	return &usecase.HolidaySearcherOutput{
		Holidays: rs,
		Lisence:  "CC-BY",
	}, nil
}
