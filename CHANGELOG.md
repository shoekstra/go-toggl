# Changelog

## [0.2.0](https://github.com/shoekstra/go-toggl/compare/v0.1.0...v0.2.0) (2026-03-25)


### 🤖 CI

* Add self-hosted Renovate for dependency and Actions updates ([#15](https://github.com/shoekstra/go-toggl/issues/15)) ([9c52c57](https://github.com/shoekstra/go-toggl/commit/9c52c57986c212c3e5b7237ceb391471f845ebcd))
* Fix Renovate Action version [skip-ci] ([#17](https://github.com/shoekstra/go-toggl/issues/17)) ([3effb43](https://github.com/shoekstra/go-toggl/commit/3effb43f77b6096db8864e3fc933f7369b72204f))
* switch Renovate to Mend GitHub App ([#22](https://github.com/shoekstra/go-toggl/issues/22)) ([852fc2a](https://github.com/shoekstra/go-toggl/commit/852fc2abec215292ace165a6cd2f73fb8677a7a2))
* Tell Renovate to autodiscover this repository [skip ci] ([#18](https://github.com/shoekstra/go-toggl/issues/18)) ([15c5029](https://github.com/shoekstra/go-toggl/commit/15c50290b5812b929710138801869dd43dcfc23f))


### 🚀 Features

* add MeService with GetMe ([#28](https://github.com/shoekstra/go-toggl/issues/28)) ([032683c](https://github.com/shoekstra/go-toggl/commit/032683c09f022c65d05f1e3dbc28191a3891af62))
* add TagsService.GetTag ([#27](https://github.com/shoekstra/go-toggl/issues/27)) ([2867218](https://github.com/shoekstra/go-toggl/commit/28672188ff1415559abf956fe1854065844f351b))


### 🐛 Fixes

* improve ErrorResponse.Error() message quality ([#29](https://github.com/shoekstra/go-toggl/issues/29)) ([757848f](https://github.com/shoekstra/go-toggl/commit/757848f759f72c6b118de1808e32686b600f97ae))
* remove GetTag and correct GetRunningTimeEntry null response ([#34](https://github.com/shoekstra/go-toggl/issues/34)) ([4b9c73e](https://github.com/shoekstra/go-toggl/commit/4b9c73e4c7f3f5e10db382b3f2a21c155b36bc41))


### 🧼 Refactoring

* rename TogglClient to WorkspaceClient ([#26](https://github.com/shoekstra/go-toggl/issues/26)) ([5fa29e8](https://github.com/shoekstra/go-toggl/commit/5fa29e813b83c2c699689e4c6904ffc2d51228ec))

## 0.1.0 (2026-03-23)


### 🤖 CI

* add release-please workflow ([#13](https://github.com/shoekstra/go-toggl/issues/13)) ([fb459cb](https://github.com/shoekstra/go-toggl/commit/fb459cbed651fba2bfe007a5fdf7a9b97b8d26dd))


### 🚀 Features

* add pagination support ([#12](https://github.com/shoekstra/go-toggl/issues/12)) ([75534ad](https://github.com/shoekstra/go-toggl/commit/75534adcacd66fb725a8d764192b162230a928ce))
* **ClientsService:** add List, Get, Create, Update, Delete, Archive, and Restore operations ([#7](https://github.com/shoekstra/go-toggl/issues/7)) ([c492aaa](https://github.com/shoekstra/go-toggl/commit/c492aaae51935b7c898bb9c323d0ea6e3e1ecc7b))
* implement ReportsService ([#9](https://github.com/shoekstra/go-toggl/issues/9)) ([aada650](https://github.com/shoekstra/go-toggl/commit/aada650c7c12cfceceef480f0ef83e849830278a))
* **ProjectsService:** add List, Get, Create, Update, and Delete operations ([#5](https://github.com/shoekstra/go-toggl/issues/5)) ([f6b8744](https://github.com/shoekstra/go-toggl/commit/f6b8744b6c7637daec2328c3a12cd9d4f1121bf9))
* **TagsService:** add List, Create, Update, and Delete operations ([#8](https://github.com/shoekstra/go-toggl/issues/8)) ([c2a046f](https://github.com/shoekstra/go-toggl/commit/c2a046fbca5f9b4d34d2dcdf6b9286d0f23dafc4))
* **TimeEntriesService:** add full CRUD operations ([#3](https://github.com/shoekstra/go-toggl/issues/3)) ([3011124](https://github.com/shoekstra/go-toggl/commit/3011124b4d3403dc92eebac16e64ab13cfb51ee9))
* **WorkspacesService:** add List, Get, and Update operations ([#4](https://github.com/shoekstra/go-toggl/issues/4)) ([079bd18](https://github.com/shoekstra/go-toggl/commit/079bd1802283f60e5a90a97a0fa0ccbe2b7f323f))


### 🧼 Refactoring

* **client:** remove post-construction mutators, fix WithTimeout ordering, drop version from User-Agent ([#11](https://github.com/shoekstra/go-toggl/issues/11)) ([19ee53c](https://github.com/shoekstra/go-toggl/commit/19ee53cd12e9b4ba18d33b3cb9a1441bfc9d8673))
