//go:generate mockgen -destination=internal/dao/mocks_test.go -package=dao github.com/goverland-labs/goverland-core-storage/internal/dao DataProvider,Publisher,DaoIDProvider

package main
