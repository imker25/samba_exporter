# GitHub Actions and Release Process

This page give information on the GitHub Actions and Release Process used by the project.

## CI/CD Pipeline

For continuous integration and deployment this project uses [GitHub Actions](https://github.com/imker25/samba_exporter/actions). The main pipeline is defined in `.github/workflows/ci-jobs.yml`. This pipeline will will start on every push to github and then run the steps shown below:

```mermaid
%%{init: {'theme':'forest'}}%%
graph TD;
    push(Developer push to GitHub)
    build[Build and unit tests]
    docs[Test doc pages creation]
    insTest[Installation tests]
    intTest[Integration tests]
    checkBranch{Check branch}
    mainB((main))
    releaseB((release/*))
    otherB((other branch))
    preRel[Create GitHub -pre release]
    releaseP[Create GitHub release]
    done(Pipeline end)

    push-->build;
    push-->docs;
    build-->insTest;
    build-->intTest;
    docs-->checkBranch;
    insTest-->checkBranch;
    intTest-->checkBranch;
    checkBranch-->mainB
    checkBranch-->releaseB
    checkBranch-->otherB
    releaseB-->releaseP
    mainB-->preRel
    otherB-->done
    releaseP-->done
    preRel-->done
```

## Release Pipeline

After a GitHub release (also -pre) is done from the the CI/CD pipeline the `.github/workflows/release-jobs.yml` will be triggered. This job does the following workflow:

```mermaid
%%{init: {'theme':'forest'}}%%
graph TD;
    release(GitHub release created)
    buildFocal[Build focal *.deb package]
    buildImpish[Build impish *.deb package]
    buildBuster[Build buster *.deb package]
    buildBullseye[Build bullseye *.deb package]
    docs[Documentation pages creation]
    releaseFocalLP[Push focal *.deb to Launchpad]
    releaseImpishLP[Push impish *.deb to Launchpad]
    releaseFocalGR[Add focal *.deb to GitHub release]
    releaseImpishGR[Add impish *.deb to GitHub release]    
    releaseBullseyeGR[Add bullseye *.deb to GitHub release]
    releaseBusterGR[Add buster *.deb to GitHub release] 
    pagesRelease[Documentation release on Github pages]
    done(Pipeline end)
    checkRelease1{Check release}
    preRelease1(( -pre release))
    fullRelease1((release))
    checkRelease2{Check release}
    preRelease2(( -pre release))
    fullRelease2((release))
    checkRelease3{Check release}
    preRelease3(( -pre release))
    fullRelease3((release))

    release-->buildFocal
    buildFocal-->checkRelease1
    checkRelease1-->preRelease1
    checkRelease1-->fullRelease1
    fullRelease1-->releaseFocalLP
    releaseFocalLP-->buildImpish
    preRelease1-->buildImpish

    buildImpish-->checkRelease2
    checkRelease2-->preRelease2
    checkRelease2-->fullRelease2
    fullRelease2-->releaseImpishLP
    releaseImpishLP-->buildBullseye
    preRelease2-->buildBullseye

    buildBullseye-->buildBuster
    buildBuster-->docs
    docs-->releaseFocalGR
    releaseFocalGR-->releaseImpishGR
    releaseImpishGR-->releaseBullseyeGR
    releaseBullseyeGR-->releaseBusterGR
    releaseBusterGR-->checkRelease3

    checkRelease3-->preRelease3
    checkRelease3-->fullRelease3
    fullRelease3-->pagesRelease
    pagesRelease-->done
    preRelease3-->done
```

Whenever a *.deb package is uploaded to the [samba-exporter Launchpad PPA](https://launchpad.net/~imker/+archive/ubuntu/samba-exporter-ppa) launchpad will start a own release process. When this process is finished (usually takes about an hour), users can download and install the new package version from the PPA.

## Creation of release branches

The release process of this project is fully automated. To create a new release (not -pre) of the software use the script `build/PrepareRelease.sh`. Before running the script ensure you are on `main` branch and got the latest changes from GitHub origin. This script will:

- Create a **release** branch from the current state at the main branch
- Update the `VersionMaster.txt` with a new increment version on **main** branch
- Update the `changelog` with a stub entry for the new version on **main** branch
- Commit the changes on the main branch
- Push all changes on **main** and the **new release** branch to GitHub

Once this changes are pushed to github the CI/CD pipeline will start to run for both, `main` and the new `release` branch and create a new *-pre Release* from main as well as a new *full Release* from the new release branch. 
