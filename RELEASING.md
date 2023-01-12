# Creating an `ae` Release

With goreleaser, the release process should be fairly straightforward as it does majority of the work. Below is a list of steps that should get the repo prepped for release and then generate a release.

1. Once you are happy with the state of `main` for the release cut, create a tag with the version number you wish to release.
    1. Ex. If you are releasing 0.0.2, make a tag for `v0.0.2`.
2. This will trigger the Github Action for the goreleaser, and once it finishes successfully, you will have a Github Release for the tag.
3. The release will have a changelog of just all the commits since the last tag. You can edit it as needed.
4. You can also edit the release to add additional information, a preamble, whatever you need.

If you wish to run goreleaser locally for testing the flow, you can run `make goreleaser`. If you want to adjust the default parameters you can change them by setting the `GORELEASER_FLAGS` variable. `make GORELEASER_FLAGS=--snapshot goreleaser`.
