package internal

type Credentials struct {
	ClientID string
	ClientSecret string
}

type Item struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description"`
	PictureURL  string  `json:"picture_url"`
	Quantity    int     `json:"quantity" validate:"required"`
	UnitPrice   float64 `json:"unit_price" validate:"required"`
}

type Payer struct {
	Name    string `json:"name" validate:"required"`
	Surname string `json:"surname"`
	Email   string `json:"email" validate:"required"`
	Phone Phone `json:"phone" validate:"required"`
	Address Address `json:"address" validate:"required"`
	CreatedAt string `json:"date_created" validate:"required"`
}

type Phone struct {
	AreaCode string `json:"area_code"`
	Number   string `json:"number" validate:"required"`
}

type Address struct {
	ZipCode string `json:"zip_code"`
	Street  string `json:"street" validate:"required"`
	Number  int    `json:"number" validate:"required"`
}

type Redirect struct {
	Success string `json:"success"`
	Pending string `json:"pending"`
	failure string `json:"failure"`
}

type NewPreference struct {
	Items []Item `json:"items" validate:"required,min=1"`
	Payer Payer `json:"payer" validate:"required"`
	Redirect Redirect `json:"back_urls"`
	AutoReturn bool `json:"auto_return"`
}

