package listproblemdifficulty

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"

	"emperror.dev/errors"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db: db,
	}
}

func (g *GormRepository) GetAllProblemDifficulties(ctx context.Context) ([]ProblemDifficulty, error) {
	var problemDifficulties []database.ProblemDifficulty
	if err := g.db.WithContext(ctx).Model(&database.ProblemDifficulty{}).Preload("DisplayNames").Find(&problemDifficulties).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get all problem difficulties")
	}

	result := make([]ProblemDifficulty, 0, len(problemDifficulties))
	for _, problemDifficulty := range problemDifficulties {
		pd := ProblemDifficulty{
			ProblemDifficultyID: problemDifficulty.ProblemDifficultyID,
			DisplayNames:        make([]DisplayName, 0, len(problemDifficulty.DisplayNames)),
		}

		for _, displayName := range problemDifficulty.DisplayNames {
			pd.DisplayNames = append(pd.DisplayNames, DisplayName{
				Language: displayName.Language,
				Name:     displayName.DisplayName,
			})
		}

		result = append(result, pd)
	}

	return result, nil
}
