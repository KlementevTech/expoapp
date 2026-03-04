package storage

import (
	"context"
	"time"

	"expo/internal/api/handler"
	"expo/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

// PartStorageSuite определяет структуру сюиты.
// Мы не создаем её напрямую, а встраиваем в тесты конкретных реализаций.
type PartStorageSuite struct {
	suite.Suite

	Repo handler.PartRepository
}

// TestCreateAndFindByID Все методы, начинающиеся с "Test", будут запущены автоматически.
func (s *PartStorageSuite) TestCreateAndFindByID() {
	ctx := context.Background()
	newPart := model.NewPart(newPartID(), "Engine")

	// Create
	err := s.Repo.Create(ctx, newPart)
	s.Require().NoError(err)

	// Find
	found, err := s.Repo.Get(ctx, newPart.ID)
	s.Require().NoError(err)
	s.Equal("Engine", found.Name)

	// NotFound
	_, err = s.Repo.Get(ctx, uuid.New())
	s.Require().ErrorIs(err, model.ErrNotFound)
}

func (s *PartStorageSuite) TestDuplicateID() {
	ctx := context.Background()
	p := model.NewPart(newPartID(), "Unique")

	_ = s.Repo.Create(ctx, p)
	err := s.Repo.Create(ctx, p)

	s.Require().ErrorIs(err, model.ErrAlreadyExists)
}

func (s *PartStorageSuite) TestGetListSorting() {
	ctx := context.Background()
	baseTime := time.Now()
	p1 := &model.Part{ID: newPartID(), Name: "Bolt", Version: 1, CreatedAt: baseTime}
	p2 := &model.Part{ID: newPartID(), Name: "Bar", Version: 1, CreatedAt: baseTime.Add(time.Second)}
	p3 := &model.Part{ID: newPartID(), Name: "Nut", Version: 1, CreatedAt: baseTime.Add(time.Minute)}

	s.NoError(s.Repo.Create(ctx, p1))
	s.NoError(s.Repo.Create(ctx, p2))
	s.NoError(s.Repo.Create(ctx, p3))

	list, err := s.Repo.List(ctx)
	s.Require().NoError(err)
	const partsNumber = 3
	s.Len(list, partsNumber)
	s.Equal("Nut", list[0].Name)
	s.Equal("Bar", list[1].Name)
	s.Equal("Bolt", list[2].Name)
}

func (s *PartStorageSuite) TestContextCancellation() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	p := model.NewPart(newPartID(), "Bud")
	err := s.Repo.Create(ctx, p)

	// Require() остановит тест здесь, если err == nil
	s.Require().Error(err, "Репозиторий должен вернуть ошибку при отмененном контексте")

	// Теперь безопасно проверять тип ошибки
	s.ErrorIs(err, context.Canceled)
}

func newPartID() model.PartID {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return id
}
