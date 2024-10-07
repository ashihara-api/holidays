package usecase

import (
	"context"

	"github.com/ashihara-api/core/domain"
	"github.com/ashihara-api/holidays/core/domain/entity"
)

type (
	DateRange struct {
		From string
		To   string
	}

	HolidayFinderInput struct {
		ID int
	}

	HolidayFinderOutput struct {
		Holiday *entity.Holiday
		Lisence string
	}

	HolidaySearcherInput struct {
		Fields    []string
		DateRange *DateRange
		Limit     int
		Offset    int
	}

	HolidaySearcherOutput struct {
		Holidays []*entity.Holiday
		Lisence  string
	}

	HolidayFinder = domain.Usecase[HolidayFinderInput, HolidayFinderOutput]

	HolidaySearcher = domain.Usecase[HolidaySearcherInput, HolidaySearcherOutput]

	HolidayImporter interface {
		Execute(ctx context.Context) (err error)
	}
)
