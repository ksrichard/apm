package project

type ProjectDetails struct {
	Board        *ProjectBoard        `json:"board"`
	Dependencies []ProjectDependency `json:"dependencies"`
}

type ProjectBoard struct {
	Package      string `json:"package,omitempty"`
	Architecture string `json:"architecture,omitempty"`
	Version      string `json:"version,omitempty"`
	BoardManagerUrl string `json:"board_manager_url,omitempty"`
}

type ProjectDependency struct {
	Library string `json:"library,omitempty"`
	Version string `json:"version,omitempty"`
	Git     string `json:"git,omitempty"`
	Zip     string `json:"zip,omitempty"`
}
