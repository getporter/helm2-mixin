module get.porter.sh/mixin/helm

go 1.13

require (
	get.porter.sh/porter v0.25.0-beta.1
	github.com/Masterminds/semver v1.5.0
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/gobuffalo/packr/v2 v2.7.1
	github.com/hashicorp/go-multierror v1.0.0
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/nwaples/rardecode v1.0.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.4.0
	github.com/xeipuuv/gojsonschema v1.2.0
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	gopkg.in/yaml.v2 v2.2.4
	k8s.io/apimachinery v0.0.0-20191004115801-a2eda9f80ab8
	k8s.io/client-go v0.0.0-20191016111102-bec269661e48
)

replace github.com/hashicorp/go-plugin => github.com/carolynvs/go-plugin v1.0.1-acceptstdin
