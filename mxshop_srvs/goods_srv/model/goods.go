package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	Id        int32     `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	DeleteAt  gorm.DeletedAt
	IsDelete  bool
}

type GormList []string

func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

type Category struct {
	BaseModel
	Name             string `gorm:"type:varchar(20);not null;"`
	Level            int32  `gorm:"type:int;not null;default:1"`
	IsTab            bool   `gorm:"type:boolean;not null;default:false"`
	ParentCategoryId int32  `gorm:"type:int"`
	ParentCategory   *Category `json:"-"`
}

type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(20);not null;"`
	Logo string `gorm:"type:varchar(200);default:''"`
}

type GoodsCategoryBrand struct {
	BaseModel
	CategoryId int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Category   Category
	BrandsId   int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Brands     Brands
}

func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);not null"`
	Index int32  `gorm:"type:int;default 1;not null"`
}

type Goods struct {
	BaseModel
	CategoryId int32 `gorm:"type:int;not null"`
	Category   Category
	BrandsId   int32 `gorm:"type:int;not null"`
	Brands     Brands
	OnSale     bool `gorm:"default:false;not null"`
	ShipFree   bool `gorm:"default:false;not null"`
	IsNew      bool `gorm:"default:false;not null"`
	IsHot      bool `gorm:"default:false;not nul"`

	Name     string `gorm:"type:varchar(200);not null"`
	ClickNum int32  `gorm:"type:int;default:0;not null"`
	SoldNum  int32  `gorm:"type:int;default:0;not null"`
	FavNum   int32  `gorm:"type:int;default:0;not null"`
	MarketPrice float32 `gorm:"not null"`
	ShopPrice float32 `gorm:"not null"`
	GoodsBrief string `gorm:"type:varchar(1000);not null"`
	Images GormList `gorm:"type:varchar(1000);not null"`
	DescImages GormList
	GoodsFrontImage string `gorm:"type:varchar(200);not null"`
}
