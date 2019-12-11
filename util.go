package ign2to3

func strP(in string) *string {
	if in == "" {
		return nil
	}
	return &in
}

func boolP(in bool) *bool {
	if !in {
		return nil
	}
	return &in
}
