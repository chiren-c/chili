package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type ArticleAuthorDAO interface {
	Insert(ctx context.Context, art ArticleAuthor) (int64, error)
	UpdateById(ctx context.Context, art ArticleAuthor) error
	GetByAuthor(ctx context.Context, author int64, offset, limit int) ([]ArticleAuthor, error)
	GetById(ctx context.Context, id int64) (ArticleAuthor, error)
	Sync(ctx context.Context, art ArticleAuthor) (int64, error)
	SyncStatus(ctx context.Context, author, id int64, status uint8) error
}

type GORMArticleAuthorDAO struct {
	db *gorm.DB
}

func (g *GORMArticleAuthorDAO) Insert(ctx context.Context, art ArticleAuthor) (int64, error) {
	now := time.Now().UnixMilli()
	art.Utime = now
	art.Ctime = now
	err := g.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

func (g *GORMArticleAuthorDAO) UpdateById(ctx context.Context, art ArticleAuthor) error {
	now := time.Now().UnixMilli()
	res := g.db.WithContext(ctx).
		Where("id=? and author_id=?", art.Id, art.AuthorId).
		Updates(map[string]any{
			"title":   art.Title,
			"content": art.Content,
			"status":  art.Status,
			"utime":   now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("更新数据失败")
	}
	return nil
}

func (g *GORMArticleAuthorDAO) GetByAuthor(ctx context.Context, author int64, offset, limit int) ([]ArticleAuthor, error) {
	var arts []ArticleAuthor
	err := g.db.WithContext(ctx).
		Model(&ArticleAuthor{}).
		Where("author_id=?", author).
		Offset(offset).
		Limit(limit).
		Order("utime DES").
		Find(&arts).Error
	return arts, err
}

func (g *GORMArticleAuthorDAO) GetById(ctx context.Context, id int64) (ArticleAuthor, error) {
	var art ArticleAuthor
	err := g.db.WithContext(ctx).Where("id=?", id).First(&art).Error
	return art, err
}

func (g *GORMArticleAuthorDAO) Sync(ctx context.Context, art ArticleAuthor) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GORMArticleAuthorDAO) SyncStatus(ctx context.Context, author, id int64, status uint8) error {
	//TODO implement me
	panic("implement me")
}

func NewGORMArticleAuthorDAO(db *gorm.DB) ArticleAuthorDAO {
	return &GORMArticleAuthorDAO{
		db: db,
	}
}

type ArticleAuthor struct {
	Id      int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Title   string `gorm:"type=varchar(4096)" bson:"title,omitempty"`
	Content string `gorm:"type=BLOB" bson:"content,omitempty"`
	// 作者
	AuthorId int64 `gorm:"index" bson:"author_id,omitempty"`
	Status   uint8 `bson:"status,omitempty"`
	Ctime    int64 `bson:"ctime,omitempty"`
	Utime    int64 `bson:"utime,omitempty" gorm:"index"`
}

func (ArticleAuthor) TableName() string {
	return "article_author"
}
