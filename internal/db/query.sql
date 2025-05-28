-- name: GetWeatherByLocation :one
SELECT id, city, query_time, weather_data
FROM weather_query_history
WHERE city = $1
ORDER BY query_time DESC
LIMIT 1;

-- name: InsertWeatherQuery :one
INSERT INTO weather_query_history (city, weather_data)
VALUES ($1, $2::jsonb)
RETURNING id, city, query_time, weather_data;

-- name: GetRecentWeather :many
SELECT id, city, query_time, weather_data
FROM weather_query_history
WHERE city = $1 AND query_time >= $2
ORDER BY query_time DESC;