# procon-gardener

### 概要
提出したACしたコードを自動的に取得してローカルのディレクトリに保存するコマンドラインツールです。

### インストール方法
インストールするにはGoが必要です。
```
go get github.com/togatoga/procon-gardener
```

### サポート環境
- Linux
- macOs

Windowsは動作確認してません。

### 使い方

1. 設定ファイルの初期化
必要な設定ファイルの作成を行います。`procon-gardener init`を実行してください。
```
% procon-gardener init  
2020/03/21 17:18:36 Initialize your config...
2020/03/21 17:18:36 Initialized your config at  /home/togatoga/.procon-gardener/config.json
```

2. 設定ファイルの編集
初期化した設定ファイルは以下のとおりです。
```
{
	"atcoder": {
		"repository_path": "",
		"user_id": "",
		"user_email": ""
	}
}
```

- `repository_path` ACしたコードを保存先のディレクトリを指定してください
- `user_id` 保存したいユーザーIDを入力してください
- `user_email` `repository_path`が`Git`リポジトリの場合、`git commit`時のメールアドレスに指定されます
こちらの`user_email`を設定しないとGitHubのアクティビティには反映されません。

3. ソースコードのアーカイブ
