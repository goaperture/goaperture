package client

type Permissions []string

func (p *Permissions) Has(permissions ...string) bool {
	if p == nil {
		return false
	}

	for _, sysPerm := range permissions {
		pass := false

		for _, userPerm := range *p {
			if sysPerm == userPerm {
				pass = true
				break
			}
		}

		if !pass {
			return false
		}
	}

	return true
}

func (p *Permissions) Some(permissions ...string) bool {
	if p == nil {
		return false
	}

	for _, sysPerm := range permissions {
		for _, userPerm := range *p {
			if sysPerm == userPerm {
				return true
			}
		}
	}

	return false
}
