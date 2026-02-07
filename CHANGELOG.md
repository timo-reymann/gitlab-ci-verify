## [2.9.0](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.8.4...v2.9.0) (2026-02-07)

### Features

* Add `--no-lint-api` flag to allow skipping CI lint locally  ([#169](https://github.com/timo-reymann/gitlab-ci-verify/issues/169)) ([a528957](https://github.com/timo-reymann/gitlab-ci-verify/commit/a528957ec2b4a398d83386427e16eb3dfe0aa04c))

## [2.8.4](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.8.3...v2.8.4) (2026-02-06)

### Bug Fixes

* Dont require login token in netrc line  ([#168](https://github.com/timo-reymann/gitlab-ci-verify/issues/168)) ([fd55bed](https://github.com/timo-reymann/gitlab-ci-verify/commit/fd55bedd40af6cbbfa4cf36585590dc17452f095))

## [2.8.3](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.8.2...v2.8.3) (2026-01-30)

### Bug Fixes

* **deps:** update module github.com/open-policy-agent/opa to v1.13.1 ([7c300db](https://github.com/timo-reymann/gitlab-ci-verify/commit/7c300dba5e24ae130455a02d920dd0ecab834661))

## [2.8.2](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.8.1...v2.8.2) (2026-01-29)

### Bug Fixes

* **deps:** update module github.com/open-policy-agent/opa to v1.13.0 ([b4ad824](https://github.com/timo-reymann/gitlab-ci-verify/commit/b4ad824a50d0c37e5667aba2a25317bd94730dee))

## [2.8.1](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.8.0...v2.8.1) (2026-01-27)

### Bug Fixes

* **deps:** update module github.com/bmatcuk/doublestar/v4 to v4.10.0 ([d8becfa](https://github.com/timo-reymann/gitlab-ci-verify/commit/d8becface5816cefc8ea72f1eb3549d8da134809))

## [2.8.0](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.7.3...v2.8.0) (2026-01-16)

### Features

* Add docker image for CI lint API proxy ([2fba0b5](https://github.com/timo-reymann/gitlab-ci-verify/commit/2fba0b582c37f132ec6d7daadf0f9efe2ff89330))
* Add error logging support to console logger ([fbaeae4](https://github.com/timo-reymann/gitlab-ci-verify/commit/fbaeae493c3f8b5d6c341819e526ff7201f1468a))
* Implement GitLab CI Lint API proxy with routing, error handling, and test coverage ([d72aa58](https://github.com/timo-reymann/gitlab-ci-verify/commit/d72aa585d0106a2d8b9bd8449c45e4368af7b70a))

### Bug Fixes

* Ensure shell checker is properly closed to avoid resource leaks ([38c18fd](https://github.com/timo-reymann/gitlab-ci-verify/commit/38c18fde97307b591edb3745b66ea963543f2706))

## [2.7.3](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.7.2...v2.7.3) (2025-12-19)

### Bug Fixes

* **deps:** update module github.com/open-policy-agent/opa to v1.12.1 ([bf712b8](https://github.com/timo-reymann/gitlab-ci-verify/commit/bf712b89c6d748c9eac0123dc60edcfd06bf59ea))

## [2.7.2](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.7.1...v2.7.2) (2025-12-17)

### Bug Fixes

* **deps:** update module github.com/open-policy-agent/opa to v1.11.1 ([e36d972](https://github.com/timo-reymann/gitlab-ci-verify/commit/e36d9721aeaca465e7e4f6da29227b53143bc9d6))

## [2.7.1](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.7.0...v2.7.1) (2025-12-14)

### Bug Fixes

* **release-process:** Adjust commit message ([ba83310](https://github.com/timo-reymann/gitlab-ci-verify/commit/ba833103af24b1524c21f58d8d2dd5334dc208d6))

## [2.7.0](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.6.1...v2.7.0) (2025-12-14)

### Features

* Release v2 template with each new semantic version ([31b5cd0](https://github.com/timo-reymann/gitlab-ci-verify/commit/31b5cd090c7be19a1a686f6cb1a3c3f0c8ab7cd8))

## [2.6.1](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.6.0...v2.6.1) (2025-12-14)

### Bug Fixes

* Use relative path for code quality reports ([5f58a35](https://github.com/timo-reymann/gitlab-ci-verify/commit/5f58a35b06b1d029687590c91879ef9b47528789))

## [2.6.0](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.5.0...v2.6.0) (2025-12-14)

### Features

* Add support for gitlab code quality format ([9144df6](https://github.com/timo-reymann/gitlab-ci-verify/commit/9144df63b46e20c2576f1b4697f553e43ab31b4b))
* Add support to write report to file ([65605d1](https://github.com/timo-reymann/gitlab-ci-verify/commit/65605d1361568e674ecf56ccd6936782c815c4ee))

## [2.5.0](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.4.6...v2.5.0) (2025-12-04)

### Features

* Add severity field to VirtualFileWarning ([26116bc](https://github.com/timo-reymann/gitlab-ci-verify/commit/26116bc16e47987d4afb893c7cf5f3abfa624813))
* Add support for wildcard includes in local includes ([9d619f8](https://github.com/timo-reymann/gitlab-ci-verify/commit/9d619f8f32dd4a10a89112b307a182dc7eb95c43)), closes [#142](https://github.com/timo-reymann/gitlab-ci-verify/issues/142)
* Add warning for glob patterns that match no files ([8738529](https://github.com/timo-reymann/gitlab-ci-verify/commit/8738529e813b99ca8a273bbfa5d9900c3632e6a4))

### Bug Fixes

* **deps:** update module github.com/bmatcuk/doublestar/v4 to v4.9.1 ([d2aa7e2](https://github.com/timo-reymann/gitlab-ci-verify/commit/d2aa7e2898c99bf98a468dc2c862ec1efe26a3c1))
* Fix supported glob patterns ([a4ef1f7](https://github.com/timo-reymann/gitlab-ci-verify/commit/a4ef1f7fb8c749716bb88baf57f15d20b1e9c784))

## [2.4.6](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.4.5...v2.4.6) (2025-11-26)

### Bug Fixes

* **deps:** update module github.com/open-policy-agent/opa to v1.11.0 ([523298f](https://github.com/timo-reymann/gitlab-ci-verify/commit/523298f2cdcb6ac3d324267b7ae40284aeb4e3ec))

## [2.4.5](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.4.4...v2.4.5) (2025-11-24)

### Bug Fixes

* **deps:** update module github.com/go-git/go-git/v5 to v5.16.4 ([4593412](https://github.com/timo-reymann/gitlab-ci-verify/commit/4593412d3923be240557c92494b46e5a5bfb2e8c))

## [2.4.4](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.4.3...v2.4.4) (2025-11-03)


### Bug Fixes

* **deps:** update module github.com/open-policy-agent/opa to v1.10.0 ([f8a2ec9](https://github.com/timo-reymann/gitlab-ci-verify/commit/f8a2ec997c9b1363ca0c6d9434af91b194a7ba33))

## [2.4.3](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.4.2...v2.4.3) (2025-10-06)


### Bug Fixes

* **deps:** update module github.com/go-git/go-git/v5 to v5.16.3 ([7ecd8c3](https://github.com/timo-reymann/gitlab-ci-verify/commit/7ecd8c33694b5a17b6d2504a5a6c563fed6c6af7))

## [2.4.2](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.4.1...v2.4.2) (2025-09-02)


### Bug Fixes

* **deps:** update module github.com/spf13/pflag to v1.0.10 ([24c1dcb](https://github.com/timo-reymann/gitlab-ci-verify/commit/24c1dcb15547be0331d2914677eb7f2bad201d99))

## [2.4.1](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.4.0...v2.4.1) (2025-09-02)


### Bug Fixes

* **deps:** update module github.com/spf13/pflag to v1.0.9 ([1e58c43](https://github.com/timo-reymann/gitlab-ci-verify/commit/1e58c439f8f0813ee15f9c47487fbeb2e61b23f4))

## [2.4.0](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.3.0...v2.4.0) (2025-08-01)


### Features

* **gitlab-packages-check:** Use dedicated finding reference ([22dc428](https://github.com/timo-reymann/gitlab-ci-verify/commit/22dc4284fecf972bdb787cd24bdc0ef43e80c621))
* **pipeline-api-lint-check:** Use dedicated finding reference ([9bfe4a1](https://github.com/timo-reymann/gitlab-ci-verify/commit/9bfe4a1e145d9b254e191ee5994d183aa29cd39e))

## [2.3.0](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.2.1...v2.3.0) (2025-08-01)


### Features

* [[#105](https://github.com/timo-reymann/gitlab-ci-verify/issues/105)] Add fingerprint for check findings ([f811414](https://github.com/timo-reymann/gitlab-ci-verify/commit/f8114149131b9a6e90e4571f14ab2702fb1a26eb))
* [[#105](https://github.com/timo-reymann/gitlab-ci-verify/issues/105)] Deduplicate findings for the same line ([9b8ce7d](https://github.com/timo-reymann/gitlab-ci-verify/commit/9b8ce7d6645a99f59744fcb0ed7ae81ed11ef641))

## [2.2.1](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.2.0...v2.2.1) (2025-08-01)


### Bug Fixes

* [[#105](https://github.com/timo-reymann/gitlab-ci-verify/issues/105)] Fix CI env check ([cb03358](https://github.com/timo-reymann/gitlab-ci-verify/commit/cb0335870c0909ae797a3d31937676d6a44d7218))
* **deps:** update module github.com/open-policy-agent/opa to v1.7.1 ([1da6508](https://github.com/timo-reymann/gitlab-ci-verify/commit/1da65080600f7c5a38e3bb32276ff4960356186a))

## [2.2.0](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.1.15...v2.2.0) (2025-07-31)


### Features

* [[#104](https://github.com/timo-reymann/gitlab-ci-verify/issues/104)] Add panic when resolving location fails for shellcheck finding ([44daf92](https://github.com/timo-reymann/gitlab-ci-verify/commit/44daf92ee3816c5a22ed376d2d13425e5e8badba))


### Bug Fixes

* [[#104](https://github.com/timo-reymann/gitlab-ci-verify/issues/104)] Fix virtual file concat ([2fa95c3](https://github.com/timo-reymann/gitlab-ci-verify/commit/2fa95c39b2dcde5be4fe7d26e588c80822777080))
* Only create merged ci yaml when gitlab ci yaml is valid ([9d46deb](https://github.com/timo-reymann/gitlab-ci-verify/commit/9d46deb835572944d8935682658de74857abe62a))

## [2.1.14](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.1.13...v2.1.14) (2025-07-30)


### Bug Fixes

* Set package to v2 ([57f16bb](https://github.com/timo-reymann/gitlab-ci-verify/commit/57f16bb223f2e706523e892ae51fbc0b5337dc2c))

## [2.1.13](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.1.12...v2.1.13) (2025-07-30)


### Bug Fixes

* **deps:** update module github.com/spf13/pflag to v1.0.7 ([bd9cf26](https://github.com/timo-reymann/gitlab-ci-verify/commit/bd9cf269f6ab1120f0dcad5ef5c49bcf21a1f6ca))

## [2.1.12](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.1.11...v2.1.12) (2025-06-30)


### Bug Fixes

* **deps:** update module github.com/open-policy-agent/opa to v1.6.0 ([c67ad53](https://github.com/timo-reymann/gitlab-ci-verify/commit/c67ad536188c933ebd5f5efaba437eb2c5985010))

## [2.1.11](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.1.10...v2.1.11) (2025-06-22)


### Bug Fixes

* **deps:** update module github.com/hashicorp/go-retryablehttp to v0.7.8 ([9e86eb0](https://github.com/timo-reymann/gitlab-ci-verify/commit/9e86eb06fa15461b1d62f349a75023f4515b69ba))

## [2.1.10](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.1.9...v2.1.10) (2025-06-12)


### Bug Fixes

* **deps:** update dependency coverage to ==7.9.* ([918779b](https://github.com/timo-reymann/gitlab-ci-verify/commit/918779b7c309aa2a9d633e27560d6270d48996c3))
* **deps:** update module github.com/go-git/go-git/v5 to v5.16.2 ([e8d340e](https://github.com/timo-reymann/gitlab-ci-verify/commit/e8d340e7053b2a3951208eeab3a5f62dad4249be))

## [2.1.9](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.1.8...v2.1.9) (2025-06-04)


### Bug Fixes

* **deps:** update module github.com/go-git/go-git/v5 to v5.16.1 ([e645b33](https://github.com/timo-reymann/gitlab-ci-verify/commit/e645b331817920a37ac76fdb69826ff56acee1d1))
* **deps:** update module github.com/open-policy-agent/opa to v1.5.1 ([850726f](https://github.com/timo-reymann/gitlab-ci-verify/commit/850726f05891c474dceaa070a50dc5d007e8b478))

## [2.1.8](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.1.7...v2.1.8) (2025-05-29)


### Bug Fixes

* **deps:** update module github.com/open-policy-agent/opa to v1.5.0 ([a7644e5](https://github.com/timo-reymann/gitlab-ci-verify/commit/a7644e5422ba3875576b8f0485581fc3d34db0fb))

## [2.1.7](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.1.6...v2.1.7) (2025-05-23)


### Bug Fixes

* **deps:** update dependency pydoctor to v25 ([4aec541](https://github.com/timo-reymann/gitlab-ci-verify/commit/4aec541de5d43d3d69c5bfb9a438edd42e0553df))

## [2.1.6](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.1.5...v2.1.6) (2025-05-05)


### Bug Fixes

* **deps:** update module github.com/open-policy-agent/opa to v1.4.2 ([7862b77](https://github.com/timo-reymann/gitlab-ci-verify/commit/7862b77f98f8f1cc96fe5f8642950777c549e2a8))

## [2.1.5](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.1.4...v2.1.5) (2025-04-16)


### Bug Fixes

* **deps:** update module github.com/go-git/go-git/v5 to v5.16.0 ([ec5000b](https://github.com/timo-reymann/gitlab-ci-verify/commit/ec5000b45e499c654381337334bcb86564bea0ce))

## [2.1.4](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.1.3...v2.1.4) (2025-04-11)


### Bug Fixes

* **deps:** update module github.com/go-git/go-git/v5 to v5.15.0 ([e69fe47](https://github.com/timo-reymann/gitlab-ci-verify/commit/e69fe476b5264e17b92d183cb486bf93b7962aab))

## [2.1.3](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.1.2...v2.1.3) (2025-04-01)


### Bug Fixes

* **deps:** update dependency coverage to ==7.8.* ([f1ba667](https://github.com/timo-reymann/gitlab-ci-verify/commit/f1ba667c0ff6c3c438891f98d6dcbcc533afeea7))

## [2.1.2](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.1.1...v2.1.2) (2025-03-27)


### Bug Fixes

* **deps:** update module github.com/open-policy-agent/opa to v1.3.0 ([f2a62ea](https://github.com/timo-reymann/gitlab-ci-verify/commit/f2a62ea9c83f66f0ffa5b017893915f9165d0808))

## [2.1.1](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.1.0...v2.1.1) (2025-03-22)


### Bug Fixes

* **python-wrapper:** Downgrade setuptools ([961f83c](https://github.com/timo-reymann/gitlab-ci-verify/commit/961f83c963523fcde0edbf140d952ccfee7f4ce7))

## [2.1.0](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.0.2...v2.1.0) (2025-03-22)


### Features

* [[#71](https://github.com/timo-reymann/gitlab-ci-verify/issues/71)] Add support for ignoring check findings with YAML comments ([a8a4916](https://github.com/timo-reymann/gitlab-ci-verify/commit/a8a49165ff8b97266caf5ed9030e3251473011e1))


### Bug Fixes

* **deps:** update dependency setuptools to v77 ([450592f](https://github.com/timo-reymann/gitlab-ci-verify/commit/450592f6d75fc50ee07421b3373b8545c7374cdf))
* Fix CI lint not done in some cases ([6b56b58](https://github.com/timo-reymann/gitlab-ci-verify/commit/6b56b58f4e0bf8d99ce962c7fad0a8fdcf95588f))
* Fix fail severity check ([b117d45](https://github.com/timo-reymann/gitlab-ci-verify/commit/b117d45b27ddfd85a82c8a2f0545f3f42d91f7bc))

## [2.0.2](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.0.1...v2.0.2) (2025-03-17)


### Bug Fixes

* **deps:** update dependency coverage to ==7.7.* ([44a56c8](https://github.com/timo-reymann/gitlab-ci-verify/commit/44a56c8ac952c01c68aead68b6ee5bd9acbc75ad))

## [2.0.1](https://github.com/timo-reymann/gitlab-ci-verify/compare/v2.0.0...v2.0.1) (2025-03-15)


### Bug Fixes

* **deps:** update dependency gitlab-ci-verify-bin to v2 ([9544e41](https://github.com/timo-reymann/gitlab-ci-verify/commit/9544e413cfabc595399be4a9ecf96a8b435a4898))

## [2.0.0](https://github.com/timo-reymann/gitlab-ci-verify/compare/v1.2.11...v2.0.0) (2025-03-15)


### ⚠ BREAKING CHANGES

* Migrate check input to virtual ci yaml

### Features

* Add rego file and line resolve ([daf90e3](https://github.com/timo-reymann/gitlab-ci-verify/commit/daf90e3c89d57aff578dd27231eb2469236c439f))
* Add resolve for virtual file ([0d8076e](https://github.com/timo-reymann/gitlab-ci-verify/commit/0d8076e12c6c0a42b1443ffa5b87a29cfad256ad))
* add resolve util to check input ([ba27491](https://github.com/timo-reymann/gitlab-ci-verify/commit/ba27491cacaf8e6313b27613f9ac0d744347c32d))
* Add virtual ci yaml file ([c7d7c16](https://github.com/timo-reymann/gitlab-ci-verify/commit/c7d7c16e706bb2aac3e08f6c7cc407e00df352ea))
* **includes:** Add base code ([5efc251](https://github.com/timo-reymann/gitlab-ci-verify/commit/5efc251dda13ef441b0f98e576808f1979da0f0e))
* **includes:** Add component support ([7fe8614](https://github.com/timo-reymann/gitlab-ci-verify/commit/7fe8614dac834b5ab71af237d819df065b4ed600))
* **includes:** Add local support ([ca441f0](https://github.com/timo-reymann/gitlab-ci-verify/commit/ca441f0877b21e8df16d40575dcdddd2034a50b7))
* **includes:** Add project support ([8bba82f](https://github.com/timo-reymann/gitlab-ci-verify/commit/8bba82f5f7cad84cd3641528879ebaafb87641d4))
* **includes:** Add remote support ([c09a0b6](https://github.com/timo-reymann/gitlab-ci-verify/commit/c09a0b6284e14dd4e3ed551a3c09b7141c86ea57))
* **includes:** Add template support ([3322de3](https://github.com/timo-reymann/gitlab-ci-verify/commit/3322de31512b6c817d8d08157e1790ab115c3fcf))
* Migrate check input to virtual ci yaml ([25d48a7](https://github.com/timo-reymann/gitlab-ci-verify/commit/25d48a75f377af9ce89cf240135a7fe634db12ac))
* Use decoder with non unique keys for virtual ci file ([fd146bd](https://github.com/timo-reymann/gitlab-ci-verify/commit/fd146bd99abc249bc647804323e4c0580b864be0))


### Bug Fixes

* Adjust yamlpath line numbers ([5276daf](https://github.com/timo-reymann/gitlab-ci-verify/commit/5276dafc15a600f45a974974247d9789a7ce7e08))
* Fix node attachment for local include ([20aa395](https://github.com/timo-reymann/gitlab-ci-verify/commit/20aa39520b7c22e3274a31ac10dcdb2d0b9dab88))
* Fix virtual file appendix ([a33accb](https://github.com/timo-reymann/gitlab-ci-verify/commit/a33accb541c98db24ea8a140d888bfed786fab58))
* Resolve ci yaml to project path ([d37f81d](https://github.com/timo-reymann/gitlab-ci-verify/commit/d37f81dfe827d6a30fec50892494da432b7dba9a))
* Use combined file content for lint ([4b3ed30](https://github.com/timo-reymann/gitlab-ci-verify/commit/4b3ed30ee8950a97460ca833a53793ffb404e121))

## [1.2.11](https://github.com/timo-reymann/gitlab-ci-verify/compare/v1.2.10...v1.2.11) (2025-03-11)


### Bug Fixes

* **deps:** update dependency setuptools to v76 ([19de2c8](https://github.com/timo-reymann/gitlab-ci-verify/commit/19de2c8cb443d9f9730f327a0f91d3820b01e14d))

## [1.2.10](https://github.com/timo-reymann/gitlab-ci-verify/compare/v1.2.9...v1.2.10) (2025-02-28)


### Bug Fixes

* **deps:** update module github.com/open-policy-agent/opa to v1.2.0 ([fb6d875](https://github.com/timo-reymann/gitlab-ci-verify/commit/fb6d875493e291fdd8c95d41177f6ed33176d202))

## [1.2.9](https://github.com/timo-reymann/gitlab-ci-verify/compare/v1.2.8...v1.2.9) (2025-02-27)


### Bug Fixes

* **deps:** update module github.com/go-git/go-git/v5 to v5.14.0 ([78490dd](https://github.com/timo-reymann/gitlab-ci-verify/commit/78490ddc1f0e9b64c164e350f65de1543947628d))

## [1.2.8](https://github.com/timo-reymann/gitlab-ci-verify/compare/v1.2.7...v1.2.8) (2025-02-22)


### Bug Fixes

* Fix pypi release ([82e333a](https://github.com/timo-reymann/gitlab-ci-verify/commit/82e333a3e6d35739c0e515b9aa6d2ea0e7970a08))

## [1.2.6](https://github.com/timo-reymann/gitlab-ci-verify/compare/v1.2.5...v1.2.6) (2025-01-29)


### Bug Fixes

* **deps:** update module github.com/spf13/pflag to v1.0.6 ([4af8dc0](https://github.com/timo-reymann/gitlab-ci-verify/commit/4af8dc099b3d90376bb2c6defa9c2e7a08a78be3))

## [1.2.5](https://github.com/timo-reymann/gitlab-ci-verify/compare/v1.2.4...v1.2.5) (2025-01-27)


### Bug Fixes

* **deps:** update module github.com/open-policy-agent/opa to v1.1.0 ([0661f32](https://github.com/timo-reymann/gitlab-ci-verify/commit/0661f32cae97f425665824e33a121be289946a89))

## [1.2.4](https://github.com/timo-reymann/gitlab-ci-verify/compare/v1.2.3...v1.2.4) (2025-01-23)


### Bug Fixes

* **deps:** update module github.com/go-git/go-git/v5 to v5.13.2 ([fd91c33](https://github.com/timo-reymann/gitlab-ci-verify/commit/fd91c33b926137ee6e6de89a2a6095a1313f4f13))

## [1.2.3](https://github.com/timo-reymann/gitlab-ci-verify/compare/v1.2.2...v1.2.3) (2025-01-21)


### Bug Fixes

* **deps:** update dependency twine to ==6.1.* ([b928b9d](https://github.com/timo-reymann/gitlab-ci-verify/commit/b928b9dec6c0ff6246fb23cfd341f5277368b566))

## [1.2.2](https://github.com/timo-reymann/gitlab-ci-verify/compare/v1.2.1...v1.2.2) (2025-01-21)


### Bug Fixes

* **deps:** update module github.com/open-policy-agent/opa to v1.0.1 ([23cd4a4](https://github.com/timo-reymann/gitlab-ci-verify/commit/23cd4a44a58705d129622b4beafe9e6762236d80))

## [1.2.1](https://github.com/timo-reymann/gitlab-ci-verify/compare/v1.2.0...v1.2.1) (2025-01-20)


### Bug Fixes

* [#43](https://github.com/timo-reymann/gitlab-ci-verify/issues/43) Handle yaml multiline list item continuation ([1248f29](https://github.com/timo-reymann/gitlab-ci-verify/commit/1248f298900a361392e9cc6e8a1d0c5e8616d4cb))

## [1.2.0](https://github.com/timo-reymann/gitlab-ci-verify/compare/v1.1.1...v1.2.0) (2025-01-02)


### Features

* Add `--include-opa-bundle` flag ([537e0ab](https://github.com/timo-reymann/gitlab-ci-verify/commit/537e0abec20895a6e28daf27a96ea4034aba59d4))
* **rego-policy:** Allow loading remote bundle ([a2bb481](https://github.com/timo-reymann/gitlab-ci-verify/commit/a2bb481156c4b419909308a1b4e224acf64901ab))


### Bug Fixes

* Add missing integration tests file ([59866bd](https://github.com/timo-reymann/gitlab-ci-verify/commit/59866bd4277e7de3afcae4fe0d6a8250bae3473c))
* **deps:** update module github.com/go-git/go-git/v5 to v5.13.1 ([8f303a9](https://github.com/timo-reymann/gitlab-ci-verify/commit/8f303a9121799cdd3227f30d31514ae2bf29e477))

## [1.1.1](https://github.com/timo-reymann/gitlab-ci-verify/compare/v1.1.0...v1.1.1) (2024-12-29)


### Bug Fixes

* **deps:** update module github.com/go-git/go-git/v5 to v5.13.0 ([3f0c85f](https://github.com/timo-reymann/gitlab-ci-verify/commit/3f0c85fefcca7e8d3c24008e53ae85e4bd93a71c))

## [1.1.0](https://github.com/timo-reymann/gitlab-ci-verify/compare/v1.0.1...v1.1.0) (2024-12-28)


### Features

* Add draft for rego bundle check ([6e93844](https://github.com/timo-reymann/gitlab-ci-verify/commit/6e938440ec730e3051cec745213e338f2ea5551a))
* Add more verbose rego logging ([3be1ac8](https://github.com/timo-reymann/gitlab-ci-verify/commit/3be1ac894fb91523e6bb9432de1eb82a5c99e548))
* Add rego policy library code ([25e0fe9](https://github.com/timo-reymann/gitlab-ci-verify/commit/25e0fe9332a2ca35f551a89a0c2388fe05579de2))
* Add support for project local policies ([050bd58](https://github.com/timo-reymann/gitlab-ci-verify/commit/050bd58d5c97a57ac4f61ce2569dc9a0361cc7fe))
* Add yaml path utils to get line numbers from yamlpath query ([26665ea](https://github.com/timo-reymann/gitlab-ci-verify/commit/26665ea907bf4bf6e03f5c53f524408edb30cde8))
* Port gitlab pages check to rego ([ad0967b](https://github.com/timo-reymann/gitlab-ci-verify/commit/ad0967bb439e0e36f05c233377e6b78403211cd7))


### Bug Fixes

* Make sure merged YAML is included ([2b3536c](https://github.com/timo-reymann/gitlab-ci-verify/commit/2b3536c49b782dff177f00887e4732bc1cced448))

## [1.0.1](https://github.com/timo-reymann/gitlab-ci-verify/compare/v1.0.0...v1.0.1) (2024-12-22)


### Bug Fixes

* **deps:** update dependency gitlab-ci-verify-bin to v1 ([#36](https://github.com/timo-reymann/gitlab-ci-verify/issues/36)) ([eab00fa](https://github.com/timo-reymann/gitlab-ci-verify/commit/eab00fa2d70ffaabbe3ca346219a90c90b560b01))

## [1.0.0](https://github.com/timo-reymann/gitlab-ci-verify/compare/v0.6.0...v1.0.0) (2024-12-21)


### ⚠ BREAKING CHANGES

* --no-syntax-validate-in-ci has been removed in favor of --no-lint-api-in-ci, resulting in same behaviour

### Features

* Add more details for failed lint requests ([339b281](https://github.com/timo-reymann/gitlab-ci-verify/commit/339b281ae483430c273212f76f88b3a04f3e2587))
* Expose ci lint api result to checks ([ce73c3a](https://github.com/timo-reymann/gitlab-ci-verify/commit/ce73c3acdc907c6d748f754023a25c15343370bf))


### Bug Fixes

* Handle forbidden errors correctly ([fa8317b](https://github.com/timo-reymann/gitlab-ci-verify/commit/fa8317b7f4740f6d367ac5912b653eff23a20bae))

## [0.6.0](https://github.com/timo-reymann/gitlab-ci-verify/compare/v0.5.4...v0.6.0) (2024-12-19)


### Features

* Add retries to http client ([74d9a07](https://github.com/timo-reymann/gitlab-ci-verify/commit/74d9a0731f673a23df6b2f0730615c261f8f5493))

## [0.5.4](https://github.com/timo-reymann/gitlab-ci-verify/compare/v0.5.3...v0.5.4) (2024-12-10)


### Bug Fixes

* **python-wrapper:** Use unpack for validator config ([491a1ec](https://github.com/timo-reymann/gitlab-ci-verify/commit/491a1ec140c824c4eb7d1a766802bc6560baeccb))

## [0.5.3](https://github.com/timo-reymann/gitlab-ci-verify/compare/v0.5.2...v0.5.3) (2024-12-02)


### Bug Fixes

* **deps:** update dependency pydoctor to ==24.11.* ([94da34b](https://github.com/timo-reymann/gitlab-ci-verify/commit/94da34be6ba83cc24281db1dd45b6f518b7c4fdb))

## [0.5.2](https://github.com/timo-reymann/gitlab-ci-verify/compare/v0.5.1...v0.5.2) (2024-11-30)


### Bug Fixes

* **deps:** update dependency twine to v6 ([ecdc80d](https://github.com/timo-reymann/gitlab-ci-verify/commit/ecdc80d9f4be90b716fa75b7eba59e9af07b1d5f))

## [0.5.1](https://github.com/timo-reymann/gitlab-ci-verify/compare/v0.5.0...v0.5.1) (2024-11-09)


### Bug Fixes

* **deps:** update dependency wheel to ==0.45.* ([30efaa7](https://github.com/timo-reymann/gitlab-ci-verify/commit/30efaa7b7cd9757eee74b662ebfce4510b628709))

## [0.5.0](https://github.com/timo-reymann/gitlab-ci-verify/compare/v0.4.2...v0.5.0) (2024-10-25)


### Features

* Add --no-syntax-validate-in-ci flag ([a0acdad](https://github.com/timo-reymann/gitlab-ci-verify/commit/a0acdad0b9c757d3243aabfd2e0cc640d71f45a8))

## [0.4.2](https://github.com/timo-reymann/gitlab-ci-verify/compare/v0.4.1...v0.4.2) (2024-10-25)


### Bug Fixes

* **vault:** Fix support for versioned k/v store ([164cddf](https://github.com/timo-reymann/gitlab-ci-verify/commit/164cddf42289c42aac15d5fa3f81c15518456331))

## [0.4.1](https://github.com/timo-reymann/gitlab-ci-verify/compare/v0.4.0...v0.4.1) (2024-10-22)


### Bug Fixes

* **deps:** update module github.com/fatih/color to v1.18.0 ([d529bfa](https://github.com/timo-reymann/gitlab-ci-verify/commit/d529bfa63ba33a6ee2dbbca9ac9834d667a1432a))

## [0.4.0](https://github.com/timo-reymann/gitlab-ci-verify/compare/v0.3.0...v0.4.0) (2024-10-14)


### Features

* Add token source abstraction ([10f1c37](https://github.com/timo-reymann/gitlab-ci-verify/commit/10f1c375f13d01e79feb28ffd2f9754a805910c4))
* Add vault functionality ([d42dd74](https://github.com/timo-reymann/gitlab-ci-verify/commit/d42dd7426cd35bccb33691705da0aa67eaf43391))
* Support vault for gitlab api client ([3216a91](https://github.com/timo-reymann/gitlab-ci-verify/commit/3216a917dfb96aa9f78715a499abceaa23689777))

## [0.3.0](https://github.com/timo-reymann/gitlab-ci-verify/compare/v0.2.2...v0.3.0) (2024-09-29)


### Features

* Upload python wrapper dist to gh release ([e82352a](https://github.com/timo-reymann/gitlab-ci-verify/commit/e82352a6be599094290c03859881dddbd3de01b1))

## [0.2.2](https://github.com/timo-reymann/gitlab-ci-verify/compare/v0.2.1...v0.2.2) (2024-09-27)


### Bug Fixes

* Checkout for gh release creation ([6c1eab4](https://github.com/timo-reymann/gitlab-ci-verify/commit/6c1eab45440a389939e92ba687fedf14ea89599a))

## [0.2.1](https://github.com/timo-reymann/gitlab-ci-verify/compare/v0.2.0...v0.2.1) (2024-09-27)


### Bug Fixes

* Fix binary upload for release ([1a74bd1](https://github.com/timo-reymann/gitlab-ci-verify/commit/1a74bd10023100439302baab70195978fdf42875))

## [0.2.0](https://github.com/timo-reymann/gitlab-ci-verify/compare/v0.1.0...v0.2.0) (2024-09-27)


### Features

* Attach binaries to github release and set recommended installation order ([710f9a5](https://github.com/timo-reymann/gitlab-ci-verify/commit/710f9a54cc98434c3678ff8b9b35b29ca50ae1a7))

## [0.1.0](https://github.com/timo-reymann/gitlab-ci-verify/compare/v0.0.18...v0.1.0) (2024-09-27)


### Features

* Trigger first release ([5a60ea6](https://github.com/timo-reymann/gitlab-ci-verify/commit/5a60ea66e6dadde0efe2a696dd86142f60578d9a))
