# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## \[Unreleased\]

### Fixed

- Licenses format

## \[[v0.2](https://github.com/mbaraa/eloi/releases/tag/v0.2)\] 2023-03-06

### Added

- Ebuild's description

### Fixed

- Fetching ebuild's homepage

## \[[v0.1.1](https://github.com/mbaraa/eloi/releases/tag/v0.1.1)\] - 2023-02-18

### Added

- added simple overlays fetching, for just overlays' metadata fetching

## \[[v0.1](https://github.com/mbaraa/eloi/releases/tag/v0.1)\] - 2023-02-01

### Added 

- added a changelog :)

### Changed

- changed the scrapper to use concurrent jobs, which made fetching overlays x6 faster
- changed the structure of the overlay model to represent some data more elegantly

### Fixed

- made the scrapper return any occurring error, because you never know...
