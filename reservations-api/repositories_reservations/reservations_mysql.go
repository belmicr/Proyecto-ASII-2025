package repositories_reservations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"reservations/domain_reservations"

	// Import por efectos (registra el driver "mysql" en database/sql)
	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	db *sql.DB
}

type MySQLConfig struct {
	Host   string
	Port   string
	User   string
	Pass   string
	DB     string
	Params string // ej: "parseTime=true&charset=utf8mb4"
}

func NewMySQL(cfg MySQLConfig) (*MySQL, error) {
	if cfg.Params == "" {
		cfg.Params = "parseTime=true&charset=utf8mb4"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.DB, cfg.Params)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &MySQL{db: db}, nil
}

func (m *MySQL) Create(ctx context.Context, r domain_reservations.Reservation) (domain_reservations.Reservation, error) {
	if strings.TrimSpace(r.ID) == "" {
		r.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	const q = `
		INSERT INTO reservations
		(id, hotel_id, user_id, check_in, check_out, guests, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := m.db.ExecContext(ctx, q,
		r.ID, r.HotelID, r.UserID, r.CheckIn, r.CheckOut, r.Guests, r.Status)
	if err != nil {
		return domain_reservations.Reservation{}, err
	}
	return r, nil
}

func (m *MySQL) Update(ctx context.Context, id string, r domain_reservations.Reservation) (domain_reservations.Reservation, error) {
	const q = `
		UPDATE reservations
		SET hotel_id = ?, user_id = ?, check_in = ?, check_out = ?, guests = ?, status = ?
		WHERE id = ?`
	res, err := m.db.ExecContext(ctx, q,
		r.HotelID, r.UserID, r.CheckIn, r.CheckOut, r.Guests, r.Status, id)
	if err != nil {
		return domain_reservations.Reservation{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain_reservations.Reservation{}, sql.ErrNoRows
	}
	r.ID = id
	return r, nil
}

func (m *MySQL) GetByID(ctx context.Context, id string) (domain_reservations.Reservation, error) {
	const q = `
		SELECT id, hotel_id, user_id, check_in, check_out, guests, status
		FROM reservations
		WHERE id = ?`
	var r domain_reservations.Reservation
	err := m.db.QueryRowContext(ctx, q, id).
		Scan(&r.ID, &r.HotelID, &r.UserID, &r.CheckIn, &r.CheckOut, &r.Guests, &r.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain_reservations.Reservation{}, err
		}
		return domain_reservations.Reservation{}, err
	}
	return r, nil
}

func (m *MySQL) List(ctx context.Context, hotelID, userID, status string) ([]domain_reservations.Reservation, error) {
	var sb strings.Builder
	args := []any{}

	sb.WriteString(`SELECT id, hotel_id, user_id, check_in, check_out, guests, status
	                FROM reservations WHERE 1=1`)
	if hotelID != "" {
		sb.WriteString(" AND hotel_id = ?")
		args = append(args, hotelID)
	}
	if userID != "" {
		sb.WriteString(" AND user_id = ?")
		args = append(args, userID)
	}
	if status != "" {
		sb.WriteString(" AND status = ?")
		args = append(args, status)
	}
	sb.WriteString(" ORDER BY created_at DESC")

	rows, err := m.db.QueryContext(ctx, sb.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain_reservations.Reservation
	for rows.Next() {
		var r domain_reservations.Reservation
		if err := rows.Scan(&r.ID, &r.HotelID, &r.UserID, &r.CheckIn, &r.CheckOut, &r.Guests, &r.Status); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}
