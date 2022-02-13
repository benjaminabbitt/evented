#!/usr/bin/env python
import argparse

from git import Repo

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Build a version string based on the git short SHA, if its relevant and valid, or a provided string if its invalid")
    parser.add_argument('--git_root', dest='root', default=".")
    parser.add_argument('--changed_text', dest='changed', default="dynamic")
    args = parser.parse_args()
    repo = Repo(args.root)
    changes = [item.a_path for item in repo.index.diff(None)]
    if len(changes) > 0:
        print(args.changed)
    else:
        sha = repo.head.commit.hexsha
        print(repo.git.rev_parse(sha, short=7))
