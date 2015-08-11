package integration_test

import (
	"os/exec"

	"github.com/glestaris/gofigure"
	"github.com/glestaris/gofigure/providers/rabbitmq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Main", func() {
	var (
		ceBinPath string
		stdout    *gbytes.Buffer
		cmd       *exec.Cmd
	)

	JustBeforeEach(func() {
		stdout = gbytes.NewBuffer()

		cmd = exec.Command(ceBinPath)
		cmd.Stdout = stdout
		cmd.Stderr = GinkgoWriter

		Ω(cmd.Start()).Should(Succeed())
	})

	AfterEach(func() {
		Ω(cmd.Process.Kill()).Should(Succeed())
	})

	Context("when run with RabbitMQ provider", func() {
		BeforeEach(func() {
			ceBinPath = RabbitMQCeBinPath
		})

		It("works", func() {
			creds := rabbitmq.AMQPCredentials{
				Host:     "localhost",
				Port:     5672,
				Username: "guest",
				Password: "guest",
			}
			queue := "hello"

			sender, err := rabbitmq.NewSender(creds, queue)
			Ω(err).ShouldNot(HaveOccurred())

			ch, err := gofigure.OutboundChannel(sender)
			Ω(err).ShouldNot(HaveOccurred())

			ch <- "hello"
			Eventually(stdout, "2s").Should(gbytes.Say("hello"))
		})

		It("sends ints", func() {
			creds := rabbitmq.AMQPCredentials{
				Host:     "localhost",
				Port:     5672,
				Username: "guest",
				Password: "guest",
			}
			queue := "hello"

			sender, err := rabbitmq.NewSender(creds, queue)
			Ω(err).ShouldNot(HaveOccurred())

			ch, err := gofigure.OutboundChannel(sender)
			Ω(err).ShouldNot(HaveOccurred())

			ch <- 12
			Eventually(stdout, "2s").Should(gbytes.Say("12"))
		})
	})
})
