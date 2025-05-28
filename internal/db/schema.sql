CREATE TABLE weather_query_history (
    id SERIAL PRIMARY KEY,
    city VARCHAR(100) NOT NULL,
    query_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    weather_data JSONB NOT NULL
);

CREATE INDEX idx_weather_query_history_city ON weather_query_history (city);
CREATE INDEX idx_weather_query_history_query_time ON weather_query_history (query_time);