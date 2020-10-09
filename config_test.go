package proxanne_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/waaaaargh/proxanne"
)

var _ = Describe("ParseConfig", func() {
	When("presented wih a syntactically invalid configuration file", func() {
		configBytes := []byte("@NotYAMLQQQQQQ")
		It("should throw an error", func() {
			_, err := proxanne.ParseConfig(configBytes)
			Expect(err).ToNot(BeNil())
		})
	})

	When("presented with a syntactically valid configuration file", func() {
		var (
			err    error
			config *proxanne.Config
		)

		configBytes := []byte(`
routes:
- matches: "foo/bar/"
  target: "https://example.com/foo"
- matches: "baz/qux/"
  target: "https://example.com/baz"`)

		BeforeEach(func() {
			config, err = proxanne.ParseConfig(configBytes)
		})

		It("should return a configuration object", func() {
			Expect(config).ToNot(BeNil())
		})

		It("should not throw an error", func() {
			Expect(err).To(BeNil())
		})
	})
})

var _ = Describe("BuildRouter", func() {
	When("presented with a syntactically valid, but semantically invalid configuration", func() {
		config := &proxanne.Config{}

		var (
			err    error
			router proxanne.Router
		)

		BeforeEach(func() {
			router, err = proxanne.BuildRouter(config)
		})

		It("should emit an error", func() {
			Expect(err).ToNot(BeNil())
		})

		It("should emit an emptry router", func() {
			Expect(router).To(BeNil())
		})
	})

	DescribeTable("Catching invalid Configurations",
		func(config *proxanne.Config) {
			_, err := proxanne.BuildRouter(config)
			Expect(err).ToNot(BeNil())
		},
		Entry(
			"invalid regexp",
			&proxanne.Config{[]struct {
				Matches string `yaml:"matches"`
				Target  string `yaml:"target"`
			}{{"([", "foobar"}}},
		),
		Entry(
			"invalid target url",
			&proxanne.Config{[]struct {
				Matches string `yaml:"matches"`
				Target  string `yaml:"target"`
			}{{"foobar", ":"}}},
		),
	)

	When("presented with a valid configuration", func() {
		config := &proxanne.Config{[]struct {
			Matches string `yaml:"matches"`
			Target  string `yaml:"target"`
		}{{"/foobar/", "example.com/foobar"}}}

		var (
			router proxanne.Router
			err    error
		)

		BeforeEach(func() {
			router, err = proxanne.BuildRouter(config)
		})

		It("should return a Router", func() {
			Expect(len(router)).To(Equal(1))
		})

		It("should not return an error", func() {
			Expect(err).To(BeNil())
		})
	})
})
