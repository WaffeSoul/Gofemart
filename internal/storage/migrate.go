package storage

import "fmt"

func MigrateTables(s Store) error {
	if err := s.Users().Migrate(); err != nil {
		return fmt.Errorf("failed to migrate users: %w", err)
	}
	if err := s.Orders().Migrate(); err != nil {
		return fmt.Errorf("failed to migrate orders: %w", err)
	}
	if err := s.Withdrawals().Migrate(); err != nil {
		return fmt.Errorf("failed to migrate withdrawls: %w", err)
	}
	return nil
}

func DropTable(s Store) error {
	if err := s.Users().Drop(); err != nil {
		return fmt.Errorf("failed to drop users: %w", err)
	}
	if err := s.Orders().Drop(); err != nil {
		return fmt.Errorf("failed to drop orders: %w", err)
	}
	if err := s.Withdrawals().Drop(); err != nil {
		return fmt.Errorf("failed to drop withdrawls: %w", err)
	}
	return nil
}