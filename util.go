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

func intP(in int) *int {
	if in == 0 {
		return nil
	}
	return &in
}
