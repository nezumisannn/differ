# Differ

## 概要
DifferはGolangで記述されたCLIツールです。    
MySQLデータベースにおけるSHOW VARIABLESの値を取得して設定値の差分を確認することができます。    

## 使い方

### インストール
GitHubのリリースからバイナリをダウンロードできます。

### MySQLコンテナ起動

```
$ docker-conpose up -d
```

### 設定ファイル(differ.config.yaml)の記述
以下が記載例です。    
- user: 接続ユーザー名
- password: ユーザーパスワード
- protocol: 接続先とポート番号

```
differ:
  database01:
    user: "root"
    password: "db01"
    protocol: "tcp(127.0.0.1:3306)"
  database02:
    user: "root"
    password: "db02"
    protocol: "tcp(127.0.0.1:3307)"
```

### コマンド実行

```
$ ./differ run --config differ.config.yaml
```

### コマンドオプション

|オプション|値|デフォルト|
|---|---|---|
|--config|/path/to/differ.config.yaml|$HOME/differ.config.yaml|
|--output|table,csv|table|

### 実行例(Table)
```
$ ./differ run --config differ.config.yaml --output table
Using config file: differ.config.yaml
+---------------------+---------------------------------------------+---------------------------------------------+------------+
|      PARAMETER      |              VALUE(DATABASE01)              |              VALUE(DATABASE02)              |   STATUS   |
+---------------------+---------------------------------------------+---------------------------------------------+------------+
| general_log_file    | /var/lib/mysql/902fa48f5495.log             | /var/lib/mysql/c345747791e5.log             | Different. |
+---------------------+---------------------------------------------+---------------------------------------------+------------+
| hostname            | 902fa48f5495                                | c345747791e5                                | Different. |
+---------------------+---------------------------------------------+---------------------------------------------+------------+
| long_query_time     |                                    1.000000 |                                   10.000000 | Different. |
+---------------------+---------------------------------------------+---------------------------------------------+------------+
| relay_log           | 902fa48f5495-relay-bin                      | c345747791e5-relay-bin                      | Different. |
+---------------------+---------------------------------------------+---------------------------------------------+------------+
| relay_log_basename  | /var/lib/mysql/902fa48f5495-relay-bin       | /var/lib/mysql/c345747791e5-relay-bin       | Different. |
+---------------------+---------------------------------------------+---------------------------------------------+------------+
| relay_log_index     | /var/lib/mysql/902fa48f5495-relay-bin.index | /var/lib/mysql/c345747791e5-relay-bin.index | Different. |
+---------------------+---------------------------------------------+---------------------------------------------+------------+
| server_uuid         | a68c61c5-68a6-11ea-acc4-0242ac1e0002        | a68accfa-68a6-11ea-9f16-0242ac1e0003        | Different. |
+---------------------+---------------------------------------------+---------------------------------------------+------------+
| slow_query_log      | ON                                          | OFF                                         | Different. |
+---------------------+---------------------------------------------+---------------------------------------------+------------+
| slow_query_log_file | /var/lib/mysql/902fa48f5495-slow.log        | /var/lib/mysql/c345747791e5-slow.log        | Different. |
+---------------------+---------------------------------------------+---------------------------------------------+------------+
| timestamp           |                           1584487815.562770 |                           1584487815.573459 | Different. |
+---------------------+---------------------------------------------+---------------------------------------------+------------+
```

### 実行例(CSV)

```
$ ./differ run  --config differ.config.yaml --output csv
Using config file: differ.config.yaml
CSV File Created.

general_log_file,/var/lib/mysql/902fa48f5495.log,/var/lib/mysql/c345747791e5.log,Different.
hostname,902fa48f5495,c345747791e5,Different.
long_query_time,1.000000,10.000000,Different.
relay_log,902fa48f5495-relay-bin,c345747791e5-relay-bin,Different.
relay_log_basename,/var/lib/mysql/902fa48f5495-relay-bin,/var/lib/mysql/c345747791e5-relay-bin,Different.
relay_log_index,/var/lib/mysql/902fa48f5495-relay-bin.index,/var/lib/mysql/c345747791e5-relay-bin.index,Different.
server_uuid,a68c61c5-68a6-11ea-acc4-0242ac1e0002,a68accfa-68a6-11ea-9f16-0242ac1e0003,Different.
slow_query_log,ON,OFF,Different.
slow_query_log_file,/var/lib/mysql/902fa48f5495-slow.log,/var/lib/mysql/c345747791e5-slow.log,Different.
timestamp,1584570005.664123,1584570005.676186,Different.
```