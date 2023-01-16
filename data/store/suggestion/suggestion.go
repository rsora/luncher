package suggestion

import (
	"go.uber.org/zap"
)

// Store manages the set of API's for suggestion access.
type Store struct {
	log *zap.SugaredLogger
	// db  *sqlx.DB
}

// NewStore constructs a user store for api access.
func NewStore(log *zap.SugaredLogger) Store {
	return Store{
		log: log,
	}
}

// QueryByID retrieves a suggestion by id.
func (s Store) QueryByID(suggestionID int) (Suggestion, error) {
	// len(l) = 9 ATM
	l := []Suggestion{
		{RecipeName: "Gazpacho senza cetrioli con polpa di peperoni e crostini di pane ai cereali"},
		{RecipeName: "Spaghetti grossi con sugo rosso tonno e olive"},
		{RecipeName: "Insalata di pollo"},
		{RecipeName: "Insalata di riso"},
		{RecipeName: "Insalata tonno, pomodoro e avocado"},
		{RecipeName: "Pasta con ragu'"},
		{RecipeName: "Frittata con le zucchine"},
		{RecipeName: "Lenticchie rosse e salsiccia in umido"},
		{RecipeName: "Crescia sfogliata fatta in casa con pomodoro mozzarella crudo e guacamole"},
		{RecipeName: "Spezzatino con sedano e puree di patate e zucchine"},
	}
	return l[suggestionID], nil
}
