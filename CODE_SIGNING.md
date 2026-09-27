# MagzClicker code signing policy

MagzClicker release binaries are built from the source code in the public GitHub repository and are intended to be signed through SignPath.

**Free code signing provided by SignPath.io, certificate by SignPath Foundation.**

## Roles

- Committer and reviewer: [Magzan16](https://github.com/Magzan16)
- Release approver: [Magzan16](https://github.com/Magzan16)

Changes submitted by other contributors must be reviewed before they are included in a release.

## Release signing process

1. The release binary is built from the public repository by GitHub Actions on a GitHub-hosted Windows runner.
2. The unsigned build is uploaded as a GitHub Actions artifact.
3. The artifact is submitted to SignPath with origin verification enabled.
4. A release approver manually approves the signing request.
5. The signed `MagzClicker.exe` is verified with Windows Authenticode.
6. A SHA-256 checksum is generated from the signed executable.
7. The signed executable and checksum are published as the GitHub release assets.

## Privacy

This program will not transfer any information to other networked systems unless specifically requested by the user or the person installing or operating it.

MagzClicker contains no telemetry, analytics, advertising, automatic update service, background service, or automatic network communication.
