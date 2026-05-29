// config/mail.go
package infra

import (
	"context"
	"fmt"

	"github.com/duddy57/cotrip-api/internal/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/resend/resend-go/v3"
)

type Mail struct {
	client *resend.Client
	from   string
	q      *pgstore.Queries
}

func NewMail(pool *pgxpool.Pool, apiKey, from string) *Mail {
	return &Mail{
		client: resend.NewClient(apiKey),
		from:   from,
		q:      pgstore.New(pool),
	}
}

func (m *Mail) SendTripConfirmationEmail(tripID uuid.UUID) error {
	ctx := context.Background()
	trip, err := m.q.GetTrip(ctx, tripID)
	if err != nil {
		return err
	}

	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("Acme <%s>", m.from),
		To:      []string{trip.OwnerEmail},
		Html:    "<strong>Você deve confirmar sua viagem</strong>",
		Subject: "Confirme sua viagem",
	}

	if _, err = m.client.Emails.Send(params); err != nil {
		return err
	}

	return nil
}

func (m *Mail) SendTripConfirmedEmails(tripID uuid.UUID) error {
	participants, err := m.q.GetParticipants(context.Background(), tripID)
	if err != nil {
		return err
	}

	for _, p := range participants {
		params := &resend.SendEmailRequest{
			From:    fmt.Sprintf("Acme <%s>", m.from),
			To:      []string{p.Email},
			Html:    "<strong>Você deve confirmar sua viagem</strong>",
			Subject: "Confirme sua viagem",
		}

		if _, err = m.client.Emails.Send(params); err != nil {
			return err
		}
	}

	return nil
}

func (m *Mail) SendTripConfirmedEmail(tripID, participantID uuid.UUID) error {
	ctx := context.Background()
	participant, err := m.q.GetParticipant(ctx, participantID)
	if err != nil {
		return err
	}

	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("Acme <%s>", m.from),
		To:      []string{participant.Email},
		Html:    "<strong>Você deve confirmar sua viagem</strong>",
		Subject: "Confirme sua viagem",
	}

	if _, err = m.client.Emails.Send(params); err != nil {
		return err
	}
	
	return nil
}
