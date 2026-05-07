package domain

type Product struct {
	DBID         int     `json:"-" gorm:"column:id;primaryKey;autoIncrement"`
	ID           int     `json:"id" gorm:"column:wb_id;uniqueIndex;not null"`
	Brand        string  `json:"brand" gorm:"column:brand;not null"`
	Name         string  `json:"name" gorm:"column:name;not null"`
	SupplierId   int     `json:"supplierId" gorm:"column:supplierId;not null"`
	ReviewRating float64 `json:"reviewRating" gorm:"column:reviewRating;not null"`
	FeedBacks    int     `json:"feedbacks" gorm:"column:feedbacks;not null"`
	Pics         int     `json:"pics" gorm:"column:pics;not null"`
	Sizes        []Size  `json:"sizes" gorm:"foreignKey:ProductID"`
}

type Size struct {
	DBID      int    `json:"-" gorm:"column:id;primaryKey;autoIncrement"`
	ProductID int    `json:"-" gorm:"column:product_id"`
	Name      string `json:"name" gorm:"column:name"`

	Price struct {
		Basic   int `json:"basic" gorm:"-"`
		Product int `json:"product" gorm:"-"`
	} `json:"price" gorm:"-"`

	PriceBasic   float64 `json:"-" gorm:"column:priceBasic;not null"`
	PriceProduct float64 `json:"-" gorm:"column:priceProduct;not null"`
}
