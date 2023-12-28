#!/usr/bin/env python
import argparse
import os

from git import Repo

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Build a version string based on the git short SHA, if its relevant and valid, or a provided string if its invalid")
    parser.add_argument('--git_root', dest='root', default=".")
    parser.add_argument('--changed_text', dest='changed', default="latest")
    args = parser.parse_args()

    checking_dir = os.path.abspath(args.root)

    root = None
    while not root:
        if os.path.exists(os.path.join(checking_dir, ".git")):
            root = checking_dir
        else:
            old = checking_dir
            checking_dir =  os.path.abspath(os.path.dirname(checking_dir))
            if(old == checking_dir):
                raise BaseException("No repo found.")

    repo = Repo(root)
    changes = [item.a_path for item in repo.index.diff(None)]
    if len(changes) > 0:
        print(args.changed)
    else:
        sha = repo.head.commit.hexsha
        print(repo.git.rev_parse(sha, short=7))
