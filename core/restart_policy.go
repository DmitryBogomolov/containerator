package core

// RestartPolicy defines container restart policy.
type RestartPolicy string

// RestartPolicy values.
const (
	RestartOnFailure     RestartPolicy = "on-failure"
	RestartUnlessStopped RestartPolicy = "unless-stopped"
	RestartAlways        RestartPolicy = "always"
)
