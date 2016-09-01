
# API Server Auto Generated from JSON Schema

## 目的
JSON Schemaによって自動生成されるAPIサーバーの、自動生成部分を作成するための
サンプルアプリケーション

## Migration

migrationは以下のように行う
schema/schema.ymlの先頭にあるDB設定を適切に書き換え、まっさらなDBを用意して指定する必要がある

```sh:

go get -u gtthub.com/metaclosure/migo/cmd/migo
migo init
migo -y schema/schema.yml -s database_state.yml run

```