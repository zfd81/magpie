package cmd

type GlobalFlags struct {
	Endpoints []string
	User      string
	Password  string
}

var (
//client = http.New()
)

func url(path string) string {
	return globalFlags.Endpoints[0] + "/rock/" + path
}
