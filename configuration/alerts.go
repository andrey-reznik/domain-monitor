package configuration

type Alert int

const (
	Alert2Months Alert = iota
	Alert1Month
	Alert2Weeks
	Alert1Week
	Alert3Days
	AlertDaily
)

func (a Alert) String() string {
	return [...]string{"2 месяца", "1 месяц", "2 недели", "1 неделю", "3 дня", "Ежедневное оповещение"}[a]
}
