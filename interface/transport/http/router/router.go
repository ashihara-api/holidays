package router

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ashihara-api/core/domain/errors"
	"github.com/ashihara-api/core/interface/transport/http/router"
	"github.com/go-chi/chi/v5"

	"github.com/ashihara-api/holidays/core/domain/entity"
	"github.com/ashihara-api/holidays/core/domain/usecase"
	"github.com/ashihara-api/holidays/interface/transport/presenter"
)

type (
	// HolidaysRouter ...
	HolidaysRouter struct {
		bloc presenter.HolidayBloc
	}

	Holiday struct {
		Date         string  `json:"date,omitempty"`
		Name         string  `json:"name,omitempty"`
		DifinitionID int     `json:"difinition_id,omitempty"`
		Summary      *string `json:"summary,omitempty"`
		Description  *string `json:"description,omitempty"`
	}

	HolidaySearchRequest struct {
		Fields    []string   `json:"fields"`
		DateRange *DateRange `json:"date"`
		Limit     int        `json:"limit"`
		Offset    int        `json:"offset"`
	}

	HolidayFindResponse struct {
		Holiday *Holiday `json:"holiday"`
		Lisence string   `json:"lisence"`
	}

	HolidaySearchResponse struct {
		Holidays []*Holiday           `json:"holidays"`
		Lisence  string               `json:"lisence"`
		Request  HolidaySearchRequest `json:"request"`
	}

	DateRange struct {
		From string `json:"from,omitempty"`
		To   string `json:"to,omitempty"`
	}
)

func toResponseHoliday(input *entity.Holiday) *Holiday {
	if input == nil {
		return nil
	}
	output := Holiday(*input)
	return &output
}

func fromSearchRequest(r HolidaySearchRequest) usecase.HolidaySearcherInput {
	input := usecase.HolidaySearcherInput{
		Fields: r.Fields,
		Limit:  r.Limit,
		Offset: r.Offset,
	}

	if r.DateRange != nil {
		dr := usecase.DateRange(*r.DateRange)
		input.DateRange = &dr
	}

	return input
}

func NewHolidays(bloc presenter.HolidayBloc) router.HTTPRouter {
	return &HolidaysRouter{
		bloc: bloc,
	}
}

func (rt *HolidaysRouter) Route(mux *chi.Mux) (err error) {
	routes := router.Route{
		Endpoints: []router.EndpointPattern{
			{
				Pattern: "/holidays",
				Endpoints: map[string]router.Endpoint{
					http.MethodGet: {
						Handler: rt.Search,
					},
					http.MethodPost: {
						Handler: rt.Import,
					},
				},
			},
			{
				Pattern: "/holiday/{id}",
				Endpoints: map[string]router.Endpoint{
					http.MethodGet: {
						Handler: rt.Find,
					},
				},
			},
		},
	}
	r := router.New(routes)
	return r.Build(mux)
}

func (rt *HolidaysRouter) Find(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RenderError(w, errors.NewCause(fmt.Errorf("id is invalid"), errors.CaseBadRequest))
		return
	}
	result, err := rt.bloc.Find(r.Context(), usecase.HolidayFinderInput{
		ID: id,
	})
	if err != nil {
		RenderError(w, err)
		return
	}

	hs := toResponseHoliday(result.Holiday)

	RenderJSON(w, http.StatusOK, HolidayFindResponse{
		Holiday: hs,
		Lisence: result.Lisence,
	})
}

func (rt *HolidaysRouter) Import(w http.ResponseWriter, r *http.Request) {
	err := rt.bloc.Import(r.Context())
	if err != nil {
		RenderError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (rt *HolidaysRouter) Search(w http.ResponseWriter, r *http.Request) {
	var rq HolidaySearchRequest
	if err := BindFromJSON(r.Body, &rq); err != nil {
		RenderError(w, errors.NewCause(err, errors.CaseBadRequest))
		return
	}
	output, err := rt.bloc.Search(r.Context(), fromSearchRequest(rq))
	if err != nil {
		RenderError(w, err)
		return
	}

	var hs []*Holiday
	for _, h := range output.Holidays {
		hs = append(hs, toResponseHoliday(h))
	}

	RenderJSON(w, http.StatusOK, HolidaySearchResponse{
		Holidays: hs,
		Lisence:  output.Lisence,
		Request:  rq,
	})
}
