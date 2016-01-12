package main

import (
	"testing"

	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/stretchr/testify/suite"
)

type StatusHandlerSuite struct {
	*TestSuite
}

func TestStatusHandlerSuite(t *testing.T) {
	suiteTester := &StatusHandlerSuite{&TestSuite{}}
	suite.Run(t, suiteTester)
}

func (s *StatusHandlerSuite) TestPullParallelDownloads() {
	testSteps := []struct {
		in       *jsonmessage.JSONMessage
		expected string
	}{
		{
			&jsonmessage.JSONMessage{
				ID:     "ubuntu:latest",
				Status: "The image you are pulling has been verified",
			},
			"The image you are pulling has been verified: ubuntu:latest\n",
		},
		{
			&jsonmessage.JSONMessage{
				ID:       "511136ea3c5a",
				Status:   "Pulling fs layer",
				Progress: &jsonmessage.JSONProgress{Current: 0, Start: 0, Total: 0},
			},
			"Pulling fs layer: 511136ea3c5a\n",
		},
		{
			&jsonmessage.JSONMessage{
				ID:       "c7b7c6419568",
				Status:   "Pulling fs layer",
				Progress: &jsonmessage.JSONProgress{Current: 0, Start: 0, Total: 0},
			},
			"Pulling fs layer: c7b7c6419568\n",
		},
		{
			&jsonmessage.JSONMessage{
				ID:       "511136ea3c5a",
				Status:   "Downloading",
				Progress: &jsonmessage.JSONProgress{Current: 0, Start: 0, Total: 100},
			},
			"Downloading: 511136ea3c5a (0%)",
		},
		{
			&jsonmessage.JSONMessage{
				ID:       "511136ea3c5a",
				Status:   "Downloading",
				Progress: &jsonmessage.JSONProgress{Current: 50, Start: 0, Total: 100},
			},
			"\rDownloading: 511136ea3c5a (50%)",
		},
		{
			&jsonmessage.JSONMessage{
				ID:       "c7b7c6419568",
				Status:   "Downloading",
				Progress: &jsonmessage.JSONProgress{Current: 0, Start: 0, Total: 100},
			},
			"\rDownloading: 511136ea3c5a (50%), Downloading: c7b7c6419568 (0%)",
		},
		{
			&jsonmessage.JSONMessage{
				ID:       "511136ea3c5a",
				Status:   "Download complete",
				Progress: &jsonmessage.JSONProgress{Current: 0, Start: 0, Total: 0},
			},
			"\rDownload complete: 511136ea3c5a                                \nDownloading: c7b7c6419568 (0%)",
		},
		{
			&jsonmessage.JSONMessage{
				ID:       "c7b7c6419568",
				Status:   "Downloading",
				Progress: &jsonmessage.JSONProgress{Current: 50, Start: 0, Total: 100},
			},
			"\rDownloading: c7b7c6419568 (50%)",
		},
		{
			&jsonmessage.JSONMessage{
				ID:       "c7b7c6419568",
				Status:   "Download complete",
				Progress: &jsonmessage.JSONProgress{Current: 0, Start: 0, Total: 0},
			},
			"\rDownload complete: c7b7c6419568\n",
		},
		{
			&jsonmessage.JSONMessage{
				ID:       "511136ea3c5a",
				Status:   "Extracting",
				Progress: &jsonmessage.JSONProgress{Current: 10, Start: 0, Total: 100},
			},
			"Extracting: 511136ea3c5a (10%)",
		},
		{
			&jsonmessage.JSONMessage{
				ID:       "511136ea3c5a",
				Status:   "Pull complete",
				Progress: &jsonmessage.JSONProgress{Current: 0, Start: 0, Total: 0},
			},
			"\rPull complete: 511136ea3c5a   \n",
		},
		{
			&jsonmessage.JSONMessage{
				ID:       "c7b7c6419568",
				Status:   "Extracting",
				Progress: &jsonmessage.JSONProgress{Current: 55, Start: 0, Total: 100},
			},
			"Extracting: c7b7c6419568 (55%)",
		},
		{
			&jsonmessage.JSONMessage{
				ID:       "c7b7c6419568",
				Status:   "Pull complete",
				Progress: &jsonmessage.JSONProgress{Current: 0, Start: 0, Total: 0},
			},
			"\rPull complete: c7b7c6419568   \n",
		},
		{
			&jsonmessage.JSONMessage{
				Status: "Status: Downloaded newer image for ubuntu:latest;",
			},
			"Status: Downloaded newer image for ubuntu:latest;\n",
		},
	}

	p := NewJSONMessageProcessor()
	for _, step := range testSteps {
		actual := p.ProcessJSONMessage(step.in)
		s.Equal(actual, step.expected)
	}
}

