# Changelog

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
