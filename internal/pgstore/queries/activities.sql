-- name: CreateActivity :one
INSERT INTO activities
( "trip_id", "title", "occurs_at" ) VALUES
    ( $1, $2, $3 )
    RETURNING "id";

-- name: GetTripActivities :many
SELECT
    "id", "trip_id", "title", "occurs_at"
FROM activities
WHERE
    trip_id = $1;
