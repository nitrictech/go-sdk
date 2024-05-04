package workers

import "context"

type Worker interface {
	Start(context.Context) error
}
