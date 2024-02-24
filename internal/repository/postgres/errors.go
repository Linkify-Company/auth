package postgres

const (
	ErrUniqueViolation     = "23505" // Нарушение уникальности
	ErrForeignKeyViolation = "23503" // Нарушение внешнего ключа
)
