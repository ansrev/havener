// Copyright © 2018 The Havener
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package havener_test

import (
	"os"
	"path/filepath"

	"github.com/homeport/havener/pkg/havener"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const helmChartTempLocation = "/.havener/chartstore/chartroom/charts/tomcat/stable"
const fakeChartTempLocation = "/fake-chart"
const wrongChartTempLocation = "stable/wrong-chart"
const goodChartTempLocation = "tomcat/stable"

var _ = Describe("Charts Locations", func() {

	BeforeEach(func() {
		os.Setenv("HOME", "/tmp")
	})

	Describe("loading a Helm Chart", func() {
		Context("when the Helm Chart exists locally", func() {
			var fakeChartLocation string

			JustBeforeEach(func() {
				fakeChartLocation = os.Getenv("HOME") + fakeChartTempLocation
			})

			It("should find a helm chart when the path exist", func() {
				os.Mkdir(fakeChartLocation, os.ModePerm)
				os.Mkdir(fakeChartLocation+"/templates", os.ModePerm)
				os.Create(fakeChartLocation + "/Chart.yaml")
				os.Create(fakeChartLocation + "/values.yaml")
				helmChartPath, err := havener.PathToHelmChart(fakeChartLocation)
				if err != nil {
					havener.ExitWithError("Unable to locate helm chart location", err)
				}
				Expect(err).Should(BeNil())
				absolutePath, _ := filepath.Abs(fakeChartLocation)
				Expect(helmChartPath).Should(Equal(absolutePath))
			})
		})
		Context("when the Helm Chart exists remotely", func() {
			It("should find a remote chart and place in under the ~/.havener repo", func() {
				remoteHelmChartPath, err := havener.PathToHelmChart(goodChartTempLocation)
				if err != nil {
					havener.ExitWithError("Unable to locate helm chart location", err)
				}
				Expect(err).Should(BeNil())
				absoluteChartPath, _ := filepath.Abs(os.Getenv("HOME") + helmChartTempLocation)
				Expect(remoteHelmChartPath).Should(Equal(absoluteChartPath))
			})
		})
		Context("when the Helm Chart does not exists", func() {
			It("should exit with error when the helm chart path does not exist", func() {
				_, err := havener.PathToHelmChart(os.Getenv("HOME") + wrongChartTempLocation)
				constructedError := "unable to determine Helm Chart location of '" + os.Getenv("HOME") + wrongChartTempLocation +
					"', it is not a local path, nor is it defined in https://github.com/homeport/chartroom or https://github.com/helm/charts"
				Expect(err.Error()).Should(Equal(constructedError))
			})
		})
	})
})

var _ = AfterSuite(func() {
	os.RemoveAll(os.Getenv("HOME") + "/.havener")
	os.RemoveAll(os.Getenv("HOME") + "/fake-chart")
	os.Unsetenv("HOME")
})