package city

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kryvna/go-weather/internal/weather"
)

type Repository struct{ pool *pgxpool.Pool }

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool} }
func (r *Repository) Favorites(ctx context.Context) ([]weather.City, error) {
	rows, e := r.pool.Query(ctx, `SELECT c.name,c.country,c.latitude,c.longitude,c.timezone FROM favorite_cities f JOIN cities c ON c.id=f.city_id ORDER BY f.created_at`)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []weather.City{}
	for rows.Next() {
		var c weather.City
		if e = rows.Scan(&c.Name, &c.Country, &c.Latitude, &c.Longitude, &c.Timezone); e != nil {
			return nil, e
		}
		out = append(out, c)
	}
	return out, rows.Err()
}
func (r *Repository) AddFavorite(ctx context.Context, c weather.City) error {
	tx, e := r.pool.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(ctx)
	var id int64
	e = tx.QueryRow(ctx, `INSERT INTO cities(name,country,latitude,longitude,timezone) VALUES($1,$2,$3,$4,$5) ON CONFLICT(latitude,longitude) DO UPDATE SET name=EXCLUDED.name,country=EXCLUDED.country,timezone=EXCLUDED.timezone RETURNING id`, c.Name, c.Country, c.Latitude, c.Longitude, c.Timezone).Scan(&id)
	if e != nil {
		return e
	}
	if _, e = tx.Exec(ctx, `INSERT INTO favorite_cities(city_id) VALUES($1) ON CONFLICT(city_id) DO NOTHING`, id); e != nil {
		return e
	}
	return tx.Commit(ctx)
}
func (r *Repository) RemoveFavorite(ctx context.Context, lat, lon float64) error {
	_, e := r.pool.Exec(ctx, `DELETE FROM favorite_cities f USING cities c WHERE f.city_id=c.id AND c.latitude=$1 AND c.longitude=$2`, lat, lon)
	return e
}
