/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

// Package commit provides commit-related policies.
package commit

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/talos-systems/conform/internal/git"
	"github.com/talos-systems/conform/internal/policy"
)

// HeaderChecks is the configuration for checks on the header of a commit.
type HeaderChecks struct {
	// Length is the maximum length of the commit subject.
	Length int `mapstructure:"length"`
	// Imperative enforces the use of imperative verbs as the first word of a
	// commit message.
	Imperative bool `mapstructure:"imperative"`
	// HeaderCase is the case that the first word of the header must have ("upper" or "lower").
	Case string `mapstructure:"case"`
	// HeaderInvalidLastCharacters is a string containing all invalid last characters for the header.
	InvalidLastCharacters string `mapstructure:"invalidLastCharacters"`
	// Jira checks if the header containers a Jira project key.
	Jira *JiraChecks `mapstructure:"jira"`
}

// JiraChecks is the configuration for checks for Jira issues.
type JiraChecks struct {
	Keys []string `mapstructure:"keys"`
}

// BodyChecks is the configuration for checks on the body of a commit.
type BodyChecks struct {
	// Required enforces that the current commit has a body.
	Required bool `mapstructure:"required"`
}

// Commit implements the policy.Policy interface and enforces commit
// messages to conform the Conventional Commit standard.
type Commit struct {
	// SpellCheck enforces correct spelling.
	SpellCheck *SpellCheck `mapstructure:"spellcheck"`
	// Conventional is the user specified settings for conventional commits.
	Conventional *Conventional `mapstructure:"conventional"`
	// Header is the user specified settings for the header of each commit.
	Header *HeaderChecks `mapstructure:"header"`
	// Header is the user specified settings for the body of each commit.
	Body *BodyChecks `mapstructure:"body"`
	// DCO enables the Developer Certificate of Origin check.
	DCO bool `mapstructure:"dco"`
	// GPG enables the GPG signature check.
	GPG bool `mapstructure:"gpg"`
	// MaximumOfOneCommit enforces that the current commit is only one commit
	// ahead of a specified ref.
	MaximumOfOneCommit bool `mapstructure:"maximumOfOneCommit"`

	msg string
}

// FirstWordRegex is theregular expression used to find the first word in a
// commit.
var FirstWordRegex = regexp.MustCompile(`^\s*([a-zA-Z0-9]+)`)

// Compliance implements the policy.Policy.Compliance function.
// nolint: gocyclo
func (c *Commit) Compliance(options *policy.Options) (*policy.Report, error) {
	var err error

	report := &policy.Report{}

	// Setup the policy for all checks.

	var g *git.Git

	if g, err = git.NewGit(); err != nil {
		return report, errors.Errorf("failed to open git repo: %v", err)
	}

	var msg string

	if options.CommitMsgFile != nil {
		var contents []byte

		if contents, err = ioutil.ReadFile(*options.CommitMsgFile); err != nil {
			return report, errors.Errorf("failed to read commit message file: %v", err)
		}

		msg = string(contents)
	} else if msg, err = g.Message(); err != nil {
		return report, errors.Errorf("failed to get commit message: %v", err)
	}

	c.msg = msg

	if c.Header != nil {
		if c.Header.Length != 0 {
			report.AddCheck(c.ValidateHeaderLength())
		}

		if c.Header.Imperative {
			report.AddCheck(c.ValidateImperative())
		}

		if c.Header.Case != "" {
			report.AddCheck(c.ValidateHeaderCase())
		}

		if c.Header.InvalidLastCharacters != "" {
			report.AddCheck(c.ValidateHeaderLastCharacter())
		}

		if c.Header.Jira != nil {
			report.AddCheck(c.ValidateJiraCheck())
		}
	}

	if c.DCO {
		report.AddCheck(c.ValidateDCO())
	}

	if c.GPG {
		report.AddCheck(c.ValidateGPGSign(g))
	}

	if c.Conventional != nil {
		report.AddCheck(c.ValidateConventionalCommit())
	}

	if c.SpellCheck != nil {
		report.AddCheck(c.ValidateSpelling())
	}

	if c.MaximumOfOneCommit {
		report.AddCheck(c.ValidateNumberOfCommits(g, options.CommitRef))
	}

	if c.Body != nil {
		if c.Body.Required {
			report.AddCheck(c.ValidateBody())
		}
	}

	return report, nil
}

func (c Commit) firstWord() (string, error) {
	var (
		groups []string
		msg    string
	)

	if c.Conventional != nil {
		groups = parseHeader(c.msg)
		if len(groups) != 7 {
			return "", errors.Errorf("Invalid conventional commit format")
		}

		msg = groups[5]
	} else {
		msg = c.msg
	}

	if msg == "" {
		return "", errors.Errorf("Invalid msg: %s", msg)
	}

	if groups = FirstWordRegex.FindStringSubmatch(msg); groups == nil {
		return "", errors.Errorf("Invalid msg: %s", msg)
	}

	return groups[0], nil
}

func (c Commit) header() string {
	return strings.Split(strings.TrimPrefix(c.msg, "\n"), "\n")[0]
}
