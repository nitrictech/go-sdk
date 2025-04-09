package lifecycle_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/nitrictech/go-sdk/nitric/lifecycle"
)

var _ = Describe("Lifecycle", func() {
	BeforeEach(func() {
		// Clear environment before each test
		os.Unsetenv(lifecycle.NITRIC_ENVIRONMENT)
	})

	Describe("GetCurrentLifecycle", func() {
		Context("when environment variable is set to valid stages", func() {
			It("should return LocalRun stage", func() {
				os.Setenv(lifecycle.NITRIC_ENVIRONMENT, "run")
				stage, err := lifecycle.GetCurrentLifecycle()
				Expect(err).NotTo(HaveOccurred())
				Expect(stage).To(Equal(lifecycle.LocalRun))
			})

			It("should return Build stage", func() {
				os.Setenv(lifecycle.NITRIC_ENVIRONMENT, "build")
				stage, err := lifecycle.GetCurrentLifecycle()
				Expect(err).NotTo(HaveOccurred())
				Expect(stage).To(Equal(lifecycle.Build))
			})

			It("should return Cloud stage", func() {
				os.Setenv(lifecycle.NITRIC_ENVIRONMENT, "cloud")
				stage, err := lifecycle.GetCurrentLifecycle()
				Expect(err).NotTo(HaveOccurred())
				Expect(stage).To(Equal(lifecycle.Cloud))
			})
		})

		Context("when environment variable is not set", func() {
			It("should return an error", func() {
				stage, err := lifecycle.GetCurrentLifecycle()
				Expect(err).To(HaveOccurred())
				Expect(stage).To(BeEmpty())
			})
		})

		Context("when environment variable is set to invalid stage", func() {
			It("should return an error", func() {
				os.Setenv(lifecycle.NITRIC_ENVIRONMENT, "invalid")
				stage, err := lifecycle.GetCurrentLifecycle()
				Expect(err).To(HaveOccurred())
				Expect(stage).To(BeEmpty())
			})
		})
	})

	Describe("IsInLifecycle", func() {
		Context("when current stage matches one of the provided stages", func() {
			It("should return true for LocalRun", func() {
				os.Setenv(lifecycle.NITRIC_ENVIRONMENT, "run")
				result := lifecycle.IsInLifecycle(lifecycle.LocalRun, lifecycle.Cloud)
				Expect(result).To(BeTrue())
			})

			It("should return true for Cloud", func() {
				os.Setenv(lifecycle.NITRIC_ENVIRONMENT, "cloud")
				result := lifecycle.IsInLifecycle(lifecycle.LocalRun, lifecycle.Cloud)
				Expect(result).To(BeTrue())
			})
		})

		Context("when current stage does not match any of the provided stages", func() {
			It("should return false", func() {
				os.Setenv(lifecycle.NITRIC_ENVIRONMENT, "build")
				result := lifecycle.IsInLifecycle(lifecycle.LocalRun, lifecycle.Cloud)
				Expect(result).To(BeFalse())
			})
		})
	})

	Describe("WhenInLifecycles", func() {
		Context("when current stage matches one of the provided stages", func() {
			It("should execute the callback", func() {
				os.Setenv(lifecycle.NITRIC_ENVIRONMENT, "run")
				result := lifecycle.WhenInLifecycles([]lifecycle.LifecycleStage{lifecycle.LocalRun, lifecycle.Cloud}, func() string {
					return "executed"
				})
				Expect(result).To(Equal("executed"))
			})
		})

		Context("when current stage does not match any of the provided stages", func() {
			It("should not execute the callback and return zero value", func() {
				os.Setenv(lifecycle.NITRIC_ENVIRONMENT, "build")
				result := lifecycle.WhenInLifecycles([]lifecycle.LifecycleStage{lifecycle.LocalRun, lifecycle.Cloud}, func() string {
					return "executed"
				})
				Expect(result).To(BeEmpty())
			})
		})
	})

	Describe("IsRunning", func() {
		Context("when in LocalRun stage", func() {
			It("should return true", func() {
				os.Setenv(lifecycle.NITRIC_ENVIRONMENT, "run")
				Expect(lifecycle.IsRunning()).To(BeTrue())
			})
		})

		Context("when in Cloud stage", func() {
			It("should return true", func() {
				os.Setenv(lifecycle.NITRIC_ENVIRONMENT, "cloud")
				Expect(lifecycle.IsRunning()).To(BeTrue())
			})
		})

		Context("when in Build stage", func() {
			It("should return false", func() {
				os.Setenv(lifecycle.NITRIC_ENVIRONMENT, "build")
				Expect(lifecycle.IsRunning()).To(BeFalse())
			})
		})
	})

	Describe("IsCollecting", func() {
		Context("when in Build stage", func() {
			It("should return true", func() {
				os.Setenv(lifecycle.NITRIC_ENVIRONMENT, "build")
				Expect(lifecycle.IsCollecting()).To(BeTrue())
			})
		})

		Context("when in LocalRun stage", func() {
			It("should return false", func() {
				os.Setenv(lifecycle.NITRIC_ENVIRONMENT, "run")
				Expect(lifecycle.IsCollecting()).To(BeFalse())
			})
		})

		Context("when in Cloud stage", func() {
			It("should return false", func() {
				os.Setenv(lifecycle.NITRIC_ENVIRONMENT, "cloud")
				Expect(lifecycle.IsCollecting()).To(BeFalse())
			})
		})
	})
})
