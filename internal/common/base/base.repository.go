package base

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)


type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	CreateBatch(ctx context.Context, entities []*T) error
	FindByID(ctx context.Context, id uuid.UUID) (*T, error)
	FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*T, error)
	FindOne(ctx context.Context, conditions map[string]interface{}) (*T, error)
	FindAll(ctx context.Context, limit, offset int) ([]*T, error)
	Find(ctx context.Context, conditions map[string]interface{}, limit, offset int) ([]*T, error)
	Update(ctx context.Context, entity *T) error
	UpdateFields(ctx context.Context, id uuid.UUID, fields map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
	HardDelete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int64, error)
	CountWhere(ctx context.Context, conditions map[string]interface{}) (int64, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	ExistsWhere(ctx context.Context, conditions map[string]interface{}) (bool, error)
	WithTx(tx *gorm.DB) Repository[T]
	WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}


type BaseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *BaseRepository[T]) CreateBatch(ctx context.Context, entities []*T) error {
	if len(entities) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(entities, 100).Error
}

func (r *BaseRepository[T]) FindByID(ctx context.Context, id uuid.UUID) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&entity, "id = ?", id).
		Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) FindOne(ctx context.Context, conditions map[string]interface{}) (*T, error) {
	var entity T
	query := r.db.WithContext(ctx)
	for key, value := range conditions {
		query = query.Where(key+" = ?", value)
	}
	err := query.First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) FindAll(ctx context.Context, limit, offset int) ([]*T, error) {
	var entities []*T
	query := r.db.WithContext(ctx)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&entities).Error
	return entities, err
}

func (r *BaseRepository[T]) Find(ctx context.Context, conditions map[string]interface{}, limit, offset int) ([]*T, error) {
	var entities []*T
	query := r.db.WithContext(ctx)

	// Apply conditions
	for key, value := range conditions {
		query = query.Where(key+" = ?", value)
	}

	// Apply pagination
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&entities).Error
	return entities, err
}

func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *BaseRepository[T]) UpdateFields(ctx context.Context, id uuid.UUID, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(new(T)).Where("id = ?", id).Updates(fields).Error
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(new(T), "id = ?", id).Error
}

func (r *BaseRepository[T]) HardDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Delete(new(T), "id = ?", id).Error
}

func (r *BaseRepository[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(new(T)).Count(&count).Error
	return count, err
}

func (r *BaseRepository[T]) CountWhere(ctx context.Context, conditions map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(new(T))

	for key, value := range conditions {
		query = query.Where(key+" = ?", value)
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *BaseRepository[T]) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(new(T)).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *BaseRepository[T]) ExistsWhere(ctx context.Context, conditions map[string]interface{}) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(new(T))

	for key, value := range conditions {
		query = query.Where(key+" = ?", value)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

func (r *BaseRepository[T]) WithTx(tx *gorm.DB) Repository[T] {
	return NewBaseRepository[T](tx)
}


type QueryBuilder[T any] struct {
	db        *gorm.DB
	preloads  []string
	orderBy   string
	conditions map[string]interface{}
}

func (r *BaseRepository[T]) NewQueryBuilder() *QueryBuilder[T] {
	return &QueryBuilder[T]{
		db:         r.db,
		conditions: make(map[string]interface{}),
	}
}

func (qb *QueryBuilder[T]) Where(field string, value interface{}) *QueryBuilder[T] {
	qb.conditions[field+" = ?"] = value
	return qb
}

func (qb *QueryBuilder[T]) WhereIn(field string, values []interface{}) *QueryBuilder[T] {
	qb.db = qb.db.Where(field+" IN ?", values)
	return qb
}

func (qb *QueryBuilder[T]) WhereLike(field string, value string) *QueryBuilder[T] {
	qb.db = qb.db.Where(field+" LIKE ?", "%"+value+"%")
	return qb
}

func (qb *QueryBuilder[T]) Preload(preload string) *QueryBuilder[T] {
	qb.preloads = append(qb.preloads, preload)
	return qb
}

func (qb *QueryBuilder[T]) OrderBy(orderBy string) *QueryBuilder[T] {
	qb.orderBy = orderBy
	return qb
}

func (qb *QueryBuilder[T]) Limit(limit int) *QueryBuilder[T] {
	qb.db = qb.db.Limit(limit)
	return qb
}

func (qb *QueryBuilder[T]) Offset(offset int) *QueryBuilder[T] {
	qb.db = qb.db.Offset(offset)
	return qb
}

func (qb *QueryBuilder[T]) Find(ctx context.Context) ([]*T, error) {
	var entities []*T
	query := qb.db.WithContext(ctx)

	// Apply conditions
	for field, value := range qb.conditions {
		query = query.Where(field, value)
	}

	// Apply preloads
	for _, preload := range qb.preloads {
		query = query.Preload(preload)
	}

	// Apply order
	if qb.orderBy != "" {
		query = query.Order(qb.orderBy)
	}

	err := query.Find(&entities).Error
	return entities, err
}

func (qb *QueryBuilder[T]) FindOne(ctx context.Context) (*T, error) {
	var entity T
	query := qb.db.WithContext(ctx)

	// Apply conditions
	for field, value := range qb.conditions {
		query = query.Where(field, value)
	}

	// Apply preloads
	for _, preload := range qb.preloads {
		query = query.Preload(preload)
	}

	err := query.First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (qb *QueryBuilder[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	query := qb.db.WithContext(ctx).Model(new(T))

	// Apply conditions
	for field, value := range qb.conditions {
		query = query.Where(field, value)
	}

	err := query.Count(&count).Error
	return count, err
}


func (r *BaseRepository[T]) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

func Transaction[T any](ctx context.Context, db *gorm.DB, fn func(repo Repository[T]) error) error {
	repo := NewBaseRepository[T](db)
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(repo.WithTx(tx))
	})
}


func (r *BaseRepository[T]) Upsert(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(entity).Error
}

func (r *BaseRepository[T]) UpsertBatch(ctx context.Context, entities []*T) error {
	if len(entities) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(entities, 100).Error
}


type PaginationResult struct {
	Total       int64 `json:"total"`
	Page        int   `json:"page"`
	PerPage     int   `json:"per_page"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrevious bool  `json:"has_previous"`
}

func (r *BaseRepository[T]) Paginate(ctx context.Context, page, perPage int, conditions map[string]interface{}) ([]*T, *PaginationResult, error) {
	var entities []*T
	var total int64

	query := r.db.WithContext(ctx).Model(new(T))

	// Apply conditions
	for key, value := range conditions {
		query = query.Where(key+" = ?", value)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Calculate pagination
	offset := (page - 1) * perPage
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	// Fetch data
	if err := query.Offset(offset).Limit(perPage).Find(&entities).Error; err != nil {
		return nil, nil, err
	}

	// Build pagination metadata
	pagination := &PaginationResult{
		Total:       total,
		Page:        page,
		PerPage:     perPage,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}

	return entities, pagination, nil
}
