package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/292292yang/go-zero-mall/common/cryptx"
	"github.com/292292yang/go-zero-mall/common/errorx"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const userColumns = "id, username, password, nickname, avatar, status, created_at, updated_at"

type User struct {
	Id        int64     `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Nickname  string    `db:"nickname"`
	Avatar    string    `db:"avatar"`
	Status    int8      `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserCreate struct {
	Username string
	Password string
	Nickname string
	Avatar   string
	Status   int8
}

type UserPatch struct {
	Username *string
	Password *string
	Nickname *string
	Avatar   *string
	Status   *int8
}

type UserListQuery struct {
	Username *string
	Nickname *string
	Status   *int8
	Page     uint64
	Size     uint64
}

type UserListResult struct {
	List  []*User
	Total uint64
}

type UserRepository interface {
	Create(ctx context.Context, u UserCreate) (uint64, error)
	FindById(ctx context.Context, id uint64) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	List(ctx context.Context, u UserListQuery) (*UserListResult, error)
	UpdatePatch(ctx context.Context, id uint64, patch UserPatch) error
	Delete(ctx context.Context, id uint64) error
}

type userRepository struct {
	conn  sqlx.SqlConn
	redis *redis.Redis
	table string
}

func NewUserRepository(conn sqlx.SqlConn, rds *redis.Redis) UserRepository {
	return &userRepository{
		conn:  conn,
		redis: rds,
		table: "user",
	}
}

func (r *userRepository) Create(ctx context.Context, u UserCreate) (uint64, error) {
	password, err := cryptx.EncryptPassword(u.Password)
	if err != nil {
		return 0, errorx.NewCodeError(errorx.ServerError, "system error")
	}
	query, args, err := sq.
		Insert(r.table).
		Columns("username", "password", "nickname", "avatar", "status").
		Values(u.Username, password, u.Nickname, u.Avatar, u.Status).
		ToSql()
	if err != nil {
		return 0, err
	}
	result, err := r.conn.ExecCtx(ctx, query, args...)
	if err != nil {
		if isDuplicateEntry(err) {
			return 0, errorx.NewCodeError(errorx.UserAlreadyExists, "user already exists")
		}
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

func (r *userRepository) FindById(ctx context.Context, id uint64) (*User, error) {
	key := r.userDetailKey(id)

	cached, err := r.redis.GetCtx(ctx, key)
	if err == nil && cached != "" {
		var u User
		if json.Unmarshal([]byte(cached), &u) == nil {
			return &u, nil
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

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	u, err := r.findByUsernameFromDB(ctx, username)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepository) findByIdFromDB(ctx context.Context, id uint64) (*User, error) {
	query, args, err := sq.
		Select(userColumns).
		From(r.table).
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var u User
	err = r.conn.QueryRowCtx(ctx, &u, query, args...)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, errorx.NewCodeError(errorx.UserNotFound, "user not found")
		}
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) findByUsernameFromDB(ctx context.Context, username string) (*User, error) {
	query, args, err := sq.
		Select(userColumns).
		From(r.table).
		Where(sq.Eq{"username": username}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}
	var u User
	err = r.conn.QueryRowCtx(ctx, &u, query, args...)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, errorx.NewCodeError(errorx.UserNotFound, "user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) List(ctx context.Context, u UserListQuery) (*UserListResult, error) {
	page, size := normalizePage(u.Page, u.Size)

	applyWhere := func(b sq.SelectBuilder) sq.SelectBuilder {
		name := strings.TrimSpace(*u.Username)
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
		Select(userColumns).
		From(r.table).
		OrderBy("id DESC").
		Limit(size).
		Offset((page - 1) * size)
	listBuilder = applyWhere(listBuilder)

	listQuery, listArgs, err := listBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var list []*User
	err = r.conn.QueryRowsCtx(ctx, &list, listQuery, listArgs...)
	if err != nil {
		return nil, err
	}

	return &UserListResult{
		List:  list,
		Total: total,
	}, nil
}

func (r *userRepository) UpdatePatch(ctx context.Context, id uint64, patch UserPatch) error {
	builder := sq.Update(r.table).Where(sq.Eq{"id": id})
	hasSet := false

	if patch.Username != nil {
		builder = builder.Set("username", *patch.Username)
		hasSet = true
	}
	if patch.Nickname != nil {
		builder = builder.Set("nickname", *patch.Nickname)
		hasSet = true
	}
	if patch.Password != nil {
		builder = builder.Set("password", *patch.Password)
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

func (r *userRepository) Delete(ctx context.Context, id uint64) error {
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

func (r *userRepository) userDetailKey(id uint64) string {
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
