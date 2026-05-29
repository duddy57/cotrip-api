-- name: GetParticipant :one
SELECT
    "id", "trip_id", "email", "is_confirmed"
FROM participants
WHERE
    id = $1;

-- name: ConfirmParticipant :exec
UPDATE participants
SET "is_confirmed" = true
WHERE id = $1;

-- name: GetParticipants :many
SELECT
    "id", "trip_id", "email", "is_confirmed"
FROM participants
WHERE
    trip_id = $1;

-- name: InviteParticipantToTrip :one
INSERT INTO participants
( "trip_id", "email" ) VALUES
    ( $1, $2 )
    RETURNING "id";

-- name: InviteParticipantsToTrip :copyfrom
INSERT INTO participants
( "trip_id", "email" ) VALUES
    ( $1, $2 );
