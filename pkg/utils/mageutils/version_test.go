// Copyright 2021 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mageutils

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-github/v40/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/assert"
)

func newGithubMock(repositoryRelease *github.RepositoryRelease, repositoryTags []*github.RepositoryTag) *http.Client {
	return mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposReleasesLatestByOwnerByRepo,
			repositoryRelease,
		),
		mock.WithRequestMatch(
			mock.GetReposTagsByOwnerByRepo,
			repositoryTags,
		),
	)
}

type testCaseVersion struct {
	releaseType                        string
	expectedReleaseType                string
	repositoryRelease                  *github.RepositoryRelease
	repositoryTags                     []*github.RepositoryTag
	expectedNextReleaseVersion         string
	expectedNextReleaseVersionStripped string
	expectedNextReleaseBranchName      string
	expectedNextBetaVersion            string
	expectedNextRCVersion              string
}

func toStringPointer(value string) *string {
	return &value
}

func newMockRelease(tagName string) *github.RepositoryRelease {
	return &github.RepositoryRelease{Name: toStringPointer(tagName)}
}

func newMockTags(tags []string) []*github.RepositoryTag {
	var repositoryTags []*github.RepositoryTag

	for _, tag := range tags {
		repositoryTags = append(repositoryTags, &github.RepositoryTag{Name: toStringPointer(tag)})
	}

	return repositoryTags
}

