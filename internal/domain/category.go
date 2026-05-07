package domain

import "context"

type Category struct {
	DbID   int    `json:"-" gorm:"column:id;primaryKey;autoIncrement"`
	ID     int    `json:"id" gorm:"column:wb_id;not null"`
	Parent int    `json:"parent_id,omitempty" gorm:"column:parent_id;index"`
	Name   string `json:"name" gorm:"column:name"`
	Url    string `json:"url" gorm:"column:url"`
	Shard  string `json:"shard,omitempty" gorm:"column:shard"`
	Query  string `json:"query,omitempty" gorm:"column:query"`

	SearchQuery string `json:"searchQuery,omitempty" gorm:"-"`

	Childs []Category `json:"childs,omitempty" gorm:"-"`
}

type CategoryRepository interface {
	Create(ctx context.Context, cat *Category) error
	Delete(ctx context.Context, id int) error
}
