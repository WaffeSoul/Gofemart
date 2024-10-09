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
