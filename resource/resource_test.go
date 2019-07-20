package resource_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	oc "github.com/cloudboss/ofcourse/ofcourse"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/terminatingcode/exhibit-resource/resource"
)

var (
	testLogger = oc.NewLogger(oc.SilentLevel)
	env        = oc.NewEnvironment()
	r          = Resource{}
)

var _ = Describe("Resource", func() {
	Describe("Check", func() {
		Context("with no previous version", func() {
			It("should increment", func() {
				versions, err := r.Check(oc.Source{}, nil, env, testLogger)
				Expect(err).ToNot(HaveOccurred())
				Expect(versions).To(Equal([]oc.Version{oc.Version{"count": "1"}}))
			})
		})

		Context("with a previous version", func() {
			It("should increment", func() {
				versions, err := r.Check(oc.Source{}, oc.Version{"count": "1000"}, env, testLogger)
				Expect(err).ToNot(HaveOccurred())
				Expect(versions).To(Equal([]oc.Version{oc.Version{"count": "1001"}}))
			})
		})

		Context("with an invalid version", func() {
			It("should error", func() {
				versions, err := r.Check(oc.Source{}, oc.Version{"numero": "1"}, env, testLogger)
				Expect(err).To(HaveOccurred())
				Expect(versions).To(BeNil())
			})
		})
	})

	Describe("In", func() {
		var (
			td  string
			err error
		)

		BeforeEach(func() {
			td, err = ioutil.TempDir("", "resource-")
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			defer os.RemoveAll(td)
		})

		Context("with a valid input", func() {
			It("should return metatdata", func() {
				version := oc.Version{"count": "1"}
				_, metadata, err := r.In(td, oc.Source{}, oc.Params{}, version, env, testLogger)
				Expect(err).ToNot(HaveOccurred())
				Expect(metadata).To(ConsistOf(oc.Metadata{{Name: "a", Value: "b"}, {Name: "c", Value: "d"}}))

				path := fmt.Sprintf("%s/version", td)
				_, fileErr := os.Stat(path)
				Expect(os.IsNotExist(fileErr)).To(BeFalse())

				bytes, err := ioutil.ReadFile(path)
				Expect(err).ToNot(HaveOccurred())

				var readVersion oc.Version
				err = json.Unmarshal(bytes, &readVersion)
				Expect(readVersion).To(Equal(version))
			})
		})
	})

	Describe("Out", func() {
		var (
			td  string
			err error
		)

		BeforeEach(func() {
			td, err = ioutil.TempDir("", "resource-")
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			defer os.RemoveAll(td)
		})

		Context("with a valid version path", func() {
			It("should return metatdata", func() {
				version := oc.Version{"count": "1"}

				versionPath := "path/version"
				fullVersionPath := fmt.Sprintf("%s/%s", td, versionPath)
				err = os.MkdirAll(filepath.Dir(fullVersionPath), 0777)
				Expect(err).ToNot(HaveOccurred())

				contents, err := json.Marshal(version)
				Expect(err).ToNot(HaveOccurred())

				err = ioutil.WriteFile(fullVersionPath, []byte(contents), 0666)
				Expect(err).ToNot(HaveOccurred())

				versionOut, metadata, err := r.Out(td, oc.Source{}, oc.Params{"version_path": versionPath}, env, testLogger)
				Expect(err).ToNot(HaveOccurred())
				Expect(version).To(Equal(versionOut))
				Expect(metadata).To(ConsistOf(oc.Metadata{}))
			})
		})

		Context("with an invalid version path", func() {
			It("should error", func() {
				versionOut, metadata, err := r.Out(td, oc.Source{}, oc.Params{}, env, testLogger)
				Expect(err).To(HaveOccurred())
				Expect(versionOut).To(BeNil())
				Expect(metadata).To(BeNil())
			})
		})
	})
})
