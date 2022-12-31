# Creating an `ae` Release

With goreleaser, the release process should be fairly straightforward as it does majority of the work. Below is a list of steps that should get the repo prepped for release and then generate a release.

1. Make sure `.env` and the `Makefile` are updated with the appropriate version number.
2. Once you are happy with the state of `main` for the release cut, create a tag with the version number you wish to release.
   1. Ex. If you are releasing 0.0.2, make a tag for `0.0.2`.
3. This will trigger the Github Action for the goreleaser, and once it finishes successfully, you will have a Github Release for the tag.
4. The release will have a changelog of just all the commits since the last tag. You can edit it as needed.
5. You can also edit the release to add additional information, a preamble, whatever you need.
