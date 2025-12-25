package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/what-writers-like/backend/internal/infrastructure/config"
)

type Database struct {
	db *gorm.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	dsn := cfg.DatabaseDSN
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	if err := addConstraints(db); err != nil {
		return nil, fmt.Errorf("failed to add constraints: %w", err)
	}

	return &Database{db: db}, nil
}

func addConstraints(db *gorm.DB) error {
	// Add CHECK constraint: writer_id ≠ Work.author_id
	// PostgreSQL doesn't allow subqueries in CHECK constraints directly,
	// so we use a trigger function approach
	constraintSQL := `
		CREATE OR REPLACE FUNCTION check_writer_not_author()
		RETURNS TRIGGER AS $$
		BEGIN
			IF EXISTS (
				SELECT 1 FROM works 
				WHERE id = NEW.work_id AND author_id = NEW.writer_id
			) THEN
				RAISE EXCEPTION 'writer cannot express opinion about their own work';
			END IF;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;

		DROP TRIGGER IF EXISTS trigger_check_writer_not_author ON opinions;
		CREATE TRIGGER trigger_check_writer_not_author
			BEFORE INSERT OR UPDATE ON opinions
			FOR EACH ROW
			EXECUTE FUNCTION check_writer_not_author();
	`

	if err := db.Exec(constraintSQL).Error; err != nil {
		return fmt.Errorf("failed to create constraint trigger: %w", err)
	}

	return nil
}

func (d *Database) DB() *gorm.DB {
	return d.db
}
