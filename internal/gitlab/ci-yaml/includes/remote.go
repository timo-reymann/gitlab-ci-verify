package includes

import "gopkg.in/yaml.v3"

// RemoteInclude represents a remote include
// See more at https://docs.gitlab.com/ci/yaml/#includeremote
type RemoteInclude struct {
	Remote    string
	Integrity string
	*BaseInclude
}

func (r *RemoteInclude) Type() string {
	return "remote"
}

func (r *RemoteInclude) Equals(i Include) bool {
	remoteInclude, ok := i.(*RemoteInclude)
	if !ok {
		return false
	}

	return remoteInclude.Remote == r.Remote
}

func NewRemoteInclude(node *yaml.Node, remote, integrity string) *RemoteInclude {
	return &RemoteInclude{
		Remote:      remote,
		Integrity:   integrity,
		BaseInclude: NewBaseInclude(node),
	}
}
