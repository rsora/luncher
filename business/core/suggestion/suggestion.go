// Package suggestion provides an example of a core business API. Right now these
// calls are just wrapping the data/store layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package suggestion

import (
	"fmt"
	"time"

	"github.com/rsora/luncher/data/store/suggestion"
	"go.uber.org/zap"
)

// Core manages the set of API's for suggestion access.
type Core struct {
	log        *zap.SugaredLogger
	suggestion suggestion.Store
}

// NewCore constructs a core for suggestion api access.
func NewCore(log *zap.SugaredLogger) Core {
	return Core{
		log:        log,
		suggestion: suggestion.NewStore(log),
	}
}

// QueryByID calculates a random ID for the suggestion list to return.
func (c Core) Query() (suggestion.Suggestion, error) {

	// PERFORM PRE BUSINESS OPERATIONS
	// hardcoded len(SuggestionList) to 9
	l := 9
	suggestionID := int(time.Now().Unix() % int64(l))

	usr, err := c.suggestion.QueryByID(suggestionID)
	if err != nil {
		return suggestion.Suggestion{}, fmt.Errorf("query: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return usr, nil
}
