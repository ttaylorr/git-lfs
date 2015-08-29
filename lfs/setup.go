package lfs

import (
	"regexp"

	"github.com/github/git-lfs/git"
)

var (
	prePushHook = &git.Hook{
		Type:     PrePushHook,
		Contents: "#!/bin/sh\ncommand -v git-lfs >/dev/null 2>&1 || { echo >&2 \"\\nThis repository is configured for Git LFS but 'git-lfs' was not found on your path. If you no longer wish to use Git LFS, remove this hook by deleting .git/hooks/pre-push.\\n\"; exit 2; }\ngit lfs pre-push \"$@\"",
		Upgradeables: []string{
			"#!/bin/sh\ngit lfs push --stdin $*",
			"#!/bin/sh\ngit lfs push --stdin \"$@\"",
			"#!/bin/sh\ngit lfs pre-push \"$@\"",
			"#!/bin/sh\ncommand -v git-lfs >/dev/null 2>&1 || { echo >&2 \"\\nThis repository has been set up with Git LFS but Git LFS is not installed.\\n\"; exit 0; }\ngit lfs pre-push \"$@\"",
			"#!/bin/sh\ncommand -v git-lfs >/dev/null 2>&1 || { echo >&2 \"\\nThis repository has been set up with Git LFS but Git LFS is not installed.\\n\"; exit 2; }\ngit lfs pre-push \"$@\"",
		},
	}

	hooks = []git.Hook{
		prePushHook,
	}

	cleanFilter = &git.Filter{
		Name:    "clean",
		Command: "git lfs clean %%f",
	}

	// TODO(ttaylorr): this doesn't really make sense as a git filter.
	// Perhaps another abstraction is needed, like Attribute? This also
	// seems like a good candidate for something to throw on the Filters
	// type, if we still want to keep that around.
	requireFilter = &git.Filter{
		Name:  "required",
		Value: "true",
	}

	filters = git.Filters{
		cleanFilter,
		requireFilter,
	}

	valueRegexp = regexp.MustCompile("\\Agit[\\-\\s]media")
)

func InstallHooks(force bool) error {
	for _, h := range hooks {
		if err := hooks.Install(force); err != nil {
			return err
		}
	}

	return nil
}

func UninstallHooks() error {
	for _, h := range hooks {
		if err := hooks.Uninstall(force); err != nil {
			return err
		}
	}

	return nil
}

// SetupFilters installs filters necessary for git-lfs to process normal git
// operations. Currently, that list includes:
//   - clean filter
//
// An error will be returned if a filter is unable to be set, or if the required
// filters were not present.
func SetupFilters(force bool) error {
	filters.Setup()
	return nil
}

func TeardownFilters() error {
	filters.Teardown()
	return nil
}
