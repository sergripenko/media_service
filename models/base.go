package models

import "time"

type Base struct {
	CreatedAt time.Time `orm:"column(created_at);auto_now_add;type(timestamp without time zone);null" json:"created_at"`
	UpdatedAt time.Time `orm:"column(updated_at);auto_now;type(timestamp without time zone);null" json:"updated_at"`
	DeletedAt time.Time `orm:"column(deleted_at);type(timestamp without time zone);null" json:"deleted_at"`
}
