package pkg

type AlertHandler = func(newAlerts []Alert, allAlerts []Alert)

type AlertType int

const (
	Error           AlertType = 0
	ErrorResolved   AlertType = 1
	Started         AlertType = 7
	Running         AlertType = 2
	Down            AlertType = 3
	Away            AlertType = 4
	Warning         AlertType = 5
	WarningResolved AlertType = 6
)

type Alert struct {
	Type    AlertType
	Message string
	Key     string
}
