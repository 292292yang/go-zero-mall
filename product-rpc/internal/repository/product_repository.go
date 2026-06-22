package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/292292yang/go-zero-mall/common/errorx"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const productColumns = "id, name, description, price, stock, status, created_at, updated_at"

type Product struct {
	Id          int64     `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       int64     `db:"price"`
	Stock       int64     `db:"stock"`
	Status      int8      `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type ProductCreate struct {
	Name        string
	Description string
	Price       int64
	Stock       int64
	Status      int8
}

type ProductPatch struct {
	Name        *string
	Description *string
	Price       *int64
	Stock       *int64
	Status      *int8
}

type ProductListQuery struct {
	Name        *string
	Description *string
	Status      *int8
	Page        uint64
	Size        uint64
}

type ProductListResult struct {
	List  []*Product
	Total uint64
}

type ProductRepository interface {
	Create(ctx context.Context, p ProductCreate) (uint64, error)
	FindById(ctx context.Context, id uint64) (*Product, error)
	FindByName(ctx context.Context, name string) (*Product, error)
	List(ctx context.Context, p ProductListQuery) (*ProductListResult, error)
	UpdatePatch(ctx context.Context, id uint64, patch ProductPatch) error
	Delete(ctx context.Context, id uint64) error
}

type productRepository struct {
	conn  sqlx.SqlConn
	redis *redis.Redis
	table string
}

func NewProductRepository(conn sqlx.SqlConn, rds *redis.Redis) *productRepository {
	return &productRepository{
		conn:  conn,
		redis: rds,
		table: "product",
	}
}

func (r *productRepository) Create(ctx context.Context, p ProductCreate) (uint64, error) {
	query, args, err := sq.
		Insert(r.table).
		Columns("name", "description", "price", "stock", "status").
		Values(p.Name, p.Description, p.Price, p.Stock, p.Status).
		ToSql()
	if err != nil {
		return 0, err
	}
	result, err := r.conn.ExecCtx(ctx, query, args...)
	if err != nil {
		if isDuplicateEntry(err) {
			return 0, errorx.NewCodeError(errorx.ProductNotFound, "user already exists")
		}
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

func (r *productRepository) FindById(ctx context.Context, id uint64) (*Product, error) {
	key := r.userDetailKey(id)

	cached, err := r.redis.GetCtx(ctx, key)
	if err == nil && cached != "" {
		var p Product
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

func (r *productRepository) FindByName(ctx context.Context, username string) (*Product, error) {
	u, err := r.findByUsernameFromDB(ctx, username)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *productRepository) findByIdFromDB(ctx context.Context, id uint64) (*Product, error) {
	query, args, err := sq.
		Select(productColumns).
		From(r.table).
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var p Product
	err = r.conn.QueryRowCtx(ctx, &p, query, args...)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, errorx.NewCodeError(errorx.UserNotFound, "user not found")
		}
		return nil, err
	}

	return &p, nil
}

func (r *productRepository) findByUsernameFromDB(ctx context.Context, username string) (*Product, error) {
	query, args, err := sq.
		Select(productColumns).
		From(r.table).
		Where(sq.Eq{"username": username}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}
	var u Product
	err = r.conn.QueryRowCtx(ctx, &u, query, args...)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, errorx.NewCodeError(errorx.UserNotFound, "user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *productRepository) List(ctx context.Context, u ProductListQuery) (*ProductListResult, error) {
	page, size := normalizePage(u.Page, u.Size)

	applyWhere := func(b sq.SelectBuilder) sq.SelectBuilder {
		name := strings.TrimSpace(*u.Name)
		if name != "" {
			b = b.Where(sq.Like{"name": "%" + name + "%"})
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
		Select(productColumns).
		From(r.table).
		OrderBy("id DESC").
		Limit(size).
		Offset((page - 1) * size)
	listBuilder = applyWhere(listBuilder)

	listQuery, listArgs, err := listBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var list []*Product
	err = r.conn.QueryRowsCtx(ctx, &list, listQuery, listArgs...)
	if err != nil {
		return nil, err
	}

	return &ProductListResult{
		List:  list,
		Total: total,
	}, nil
}

func (r *productRepository) UpdatePatch(ctx context.Context, id uint64, patch ProductPatch) error {
	builder := sq.Update(r.table).Where(sq.Eq{"id": id})
	hasSet := false

	if patch.Name != nil {
		builder = builder.Set("name", *patch.Name)
		hasSet = true
	}
	if patch.Name != nil {
		builder = builder.Set("description", *patch.Name)
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
		if isDuplicateEntry(err) {
			return errorx.NewCodeError(errorx.UserAlreadyExists, "user already exists")
		}
		return err
	}

	rows, err := result.RowsAffected()
	if err == nil && rows == 0 {
		return errorx.NewCodeError(errorx.UserAlreadyExists, "user already exists")
	}

	_, _ = r.redis.DelCtx(ctx, r.userDetailKey(id))
	return nil
}

func (r *productRepository) Delete(ctx context.Context, id uint64) error {
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

	_, _ = r.redis.DelCtx(ctx, r.userDetailKey(id))
	return nil
}

func (r *productRepository) userDetailKey(id uint64) string {
	return fmt.Sprintf("user:detail:%d", id)
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
