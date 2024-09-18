package schedules

type Request interface {
	ScheduleName() string
}

type requestImpl struct {
	scheduleName string
}

func (i *requestImpl) ScheduleName() string {
	return i.scheduleName
}
