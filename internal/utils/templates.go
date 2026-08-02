// 생성할 기본 설정 파일의 내용, config.yaml에 적힐 내용
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