package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/duddy57/cotrip-api/internal/api"
	"github.com/duddy57/cotrip-api/internal/api/mocks"
	"github.com/duddy57/cotrip-api/internal/api/spec"
	"github.com/google/uuid"
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

func newTestRouter(a *api.API) http.Handler {
	return spec.Handler(a, spec.WithErrorHandler(a.ErrorHandlerFunc))
}
