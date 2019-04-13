package localvm

type CloudStemcell struct{}

func (cs CloudStemcell) CID() string {
	return "localvm"
}

func (cs CloudStemcell) Name() string {
	return "localvm"
}

func (cs CloudStemcell) Version() string {
	return "0.0.0"
}

func (cs CloudStemcell) PromoteAsCurrent() error {
	return nil
}

func (cs CloudStemcell) Delete() error {
	return nil
}
