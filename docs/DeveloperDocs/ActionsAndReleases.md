# GitHub Actions and Release Process

This page give information on the GitHub Actions and Release Process used by the project.

## CI/CD Pipeline

For continuous integration and deployment this project uses [GitHub Actions](https://github.com/imker25/samba_exporter/actions). The main pipeline is defined in `.github/workflows/ci-jobs.yml`. This pipeline will do:

- On push to any branch on github
  - Build the project and the man pages
  - Run unit tests defined in `*_test.go`
  - Run integration tests from `test/integrationTest/scripts/RunIntegrationTests.sh`
  - Run installation tests from `test/installationTest/RunInstallationTest.sh`
  - Build a debian binary package (`*.deb`)
- On push to main and release/* branch additionally
  - Upload the binary package (`*.deb`) as [GitHub Release](https://github.com/imker25/samba_exporter/releases)
  - In case it's the main branch the release will be a pre release
  - On release/* branches it will be a full public release

After a full public release is done from the the CI/CD run on the release/* branch `.github/workflows/release-jobs.yml` will be triggered. This job runs `build/PublishLaunchpadInDocker.sh` to:

- Create a binary and source package for Ubuntu 20.04 out of the just created release
  - Modifies the sources so a native debian build works
  - Create the binary package
  - Create the source package
  - Pushes the sources modified to a ubuntu-20.04 branch on [launchpad git repository](https://code.launchpad.net/~imker/samba-exporter/+git/samba-exporter)
  - Uploads the source package to [samba-exporter launchpad ppa](https://launchpad.net/~imker/+archive/ubuntu/samba-exporter-ppa)
    - Launchpad will trigger a own release workflow now and release the binary package on the ppa as well
- Create a binary and source package for Ubuntu 21.10 out of the just created release
  - Modifies the sources so a native debian build works
  - Create the binary package
  - Create the source package
  - Pushes the sources modified to a ubuntu-21.10 branch on [launchpad git repository](https://code.launchpad.net/~imker/samba-exporter/+git/samba-exporter)
  - Uploads the source package to [samba-exporter launchpad ppa](https://launchpad.net/~imker/+archive/ubuntu/samba-exporter-ppa)
    - Launchpad will trigger a own release workflow now and release the binary package on the ppa as well
- Create a binary and source package for Debian 10 out of the just created release
  - Modifies the sources so a native debian build works
  - Create the binary package
  - Pushes the sources modified to a debian-10 branch on [launchpad git repository](https://code.launchpad.net/~imker/samba-exporter/+git/samba-exporter)
- Create a binary and source package for Debian 11 out of the just created release
  - Modifies the sources so a native debian build works
  - Create the binary package
  - Pushes the sources modified to a debian-11 branch on [launchpad git repository](https://code.launchpad.net/~imker/samba-exporter/+git/samba-exporter)  

All created binary debian packages will be added as asset to the just created release, so users can download them.

## Release process

The release process of this project is fully automated. To create a new release of the software use the script `build/PrepareRelease.sh`. Before running the script ensure you got the latest changes from github origin. This script then will:

- Create a release branch from the current state at the main branch
- Update the `VersionMaster.txt` with a new increment version on main branch
- Update the `changelog` with a stub entry for the new version on main branch
- Commit the changes on the main branch
- Push all changes on main and the new release branch to github

Once this changes are pushed to github the CI/CD pipeline will start to run for both, main and the new release/* branch.
