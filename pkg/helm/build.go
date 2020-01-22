package helm

import (
	"fmt"
)

// These values may be referenced elsewhere (init.go), hence consts
const helmClientVersion string = "v2.15.2"
const helmArchiveTmpl string = "helm-%s-linux-amd64.tar.gz"
const helmDownloadURLTmpl string = "https://get.helm.sh/%s"

const getHelm string = `RUN apt-get update && \
 apt-get install -y curl && \
 curl -o helm.tgz %s && \
 tar -xzf helm.tgz && \
 mv linux-amd64/helm /usr/local/bin && \
 rm helm.tgz
RUN helm init --client-only
`

var helmArchiveVersion = fmt.Sprintf(helmArchiveTmpl, helmClientVersion)
var helmDownloadURL = fmt.Sprintf(helmDownloadURLTmpl, helmArchiveVersion)

// kubectl may be necessary; for example, to set up RBAC for Helm's Tiller component if needed
const kubeVersion string = "v1.15.3"
const getKubectl string = `RUN apt-get update && \
 apt-get install -y apt-transport-https curl && \
 curl -o kubectl https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/amd64/kubectl && \
 mv kubectl /usr/local/bin && \
 chmod a+x /usr/local/bin/kubectl`

func (m *Mixin) Build() error {
	fmt.Fprintf(m.Out, getHelm, helmDownloadURL)
	fmt.Fprintf(m.Out, getKubectl, kubeVersion)
	return nil
}
