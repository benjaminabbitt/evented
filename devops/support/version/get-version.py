from git import Repo

if __name__ == "__main__":
    repo = Repo("../../../")
    changes = [item.a_path for item in repo.index.diff(None)]
    if len(changes) > 0:
        print("dynamic")
    else:
        sha = repo.head.commit.hexsha
        print(repo.git.rev_parse(sha, short=7))

