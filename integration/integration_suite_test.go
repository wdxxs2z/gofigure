package integration_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var RabbitMQCeBinPath string

func TestIntegration(t *testing.T) {
	var beforeSuite struct {
		RabbitMQCeBinPath string
	}

	SynchronizedBeforeSuite(func() []byte {
		var err error

		beforeSuite.RabbitMQCeBinPath, err = gexec.Build("github.com/glestaris/gofigure/integration/rabbitmq_chan_echo")
		Ω(err).ShouldNot(HaveOccurred())

		b, err := json.Marshal(beforeSuite)
		Expect(err).ToNot(HaveOccurred())

		return b
	}, func(data []byte) {
		err := json.Unmarshal(data, &beforeSuite)
		Expect(err).ToNot(HaveOccurred())

		RabbitMQCeBinPath = beforeSuite.RabbitMQCeBinPath
		Expect(RabbitMQCeBinPath).NotTo(BeEmpty())
	})

	SynchronizedAfterSuite(func() {
		//noop
	}, func() {
		gexec.CleanupBuildArtifacts()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}
