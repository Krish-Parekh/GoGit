# Build Your Own Git

This is my implementation in Go for the ["Build Your Own Git" Challenge](https://codecrafters.io/challenges/git).


**Note**: Head over to [codecrafters.io](https://codecrafters.io) to try the challenge yourself.

## Stages complete

Final implementation passes all stages of [git-tester v43](https://github.com/codecrafters-io/git-tester/tree/v43):

- [x] Repository Setup
- [x] Initialize the .git directory
- [x] Read a blob object
- [x] Create a blob object
- [x] Read a tree object
- [] Write a tree object
- [] Create a commit
- [] Clone a repository

## Implemented subcommands

Only the "plumbing", low-level git commands for now (no `add`, `commit`, `status`, etc.). Just enough to complete the stages above and pass all tests.

- `init` - Does the bare minimum. Only works on current directory.
- `cat-file` - Can print size, type and content
- `hash-object` - Can calculate hash and write object to `.git/objects`
- `ls-tree` - Can list a single tree object (no recursion)
- `write-tree` - Write entire working tree recursively (no index/staging area yet)
- `commit-tree` - Write a commit object
- `clone` - Only working with remote, Smart HTTP (e.g. GitHub), repositories. Doesn't create an index yet, i.e. does just enough to pass the last stage above. Running `git checkout master` can create the index properly, though.

## To do

Continue implementing support for more subcommands as described in the [Git challenge](https://codingchallenges.fyi/challenges/challenge-git/) from [Coding Challenges](https://codingchallenges.fyi/).

**Shout out:**

- [Codecrafters](https://codecrafters.io) for the original challenge and inspiration.
- [Coding Challenges](https://codingchallenges.fyi/challenges/challenge-git/) for further ideas.

Feel free to fork, contribute, or try the challenge yourself!