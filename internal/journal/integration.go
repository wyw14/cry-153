package journal

import (
	"fmt"
	"github.com/wyw14/cry-153/internal/model"
)

func DescribeOperation(op model.Operation) string {
	return fmt.Sprintf("%s:%s:%s", op.Kind, op.IncidentID, op.Revision)
}
func PersistOperation(s *Store, op model.Operation) error {
	return s.Append(Timeline("operation", DescribeOperation(op)))
}
