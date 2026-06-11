package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/duddy57/cotrip-api/internal/api"
	"github.com/duddy57/cotrip-api/internal/api/mocks"
	"github.com/duddy57/cotrip-api/internal/api/spec"
	"github.com/duddy57/cotrip-api/internal/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestPostTrips(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validBody := map[string]any{
		"destination": "Paris",
		"owner_email": "john@example.com",
		"owner_name":  "John",
		"starts_at":   "2026-06-01T00:00:00Z",
		"ends_at":     "2026-06-10T00:00:00Z",
	}

	tests := []struct {
		name           string
		body           map[string]any
		setupMocks     func(store *mocks.MockQuerier)
		expectedStatus int
	}{
		{
			name: "deve criar viagem com sucesso",
			body: validBody,
			setupMocks: func(store *mocks.MockQuerier) {
				store.EXPECT().
					InsertTrip(gomock.Any(), gomock.Any()).
					Return(uuid.New(), nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "deve retornar 400 para destination curto",
			body: map[string]any{
				"destination": "NY",
				"owner_email": "john@example.com",
				"owner_name":  "John",
				"starts_at":   "2026-06-01T00:00:00Z",
				"ends_at":     "2026-06-10T00:00:00Z",
			},
			setupMocks:     func(store *mocks.MockQuerier) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "deve retornar 400 se InsertTrip falhar",
			body: validBody,
			setupMocks: func(store *mocks.MockQuerier) {
				store.EXPECT().
					InsertTrip(gomock.Any(), gomock.Any()).
					Return(uuid.Nil, errors.New("db error"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := mocks.NewMockQuerier(ctrl)
			tt.setupMocks(store)

			a := api.NewAPIWithStore(store, zap.NewNop())

			router := newTestRouter(&a)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/trips", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}


func TestGetTripsTripID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tripID := uuid.New()
	validTrip := pgstore.Trip{
		ID:          tripID,
		Destination: "Paris",
		OwnerEmail:  "john@example.com",
		OwnerName:   "John",
		StartsAt:    pgtype.Timestamp{Time: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC), Valid: true},
		EndsAt:      pgtype.Timestamp{Time: time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC), Valid: true},
		IsConfirmed: false,
	}

	tests := []struct {
		name           string
		tripID         string
		setupMocks     func(store *mocks.MockQuerier)
		expectedStatus int
	}{
		{
			name:   "deve retornar viagem com sucesso",
			tripID: tripID.String(),
			setupMocks: func(store *mocks.MockQuerier) {
				store.EXPECT().
					GetTrip(gomock.Any(), tripID).
					Return(validTrip, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "deve retornar 400 para UUID inválido",
			tripID: "invalid-uuid",
			setupMocks: func(store *mocks.MockQuerier) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "deve retornar 400 se viagem não encontrada",
			tripID: uuid.New().String(),
			setupMocks: func(store *mocks.MockQuerier) {
				store.EXPECT().
					GetTrip(gomock.Any(), gomock.Any()).
					Return(pgstore.Trip{}, pgx.ErrNoRows)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "deve retornar 400 se GetTrip falhar",
			tripID: tripID.String(),
			setupMocks: func(store *mocks.MockQuerier) {
				store.EXPECT().
					GetTrip(gomock.Any(), gomock.Any()).
					Return(pgstore.Trip{}, errors.New("db error"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := mocks.NewMockQuerier(ctrl)
			tt.setupMocks(store)

			a := api.NewAPIWithStore(store, zap.NewNop())

			router := newTestRouter(&a)

			req := httptest.NewRequest(http.MethodGet, "/trips/"+tt.tripID, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestPutTripsTripID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tripID := uuid.New()
	existingTrip := pgstore.Trip{
		ID:          tripID,
		Destination: "Paris",
		OwnerEmail:  "john@example.com",
		OwnerName:   "John",
		StartsAt:    pgtype.Timestamp{Time: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC), Valid: true},
		EndsAt:      pgtype.Timestamp{Time: time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC), Valid: true},
		IsConfirmed: false,
	}

	validBody := map[string]any{
		"destination": "London",
		"starts_at":   "2026-07-01T00:00:00Z",
		"ends_at":     "2026-07-10T00:00:00Z",
	}

	tests := []struct {
		name           string
		tripID         string
		body           map[string]any
		setupMocks     func(store *mocks.MockQuerier)
		expectedStatus int
	}{
		{
			name:   "deve atualizar viagem com sucesso",
			tripID: tripID.String(),
			body:   validBody,
			setupMocks: func(store *mocks.MockQuerier) {
				store.EXPECT().
					GetTrip(gomock.Any(), tripID).
					Return(existingTrip, nil)
				store.EXPECT().
					UpdateTrip(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:   "deve retornar 400 para UUID inválido",
			tripID: "invalid-uuid",
			body:   validBody,
			setupMocks: func(store *mocks.MockQuerier) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "deve retornar 400 para destination curto",
			tripID: tripID.String(),
			body: map[string]any{
				"destination": "NY",
				"starts_at":   "2026-07-01T00:00:00Z",
				"ends_at":     "2026-07-10T00:00:00Z",
			},
			setupMocks: func(store *mocks.MockQuerier) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "deve retornar 400 se viagem não encontrada",
			tripID: uuid.New().String(),
			body:   validBody,
			setupMocks: func(store *mocks.MockQuerier) {
				store.EXPECT().
					GetTrip(gomock.Any(), gomock.Any()).
					Return(pgstore.Trip{}, pgx.ErrNoRows)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "deve retornar 400 se UpdateTrip falhar",
			tripID: tripID.String(),
			body:   validBody,
			setupMocks: func(store *mocks.MockQuerier) {
				store.EXPECT().
					GetTrip(gomock.Any(), tripID).
					Return(existingTrip, nil)
				store.EXPECT().
					UpdateTrip(gomock.Any(), gomock.Any()).
					Return(errors.New("db error"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := mocks.NewMockQuerier(ctrl)
			tt.setupMocks(store)

			a := api.NewAPIWithStore(store, zap.NewNop())

			router := newTestRouter(&a)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, "/trips/"+tt.tripID, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func newTestRouter(a *api.API) http.Handler {
	return spec.Handler(a, spec.WithErrorHandler(a.ErrorHandlerFunc))
}
