package helm

import (
	"fmt"
)

const helmClientVersion string = "v2.12.3"
const helmArchiveTmpl string = "helm-%s-linux-amd64.tar.gz"
const helmDownloadURLTmpl string = "https://get.helm.sh/%s"

const dockerfileLines string = `RUN apt-get update && \
 apt-get install -y curl && \
 curl -o helm.tgz %s && \
 tar -xzf helm.tgz && \
 mv linux-amd64/helm /usr/local/bin && \
 rm helm.tgz
RUN helm init --client-only`

var helmArchiveVersion = fmt.Sprintf(helmArchiveTmpl, helmClientVersion)
var helmDownloadURL = fmt.Sprintf(helmDownloadURLTmpl, helmArchiveVersion)

func (m *Mixin) Build() error {
	fmt.Fprintf(m.Out, dockerfileLines, helmDownloadURL)
	return nil
}
