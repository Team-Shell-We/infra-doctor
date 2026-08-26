package utils

const defaultConfig = `version: 1

project:
  type: spring-boot

analysis:
  exclude:
    - build
    - .gradle
    - .idea
    - .git
    - .infra-doctor

output:
  directory: .infra-doctor/generated
  overwrite: false
`

const defaultGitignore = `cache/
tmp/
`
