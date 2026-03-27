module github.com/coreos/ign-converter

go 1.24.0

replace github.com/coreos/ign-converter/translate/v34tov35 => ./translate/v35tov34

require (
	github.com/clarketm/json v1.17.1
	github.com/coreos/ignition v0.35.0
	github.com/coreos/ignition/v2 v2.26.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/ajeddeloh/go-json v0.0.0-20200220154158-5ae607161559 // indirect
	github.com/aws/aws-sdk-go-v2 v1.41.1 // indirect
	github.com/coreos/go-json v0.0.0-20230131223807-18775e0fb4fb // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd v0.0.0-20181031085051-9002847aa142 // indirect
	github.com/coreos/go-systemd/v22 v22.6.0 // indirect
	github.com/coreos/vcontext v0.0.0-20230201181013-d72178a18687 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/vincent-petithory/dataurl v1.0.0 // indirect
	go4.org v0.0.0-20200104003542-c7e774b10ea0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
