package usecase

import (
	"context"
	"log/slog"

	"github.com/ashihara-api/core/domain/errors"

	"github.com/ashihara-api/holidays/core/domain/repository"
	"github.com/ashihara-api/holidays/core/domain/usecase"
)

type (
	findHolidayImpl struct {
		finder repository.DifinitionFinder
		logger *slog.Logger
	}
)

func FindHoliday(
	finder repository.DifinitionFinder,
	logger *slog.Logger,
) usecase.HolidayFinder {
	return &findHolidayImpl{
		finder: finder,
		logger: logger,
	}
}

func (u *findHolidayImpl) Execute(
	ctx context.Context,
	input usecase.HolidayFinderInput,
) (
	output *usecase.HolidayFinderOutput,
	err error,
) {
	r, err := u.finder.FindDefinition(ctx, input.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			return nil, errors.NewCause(err, errors.CaseNotFound)
		}
		u.logger.ErrorContext(ctx,
			"FindHoliday.Execute",
			slog.String("action", "finder.FindDefinition"),
			slog.Any("error", err),
		)
		return nil, errors.NewCause(err, errors.CaseBackendError)
	}
	return &usecase.HolidayFinderOutput{
		Holiday: r,
		Lisence: "CC-BY",
	}, nil
}
