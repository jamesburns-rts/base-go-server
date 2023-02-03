package vm

type (
	Versions struct {
		Version string     `json:"version"`
		Go      string     `json:"go"`
		Echo    string     `json:"echo"`
		Git     GitVersion `json:"git"`
	}

	GitVersion struct {
		Commit string `json:"commit"`
		Branch string `json:"branch"`
		State  string `json:"state"`
	}
)