func TestVersion(t *testing.T) {
	t.Run("should updated versions with according the expected", func(t *testing.T) {
		testCases := []testCaseVersion{
			{
				releaseType:                        "patch",
				repositoryRelease:                  newMockRelease("v1.0.0"),
				repositoryTags:                     newMockTags([]string{"v1.0.0", "v0.1.0", "v0.0.1"}),
				expectedNextReleaseVersion:         "v1.0.1",
				expectedNextReleaseVersionStripped: "1.0.1",
				expectedNextReleaseBranchName:      "release/v1.0",
				expectedNextBetaVersion:            "v1.0.1-beta.1",
				expectedNextRCVersion:              "v1.0.1-rc.1",
			},
			{
				releaseType:                        "minor",
				repositoryRelease:                  newMockRelease("v1.0.0"),
				repositoryTags:                     newMockTags([]string{"v1.0.0", "v0.1.0", "v0.0.1"}),
				expectedNextReleaseVersion:         "v1.1.0",
				expectedNextReleaseVersionStripped: "1.1.0",
				expectedNextReleaseBranchName:      "release/v1.1",
				expectedNextBetaVersion:            "v1.1.0-beta.1",
				expectedNextRCVersion:              "v1.1.0-rc.1",
			},
			{
				releaseType:                        "major",
				repositoryRelease:                  newMockRelease("v1.0.0"),
				repositoryTags:                     newMockTags([]string{"v1.0.0", "v0.1.0", "v0.0.1"}),
				expectedNextReleaseVersion:         "v2.0.0",
				expectedNextReleaseVersionStripped: "2.0.0",
				expectedNextReleaseBranchName:      "release/v2.0",
				expectedNextBetaVersion:            "v2.0.0-beta.1",
				expectedNextRCVersion:              "v2.0.0-rc.1",
			},
			{
				releaseType:                        "major",
				repositoryRelease:                  newMockRelease("v1.0.0"),
				repositoryTags:                     newMockTags([]string{"v1.0.0", "v0.1.0-rc.1", "v0.1.0-beta.1"}),
				expectedNextReleaseVersion:         "v2.0.0",
				expectedNextReleaseVersionStripped: "2.0.0",
				expectedNextReleaseBranchName:      "release/v2.0",
				expectedNextBetaVersion:            "v2.0.0-beta.1",
				expectedNextRCVersion:              "v2.0.0-rc.1",
			},
			{
				releaseType:                        "minor",
				repositoryRelease:                  newMockRelease("v2.0.0"),
				repositoryTags:                     newMockTags([]string{"v2.1.0-beta.1", "v2.1.0-rc.1", "v2.0.0", "v2.0.0-beta.1", "v2.0.0-rc.1"}),
				expectedNextReleaseVersion:         "v2.1.0",
				expectedNextReleaseVersionStripped: "2.1.0",
				expectedNextReleaseBranchName:      "release/v2.1",
				expectedNextBetaVersion:            "v2.1.0-beta.2",
				expectedNextRCVersion:              "v2.1.0-rc.2",
			},
			{
				releaseType:                        "minor",
				repositoryRelease:                  newMockRelease("v2.5.0"),
				repositoryTags:                     newMockTags([]string{"v2.6.0-beta.1", "v2.5.0-rc.1", "v2.5.0", "v2.4.0"}),
				expectedNextReleaseVersion:         "v2.6.0",
				expectedNextReleaseVersionStripped: "2.6.0",
				expectedNextReleaseBranchName:      "release/v2.6",
				expectedNextBetaVersion:            "v2.6.0-beta.2",
				expectedNextRCVersion:              "v2.6.0-rc.1",
			},
			{
				releaseType:                        "minor",
				repositoryRelease:                  newMockRelease("v3.3.3"),
				repositoryTags:                     newMockTags([]string{"v3.3.0-rc.1", "v3.3.0-beta.1", "v3.2.2", "v3.2.1", "v3.2.0", "v3.2.0-rc.2", "v3.2.0-rc.1", "v3.2.0-beta.3", "v3.2.0-beta.2", "v3.2.0-beta.1", "v3.1.1"}),
				expectedNextReleaseVersion:         "v3.4.0",
				expectedNextReleaseVersionStripped: "3.4.0",
				expectedNextReleaseBranchName:      "release/v3.4",
				expectedNextBetaVersion:            "v3.4.0-beta.1",
				expectedNextRCVersion:              "v3.4.0-rc.1",
			},
			{
				releaseType:                        "patch",
				repositoryRelease:                  newMockRelease("v3.2.2"),
				repositoryTags:                     newMockTags([]string{"v3.2.2", "v3.2.1", "v3.2.0", "v3.2.0-rc.2", "v3.2.0-rc.1", "v3.2.0-beta.3", "v3.2.0-beta.2", "v3.2.0-beta.1", "v3.1.1"}),
				expectedNextReleaseVersion:         "v3.2.3",
				expectedNextReleaseVersionStripped: "3.2.3",
				expectedNextReleaseBranchName:      "release/v3.2",
				expectedNextBetaVersion:            "v3.2.3-beta.1",
				expectedNextRCVersion:              "v3.2.3-rc.1",
			},
			{
				releaseType:                        "major",
				repositoryRelease:                  newMockRelease("v5.5.5"),
				repositoryTags:                     newMockTags([]string{"v5.5.0", "v5.5.0-rc.2", "v5.5.0-rc.1", "v5.5.0-beta.1"}),
				expectedNextReleaseVersion:         "v6.0.0",
				expectedNextReleaseVersionStripped: "6.0.0",
				expectedNextReleaseBranchName:      "release/v6.0",
				expectedNextBetaVersion:            "v6.0.0-beta.1",
				expectedNextRCVersion:              "v6.0.0-rc.1",
			},
			{
				releaseType:                        "minor",
				repositoryRelease:                  newMockRelease("v7.8.1"),
				repositoryTags:                     newMockTags([]string{"v7.9.0-rc.2", "v7.9.0-rc.1", "v7.9.0-beta.3", "v7.9.0-beta.2", "v7.9.0-beta.1", "v7.8.1", "v7.8.0", "v7.8.0-rc.1", "v7.8.0-beta.1"}),
				expectedNextReleaseVersion:         "v7.9.0",
				expectedNextReleaseVersionStripped: "7.9.0",
				expectedNextReleaseBranchName:      "release/v7.9",
				expectedNextBetaVersion:            "v7.9.0-beta.4",
				expectedNextRCVersion:              "v7.9.0-rc.3",
			},
			{
				releaseType:                        "major",
				repositoryRelease:                  newMockRelease("v7.8.1"),
				repositoryTags:                     newMockTags([]string{"v8.0.0-rc.2", "v8.0.0-rc.1", "v8.0.0-beta.3", "v8.0.0-beta.2", "v8.0.0-beta.1", "v7.8.1", "v7.8.0", "v7.8.0-rc.1", "v7.8.0-beta.1"}),
				expectedNextReleaseVersion:         "v8.0.0",
				expectedNextReleaseVersionStripped: "8.0.0",
				expectedNextReleaseBranchName:      "release/v8.0",
				expectedNextBetaVersion:            "v8.0.0-beta.4",
				expectedNextRCVersion:              "v8.0.0-rc.3",
			},
			{
				releaseType:                        "patch",
				repositoryRelease:                  newMockRelease("v7.8.1"),
				repositoryTags:                     newMockTags([]string{"v7.8.2-rc.2", "vv7.8.2-rc.1", "v7.8.2-beta.3", "v7.8.2-beta.2", "v7.8.2-beta.1", "v7.8.1", "v7.8.0", "v7.8.0-rc.1", "v7.8.0-beta.1"}),
				expectedNextReleaseVersion:         "v7.8.2",
				expectedNextReleaseVersionStripped: "7.8.2",
				expectedNextReleaseBranchName:      "release/v7.8",
				expectedNextBetaVersion:            "v7.8.2-beta.4",
				expectedNextRCVersion:              "v7.8.2-rc.3",
			},
			{
				releaseType:                        "minor",
				repositoryRelease:                  newMockRelease("v2.3.9"),
				repositoryTags:                     newMockTags([]string{"v2.4.0-beta.1", "v2.3.9"}),
				expectedNextReleaseVersion:         "v2.4.0",
				expectedNextReleaseVersionStripped: "2.4.0",
				expectedNextReleaseBranchName:      "release/v2.4",
				expectedNextBetaVersion:            "v2.4.0-beta.2",
				expectedNextRCVersion:              "v2.4.0-rc.1",
			},
		}

		for index, testCase := range testCases {
			version := &upVersions{
				githubClient: github.NewClient(newGithubMock(testCase.repositoryRelease, testCase.repositoryTags)),
				ctx:          context.Background(),
				releaseType:  testCase.releaseType,
				githubOrg:    "test",
				githubRepo:   "test",
			}

			fmt.Println("test case:", index+1)
			assert.NoError(t, version.setNextReleaseVersion())
			assert.NoError(t, version.setNextBetaAndRCVersion())
			assert.Equal(t, testCase.expectedNextReleaseVersion, version.nextReleaseVersion)
			assert.Equal(t, testCase.expectedNextReleaseBranchName, version.nextReleaseBranchName)
			assert.Equal(t, testCase.expectedNextReleaseVersionStripped, version.nextReleaseVersionStripped)
			assert.Equal(t, testCase.expectedNextBetaVersion, version.nextBetaVersion)
			assert.Equal(t, testCase.expectedNextRCVersion, version.nextRCVersion)
			fmt.Println()
		}
	})

	t.Run("should success parse abbreviated release type", func(t *testing.T) {
		testCases := []testCaseVersion{
			{
				releaseType:         "p",
				expectedReleaseType: "patch",
			},
			{
				releaseType:         "m",
				expectedReleaseType: "minor",
			},
			{
				releaseType:         "M",
				expectedReleaseType: "major",
			},
		}

		for _, testCase := range testCases {
			version := &upVersions{
				releaseType: testCase.releaseType,
			}

			version.parseAbbreviatedReleaseTypeName()

			assert.Equal(t, testCase.expectedReleaseType, version.releaseType)
		}
	})

	t.Run("should return error when invalid release type", func(t *testing.T) {
		testCases := []testCaseVersion{
			{
				releaseType: "something",
			},
			{
				releaseType: "N",
			},
			{
				releaseType: "majro",
			},
			{
				releaseType: "PATCH",
			},
			{
				releaseType: "MINOR",
			},
			{
				releaseType: "MAJOR",
			},
			{
				releaseType: "",
			},
		}

		for _, testCase := range testCases {
			err := UpVersions(testCase.releaseType)

			assert.Error(t, err)
			assert.Equal(t, fmt.Errorf(invalidReleaseType, testCase.releaseType), err)
		}
	})

	t.Run("should return error when missing repository info", func(t *testing.T) {
		err := UpVersions("p")

		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf(missingOrgAndRepositoryName, envRepositoryOrg, envRepositoryName), err)
	})

	t.Run("should return return no error when executing from actual repository", func(t *testing.T) {
		_ = os.Setenv(envRepositoryName, "horusec")
		_ = os.Setenv(envRepositoryOrg, "Fotkurz")

		assert.NoError(t, UpVersions("p"))
	})

	t.Run("should return error when trying to get the latest release from invalid repository", func(t *testing.T) {
		_ = os.Setenv(envRepositoryName, "!@#Invalid-Repo#@!")
		_ = os.Setenv(envRepositoryOrg, "!@#Invalid-Owner#@!")

		err := UpVersions("p")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Not Found")
	})

	t.Run("should return error when trying to list tags from invalid repository", func(t *testing.T) {
		version := upVersions{
			githubOrg:    "!@#Invalid-Owner#@!",
			githubRepo:   "!@#Invalid-Repo#@!",
			githubClient: github.NewClient(nil),
			ctx:          context.Background(),
		}

		err := version.setNextBetaAndRCVersion()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Not Found")
	})
}
