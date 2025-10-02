package constants

const (
	DefaultPage         = 1
	DefaultPageSize     = 10
	MaxPageSize         = 100
	MinPageSize         = 1
	ShutdownTimeout     = 30 
	DefaultGRPCPort     = "50051"
	DefaultDBDriver     = "postgres"
	DefaultDBName       = "products.db"
	DefaultDBHost       = "localhost"
	DefaultDBPort       = "5432"
	DefaultDBUser       = "postgres"
	DefaultDBPassword   = "postgres"
	DefaultDBSSLMode    = "disable"
)

const (
	ErrProductNameRequired     = "product name is required"
	ErrPriceNegative          = "price cannot be negative"
	ErrProductTypeRequired    = "product type is required"
	ErrInvalidProductID       = "invalid product ID format"
	ErrProductNotFound        = "product not found"
	ErrPlanNameRequired       = "plan name is required"
	ErrDurationPositive       = "duration must be positive"
	ErrInvalidPlanID          = "invalid subscription plan ID format"
	ErrPlanNotFound           = "subscription plan not found"
)

