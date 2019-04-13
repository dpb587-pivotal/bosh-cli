package gorenderer

// https://github.com/cloudfoundry/bosh/blob/master/src/bosh-director/lib/bosh/director/deployment_plan/instance_spec.rb
// TODO not renderer-specific
type TODOContext struct {
	Deployment string
	Job        string
	Index      int
	Bootstrap  bool
	Name       string
	ID         string
	AZ         string
	// Networks   NetworksTODOContext // TODO
	// 'properties_need_filtering',
	// 'dns_domain_name',
	// 'persistent_disk',
	Address      string
	IP           string
	Properties   map[string]interface{}
	ResourcePool string
}

func (c TODOContext) Property(key string) interface{} {
	if val, found := c.Properties[key]; found {
		return val
	}

	return nil
}
