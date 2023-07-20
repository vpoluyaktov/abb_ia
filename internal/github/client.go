package github

type GithubClient struct {
	token string
}

func NewClient(token string) (*GithubClient, error) {
	c := &GithubClient{token: token}
	return c, nil
}

func (c *GithubClient) GetLatestVer(owner string, repo string) (string, error) {
	ver := ""

	// curl -L \
	//   -H "Accept: application/vnd.github+json" \
	//   -H "Authorization: Bearer <TOKEN>" \
	//   -H "X-GitHub-Api-Version: 2022-11-28" \
	//   https://api.github.com/repos/vpoluyaktov/abb_ia/releases/latest

	return ver, nil
}
