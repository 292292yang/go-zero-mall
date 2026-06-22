package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/292292yang/go-zero-mall/common/errorx"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const orderColumns = "id, user_id, product_id, quantity, price, status, created_at, updated_at"

type Order struct {
	Id        int64     `db:"id"`
	UserId    int64     `db:"user_id"`
	ProductId int64     `db:"product_id"`
	Quantity  int64     `db:"quantity"`
	Price     int64     `db:"price"`
	Status    int8      `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type OrderCreate struct {
	UserId    int64
	ProductId int64
	Quantity  int64
	Price     int64
	Status    int8
}

type OrderPatch struct {
	Id        *int64
	UserId    *int64
	ProductId *int64
	Quantity  *int64
	Price     *int64
	Status    *int8
}

type OrderListQuery struct {
	UserId    *int64
	ProductId *int64
	Quantity  *int64
	Price     *int64
	Status    *int8
	Page      uint64
	Size      uint64
}

type OrderListResult struct {
	List  []*Order
	Total uint64
}

type OrderRepository interface {
	Create(ctx context.Context, o OrderCreate) (uint64, error)
	FindById(ctx context.Context, id uint64) (*Order, error)
	List(ctx context.Context, o OrderListQuery) (*OrderListResult, error)
	UpdatePatch(ctx context.Context, id uint64, patch OrderPatch) error
	Delete(ctx context.Context, id uint64) error
}

type orderRepository struct {
	conn  sqlx.SqlConn
	redis *redis.Redis
	table string
}

func NewOrderRepository(conn sqlx.SqlConn, rds *redis.Redis) *orderRepository {
	return &orderRepository{
		conn:  conn,
		redis: rds,
		table: "order",
	}
}

func (r *orderRepository) Create(ctx context.Context, o OrderCreate) (uint64, error) {
	query, args, err := sq.
		Insert(r.table).
		Columns("user_id", "product_id", "quantity", "price", "status").
		Values(o.UserId, o.ProductId, o.Quantity, o.Price, o.Status).
		ToSql()
	if err != nil {
		return 0, err
	}
	result, err := r.conn.ExecCtx(ctx, query, args...)
	if err != nil {
		if isDuplicateEntry(err) {
			return 0, errorx.NewCodeError(errorx.OrderCreateFailed, "order already exists")
		}
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

func (r *orderRepository) FindById(ctx context.Context, id uint64) (*Order, error) {
	key := r.orderDetailKey(id)

	cached, err := r.redis.GetCtx(ctx, key)
	if err == nil && cached != "" {
		var p Order
		if json.Unmarshal([]byte(cached), &p) == nil {
			return &p, nil
		}
	}

	p, err := r.findByIdFromDB(ctx, id)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(p)
	if err == nil {
		_ = r.redis.SetexCtx(ctx, key, string(data), 300)
	}

	return p, nil
}

func (r *orderRepository) findByIdFromDB(ctx context.Context, id uint64) (*Order, error) {
	query, args, err := sq.
		Select(orderColumns).
		From(r.table).
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var p Order
	err = r.conn.QueryRowCtx(ctx, &p, query, args...)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, errorx.NewCodeError(errorx.UserNotFound, "user not found")
		}
		return nil, err
	}

	return &p, nil
}

func (r *orderRepository) List(ctx context.Context, o OrderListQuery) (*OrderListResult, error) {
	page, size := normalizePage(o.Page, o.Size)

	applyWhere := func(b sq.SelectBuilder) sq.SelectBuilder {
		if o.UserId != nil {
			b = b.Where(sq.Eq{"user_id": o.UserId})
		}
		return b
	}

	countBuilder := sq.Select("COUNT(*)").From(r.table)
	countBuilder = applyWhere(countBuilder)

	countQuery, countArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var total uint64
	err = r.conn.QueryRowCtx(ctx, &total, countQuery, countArgs...)
	if err != nil {
		return nil, err
	}

	listBuilder := sq.
		Select(orderColumns).
		From(r.table).
		OrderBy("id DESC").
		Limit(size).
		Offset((page - 1) * size)
	listBuilder = applyWhere(listBuilder)

	listQuery, listArgs, err := listBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var list []*Order
	err = r.conn.QueryRowsCtx(ctx, &list, listQuery, listArgs...)
	if err != nil {
		return nil, err
	}

	return &OrderListResult{
		List:  list,
		Total: total,
	}, nil
}

func (r *orderRepository) UpdatePatch(ctx context.Context, id uint64, patch OrderPatch) error {
	builder := sq.Update(r.table).Where(sq.Eq{"id": id})
	hasSet := false

	if patch.ProductId != nil {
		builder = builder.Set("product_id", *patch.ProductId)
		hasSet = true
	}
	if patch.Quantity != nil {
		builder = builder.Set("quantity", *patch.Quantity)
		hasSet = true
	}
	if patch.Status != nil {
		builder = builder.Set("status", *patch.Status)
		hasSet = true
	}
	if !hasSet {
		return errorx.NewCodeError(errorx.InvalidParam, "update patch is empty")
	}
	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}
	result, err := r.conn.ExecCtx(ctx, query, args...)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err == nil && rows == 0 {
		return errorx.NewCodeError(errorx.OrderStatusInvalid, "order not exists")
	}
	_, _ = r.redis.DelCtx(ctx, r.orderDetailKey(id))
	return nil
}

func (r *orderRepository) Delete(ctx context.Context, id uint64) error {
	query, args, err := sq.Delete(r.table).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	result, err := r.conn.ExecCtx(ctx, query, args...)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err == nil && rows == 0 {
		return errorx.NewCodeError(errorx.UserNotFound, "user not found")
	}

	_, _ = r.redis.DelCtx(ctx, r.orderDetailKey(id))
	return nil
}

func (r *orderRepository) orderDetailKey(id uint64) string {
	return fmt.Sprintf("order:detail:%d", id)
}

func normalizePage(page, size uint64) (uint64, uint64) {
	if page == 0 {
		page = 1
	}
	if size == 0 || size > 100 {
		size = 20
	}
	return page, size
}

func isDuplicateEntry(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return false
}
