package monitoring

type QueryOption interface {
	Apply(*QueryOptions)
}

type QueryOptions struct {

	ResourceFilter            string
	HostName                  string
}

func NewQueryOptions() *QueryOptions {
	return &QueryOptions{}
}

type HostOption struct {
	ResourceFilter string
	HostName       string
}

func (no HostOption) Apply(o *QueryOptions) {
	o.ResourceFilter = no.ResourceFilter
	o.HostName = no.HostName
}