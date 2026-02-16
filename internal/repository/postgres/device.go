package postgres

import (
	"context"

	model "github.com/Kowari1/File-Handler/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeviceRepository struct {
	pool *pgxpool.Pool
}

func NewDeviceRepository(pool *pgxpool.Pool) *DeviceRepository {
	return &DeviceRepository{
		pool: pool,
	}
}

func (r *DeviceRepository) Save(
	ctx context.Context,
	devices []*model.Device,
) error {

	if len(devices) == 0 {
		return nil
	}

	rows := make([][]any, 0, len(devices))

	for _, d := range devices {
		rows = append(rows, []any{
			d.UnitGUID,
			d.N,
			d.MQTT,
			d.Invid,
			d.MsgID,
			d.Text,
			d.Context,
			d.Class,
			d.Level,
			d.Area,
			d.Addr,
			d.Block,
			d.Type,
			d.Bit,
			d.InvertBit,
		})
	}

	columns := []string{
		"unit_guid",
		"n",
		"mqtt",
		"invid",
		"msg_id",
		"text",
		"context",
		"class",
		"level",
		"area",
		"addr",
		"block",
		"type",
		"bit",
		"invert_bit",
	}

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Conn().CopyFrom(
		ctx,
		pgx.Identifier{"device"},
		columns,
		pgx.CopyFromRows(rows),
	)

	return err
}

func (r *DeviceRepository) FindByUnitGUID(
	ctx context.Context,
	guid uuid.UUID,
) ([]*model.Device, error) {

	query := `
		SELECT n, mqtt, invid, unit_guid, msg_id, text,
		       context, class, level, area, addr,
		       block, type, bit, invert_bit
		FROM device
		WHERE unit_guid = $1
		ORDER BY id
	`

	rows, err := r.pool.Query(ctx, query, guid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []*model.Device

	for rows.Next() {
		var d model.Device
		err := rows.Scan(
			&d.N,
			&d.MQTT,
			&d.Invid,
			&d.UnitGUID,
			&d.MsgID,
			&d.Text,
			&d.Context,
			&d.Class,
			&d.Level,
			&d.Area,
			&d.Addr,
			&d.Block,
			&d.Type,
			&d.Bit,
			&d.InvertBit,
		)
		if err != nil {
			return nil, err
		}
		devices = append(devices, &d)
	}

	return devices, nil
}

func (r *DeviceRepository) FindAll(ctx context.Context, limit, offset int) ([]model.Device, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT n, mqtt, invid, unit_guid, msg_id, text, context,
               class, level, area, addr, block, type, bit, invert_bit
        FROM device
        ORDER BY id
        LIMIT $1 OFFSET $2
    `, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanDevices(rows)
}

func (r *DeviceRepository) FindLimitedByUnitGUID(ctx context.Context, guid uuid.UUID, limit, offset int) ([]model.Device, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT n, mqtt, invid, unit_guid, msg_id, text, context,
               class, level, area, addr, block, type, bit, invert_bit
        FROM device
	   WHERE unit_guid = $1
        ORDER BY id
        LIMIT $2 OFFSET $3
    `, guid, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanDevices(rows)
}

func (r *DeviceRepository) CountAll(ctx context.Context) (int, error) {
	var total int
	err := r.pool.QueryRow(ctx, `
        SELECT COUNT(*) FROM device
    `).Scan(&total)

	return total, err
}

func (r *DeviceRepository) CountByUnitGUID(
	ctx context.Context,
	guid uuid.UUID,
) (int, error) {

	var total int

	err := r.pool.QueryRow(ctx, `
        SELECT COUNT(*)
        FROM device
        WHERE unit_guid = $1
    `, guid).Scan(&total)

	return total, err
}

func scanDevices(rows pgx.Rows) ([]model.Device, error) {
	var devices []model.Device

	for rows.Next() {
		var d model.Device

		if err := rows.Scan(
			&d.N,
			&d.MQTT,
			&d.Invid,
			&d.UnitGUID,
			&d.MsgID,
			&d.Text,
			&d.Context,
			&d.Class,
			&d.Level,
			&d.Area,
			&d.Addr,
			&d.Block,
			&d.Type,
			&d.Bit,
			&d.InvertBit,
		); err != nil {
			return nil, err
		}

		devices = append(devices, d)
	}

	return devices, rows.Err()
}
