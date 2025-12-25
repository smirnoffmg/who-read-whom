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

	if err := enableExtensions(db); err != nil {
		return nil, fmt.Errorf("failed to enable extensions: %w", err)
	}

	if err := addConstraints(db); err != nil {
		return nil, fmt.Errorf("failed to add constraints: %w", err)
	}

	if err := createSearchIndexes(db); err != nil {
		return nil, fmt.Errorf("failed to create search indexes: %w", err)
	}

	return &Database{db: db}, nil
}

func enableExtensions(db *gorm.DB) error {
	// Enable pg_trgm extension for fuzzy search
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm").Error; err != nil {
		return fmt.Errorf("failed to enable pg_trgm extension: %w", err)
	}
	return nil
}

func createSearchIndexes(db *gorm.DB) error {
	// Create GIN indexes for fuzzy search on writers table
	indexesSQL := `
		CREATE INDEX IF NOT EXISTS idx_writers_name_trgm ON writers USING gin(name gin_trgm_ops);
		CREATE INDEX IF NOT EXISTS idx_writers_bio_trgm ON writers USING gin(bio gin_trgm_ops) WHERE bio IS NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_works_title_trgm ON works USING gin(title gin_trgm_ops);
	`

	if err := db.Exec(indexesSQL).Error; err != nil {
		return fmt.Errorf("failed to create search indexes: %w", err)
	}

	return nil
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
