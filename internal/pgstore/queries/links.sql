-- name: CreateTripLink :one
INSERT INTO links
( "trip_id", "title", "url" ) VALUES
    ( $1, $2, $3 )
    RETURNING "id";

-- name: GetTripLinks :many
SELECT
    "id", "trip_id", "title", "url"
FROM links
WHERE
    trip_id = $1;