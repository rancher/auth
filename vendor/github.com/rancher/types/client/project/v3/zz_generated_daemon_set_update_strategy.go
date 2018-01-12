package client

const (
	DaemonSetUpdateStrategyType               = "daemonSetUpdateStrategy"
	DaemonSetUpdateStrategyFieldRollingUpdate = "rollingUpdate"
	DaemonSetUpdateStrategyFieldType          = "type"
)

type DaemonSetUpdateStrategy struct {
	RollingUpdate *RollingUpdateDaemonSet `json:"rollingUpdate,omitempty"`
	Type          string                  `json:"type,omitempty"`
}
