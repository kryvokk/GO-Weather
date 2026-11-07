package settings

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Settings struct {
	TemperatureUnit string `json:"temperatureUnit"`
	WindSpeedUnit   string `json:"windSpeedUnit"`
	PressureUnit    string `json:"pressureUnit"`
	Theme           string `json:"theme"`
}
type Repository struct{ pool *pgxpool.Pool }

func NewRepository(p *pgxpool.Pool) *Repository { return &Repository{p} }
func (r *Repository) Get(ctx context.Context) (Settings, error) {
	var s Settings
	e := r.pool.QueryRow(ctx, `SELECT temperature_unit,wind_speed_unit,pressure_unit,theme FROM user_settings WHERE id=1`).Scan(&s.TemperatureUnit, &s.WindSpeedUnit, &s.PressureUnit, &s.Theme)
	return s, e
}
func (r *Repository) Put(ctx context.Context, s Settings) error {
	_, e := r.pool.Exec(ctx, `INSERT INTO user_settings(id,temperature_unit,wind_speed_unit,pressure_unit,theme) VALUES(1,$1,$2,$3,$4) ON CONFLICT(id) DO UPDATE SET temperature_unit=EXCLUDED.temperature_unit,wind_speed_unit=EXCLUDED.wind_speed_unit,pressure_unit=EXCLUDED.pressure_unit,theme=EXCLUDED.theme,updated_at=now()`, s.TemperatureUnit, s.WindSpeedUnit, s.PressureUnit, s.Theme)
	return e
}