func (s *StatusHandlerSuite) TestPushParallelUploads() {
	testSteps := []struct {
		in       *jsonmessage.JSONMessage
		expected string
	}{
		{
			&jsonmessage.JSONMessage{
				Status: "The push refers to a repository [127.0.0.1:3000/bvdberg/pass] (len: 1)",
			},
			"Pushing to registry\n",
		},
		{
			&jsonmessage.JSONMessage{
				Status: "Sending image list",
			},
			"Sending image list\n",
		},
		{
			&jsonmessage.JSONMessage{
				Status: "Pushing repository 127.0.0.1:3000/bvdberg/pass (1 tags)",
			},
			"Pushing 1 tag(s)\n", // TODO
		},
		{
			&jsonmessage.JSONMessage{
				ID:       "511136ea3c5a",
				Status:   "Pushing",
				Progress: &jsonmessage.JSONProgress{Current: 0, Start: 0, Total: 0},
			},
			"Pushing: 511136ea3c5a",
		},
		{
			&jsonmessage.JSONMessage{
				ID:       "511136ea3c5a",
				Status:   "Buffering to disk",
				Progress: &jsonmessage.JSONProgress{Current: 10, Start: 0, Total: 0},
			},
			"\rBuffering to disk: 511136ea3c5a (10 B)",
		},
		// buffering done?
		{
			&jsonmessage.JSONMessage{
				ID:       "511136ea3c5a",
				Status:   "Pushing",
				Progress: &jsonmessage.JSONProgress{Current: 10, Start: 0, Total: 100},
			},
			"\rPushing: 511136ea3c5a (10%)           ",
		},
		{
			&jsonmessage.JSONMessage{
				ID:       "511136ea3c5a",
				Status:   "Image successfully pushed",
				Progress: &jsonmessage.JSONProgress{Current: 0, Start: 0, Total: 0},
			},
			"\rImage successfully pushed: 511136ea3c5a\n",
		},
		{
			&jsonmessage.JSONMessage{
				ID:       "c7b7c6419568",
				Status:   "Pushing",
				Progress: &jsonmessage.JSONProgress{Current: 0, Start: 0, Total: 0},
			},
			"Pushing: c7b7c6419568",
		},
		{
			&jsonmessage.JSONMessage{
				ID:       "c7b7c6419568",
				Status:   "Buffering to disk",
				Progress: &jsonmessage.JSONProgress{Current: 524287, Start: 0, Total: 0},
			},
			"\rBuffering to disk: c7b7c6419568 (511.9 KB)",
		},
		// Buffering done?
		{
			&jsonmessage.JSONMessage{
				ID:       "c7b7c6419568",
				Status:   "Pushing",
				Progress: &jsonmessage.JSONProgress{Current: 44, Start: 0, Total: 100},
			},
			"\rPushing: c7b7c6419568 (44%)               ",
		},
		{
			&jsonmessage.JSONMessage{
				ID:       "c7b7c6419568",
				Status:   "Image successfully pushed",
				Progress: &jsonmessage.JSONProgress{Current: 0, Start: 0, Total: 0},
			},
			"\rImage successfully pushed: c7b7c6419568\n",
		},
		{
			&jsonmessage.JSONMessage{
				Status: "Pushing tag for rev [a636b9702b50] on {http://127.0.0.1:3000/v1/repositories/bvdberg/pass/tags/build-549305dd56000d6d0700027e};",
			},
			"Pushing tag for image: a636b9702b50\n", // TODO
		},
	}

	p := NewJSONMessageProcessor()
	for _, step := range testSteps {
		actual := p.ProcessJSONMessage(step.in)
		s.Equal(actual, step.expected)
	}
}

func (s *StatusHandlerSuite) TestFormatDiskUnitBytes() {
	testSteps := []struct {
		in       int64
		expected string
	}{
		{1, "1 B"},
		{1023, "1023 B"},
		{1024, "1 KB"},
		{1025, "1 KB"},
		{1536, "1.5 KB"},
		{1048575, "1023.9 KB"},
		{1048576, "1 MB"},
		{1048577, "1 MB"},
		{1073741823, "1023.9 MB"},
		{1073741824, "1 GB"},
		{1073741825, "1 GB"},
		{2147483647, "1.9 GB"},
		{1099511628800, "1024 GB"},
		{1099511628801, "1024 GB"},
	}
	for _, step := range testSteps {
		actual := formatDiskUnit(step.in)
		s.Equal(actual, step.expected)
	}
}
