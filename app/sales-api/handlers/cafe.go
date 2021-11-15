package handlers

import (
	"context"
	"fmt"
	"github.com/Salauatinho/soft-final-project/business/auth"
	"github.com/Salauatinho/soft-final-project/business/data/cafe"
	"github.com/Salauatinho/soft-final-project/business/data/user"
	"github.com/Salauatinho/soft-final-project/foundation/web"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type cafeGroup struct {
	cafe cafe.Cafe
	auth *auth.Auth
}

func (sg cafeGroup) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	params := web.Params(r)
	pageNumber, err := strconv.Atoi(params["page"])
	if err != nil {
		return web.NewRequestError(fmt.Errorf("invalid page format: %s", params["page"]), http.StatusBadRequest)
	}

	rowsPerPage, err := strconv.Atoi(params["rows"])
	if err != nil {
		return web.NewRequestError(fmt.Errorf("invalid rows format: %s", params["rows"]), http.StatusBadRequest)
	}

	sd, err := sg.cafe.Query(ctx, v.TraceID, pageNumber, rowsPerPage)
	if err != nil {
		return errors.Wrap(err, "unable to query for cafe")
	}

	return web.Respond(ctx, w, sd, http.StatusOK)
}

func (sg cafeGroup) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	params := web.Params(r)
	sd, err := sg.cafe.QueryByID(ctx, v.TraceID, claims, params["id"])

	if err != nil {
		if err != nil {
			switch err {
			case user.ErrInvalidID:
				return web.NewRequestError(err, http.StatusBadRequest)
			case user.ErrNotFound:
				return web.NewRequestError(err, http.StatusNotFound)
			case user.ErrForbidden:
				return web.NewRequestError(err, http.StatusForbidden)
			default:
				return errors.Wrapf(err, "ID: %s", params["id"])
			}
		}
	}

	return web.Respond(ctx, w, sd, http.StatusOK)
}

func (sg cafeGroup) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	var ng cafe.NewCafe
	if err := web.Decode(r, &ng); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	ss, err := sg.cafe.Create(ctx, v.TraceID, ng, v.Now)
	if err != nil {
		return errors.Wrapf(err, "Staff: %+v", &ss)
	}

	return web.Respond(ctx, w, ss, http.StatusCreated)
}

func (sg cafeGroup) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	var gd cafe.UpdateCafe
	if err := web.Decode(r, &gd); err != nil {
		return errors.Wrap(err, "unable to decode payload")
	}

	params := web.Params(r)
	err := sg.cafe.Update(ctx, v.TraceID, claims, params["id"], gd, v.Now)
	if err != nil {
		switch err {
		case user.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case user.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case user.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "ID: %s User: %+v", params["id"], &gd)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (sg cafeGroup) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	params := web.Params(r)
	err := sg.cafe.Delete(ctx, v.TraceID, params["id"])

	if err != nil {
		if err != nil {
			switch err {
			case user.ErrInvalidID:
				return web.NewRequestError(err, http.StatusBadRequest)
			case user.ErrNotFound:
				return web.NewRequestError(err, http.StatusNotFound)
			case user.ErrForbidden:
				return web.NewRequestError(err, http.StatusForbidden)
			default:
				return errors.Wrapf(err, "ID: %s", params["id"])
			}
		}
	}
	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
