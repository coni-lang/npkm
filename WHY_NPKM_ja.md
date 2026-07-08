# 誰も信用しないスクリプトを書くのは、もうやめろ。

ちゃんと動く自動化ツールがある。

---

## 誰も直さない問題

あなたがインフラエンジニアになったのは、信頼できるインフラを作るためだ。

なのに月曜日はサーバーに1台ずつSSHして、6ヶ月前に書いたbashスクリプトを実行している。しかもそれが今でも正しく動くかどうか、もう自信がない。火曜日には、3台のサーバーが残りの4台と違う状態になっていることに気づく。いつそうなったのか、誰も分からない。水曜日には「誰が何をいつ実行したか」を調査するチケットを書く。

これはインフラじゃない。これは考古学だ。

---

## NPKMを紹介する。

**バイナリ1つ。プレイブックファイル1つ。Pythonゼロ。**

```bash
npkm -i inventory.yml playbook.yml
```

pip installなし。Galaxyアカウントなし。Ansible Towerのサブスクリプションなし。深夜2時に「virtualenvで試してみた？」というデバッグセッションなし。

Javaプロジェクトを正確に、毎回、ミリ秒単位でビルドするネイティブバイナリだけがある。全マシンで、毎回、正しく、冪等に自動化を実行するネイティブバイナリだけがある。

---

## 数字は嘘をつかない

| | Bashスクリプト | Ansible | **NPKM** |
|---|---|---|---|
| デフォルトで冪等 | ❌ 自分で実装 | ✅ あり | **✅ あり** |
| インストール | 既にある | pip + Galaxyアカウント + Python | **バイナリを1つダウンロード** |
| 適用前のドライラン | ❌ | `--check`（部分的） | **`--check` — 完全シミュレーション** |
| 実行レポート | ❌ | AWX/Tower — 有料 | **ビルトインHTML + JSON** |
| Windowsサポート | ⚠️ Batch/PSの混沌 | WinRMの苦痛 | **ネイティブPowerShell + winget** |
| エアギャップ環境 | ✅ | 困難 | **完全対応** |
| 開発用ウォッチモード | ❌ | ❌ | **`npkm watch` ビルトイン** |
| 静的解析 / リンター | ❌ | 別途インストール | **`npkm lint` ビルトイン** |
| プレイブックドキュメント | ❌ | ❌ | **`npkm --doc` — Mermaidダイアグラム** |
| 実行履歴と差分 | ❌ | ❌ | **`npkm run history diff`** |
| 学習曲線 | bashは知っている | 数日〜数週間 | **30分** |

---

## 本物の自動化とはこういうものだ

### 今のbashスクリプトはこう言っている：

```bash
#!/bin/bash
# TODO: 冪等にする
# TODO: server3で失敗する理由を調査
# TODO: 誰かが行を追加した、まだ正しいか確認
ssh user@server1 "apt-get install -y nginx"
ssh user@server2 "apt-get install -y nginx"
# server3はなぜか違う、聞かないで
ssh user@server3 "yum install -y nginx"
cp index.html user@server1:/var/www/html/
# 前回server2を忘れた
```

### NPKMはこう言う：

```yaml
- name: ウェブサーバーのセットアップ
  hosts: all
  tasks:
    - package:
        name: nginx
        state: present
    - copy:
        dest: /var/www/html/index.html
        src: files/index.html
    - service:
        name: nginx
        state: started
        enabled: true
```

**全サーバー。毎回。完全に同じ。**

---

## 本当に重要な機能

### ✅ 冪等性がビルトイン

すべてのタスクは結果を報告する：`ok`（既に完了）、`changed`（今実行した）、`skipped`（条件不一致）。同じプレイブックを10回実行しても、変更が必要なものだけを変更する。

```
TASK [ nginxをインストール ]  ok
TASK [ index.htmlをコピー ]  changed
TASK [ nginxを起動 ]  ok
```

### ✅ グループとロール — 全てを再利用

インフラをグループで定義する。タスクをロールとして一度書く。どこでも組み合わせる。

```yaml
- name: ウェブ層をプロビジョニング
  hosts: web_servers   # ← 名前付きグループを対象
  tasks:
    - include_tasks: roles/base   # ← 再利用可能なロール
    - include_tasks: roles/app
```

### ✅ group_vars — グループに従う変数

`group_vars/web_servers.edn` にファイルを置くだけで、そのグループの全ホストが自動的にそれらの変数を受け取る。コピペなし。プレイブックごとのホスト別上書きなし。

### ✅ 全てをドライラン

本番を触る前に、シミュレートする：

```bash
npkm --check -i inventory.yml deploy.yml
```

全タスクが「何をするか」を表示する。何も変わらない。自信を持ってリリースする。

### ✅ Windows？完全対応。

ネイティブPowerShell実行。`winget` と `chocolatey` パッケージ管理。ネットワーク共有からのオフラインzip展開。NPKMはLinuxをプロビジョニングするのと同じ方法でWindowsマシンをプロビジョニングする — 1つのプレイブック、1つのコマンド。

### ✅ エアギャップ環境？問題なし。

インターネット不要。ネットワーク共有から直接ツールを展開。NPKMは `apt-get` が壁に当たるロックダウンされたエンタープライズ環境でも動く。

### ✅ ビルトイン実行レポート

すべての実行でタイムスタンプ付きのダークテーマHTMLレポートを生成できる — AWXなし、Towerなし、SaaSサブスクリプションなし。

```bash
npkm --report -i inventory.yml playbook.yml
# → ~/.npkm/reports/2026-07-07_14-00-00.html
```

### ✅ 開発用ウォッチモード

タスクファイルを変更すると、NPKMが自動的に再実行する。プレイブック開発で最速のフィードバックループ。

```bash
npkm watch -i inventory.yml playbook.yml
```

### ✅ インタラクティブにステップ実行

実行前に各タスクを確認する。リスクの高い初回デプロイに最適。

```bash
npkm --step -i inventory.yml deploy.yml

TASK [ アプリケーションサーバーを停止 ]
  → このタスクを実行しますか？ [y/n/q]:
```

---

## 「でも心配なのは...」

**「すでにAnsibleを使っている。」**  
NPKMは同じYAML構文を読む。プレイブックは数日ではなく数分で移行できる。そしてPythonの依存関係チェーンを一晩で捨てられる。

**「秘密情報はどうなる？」**  
ビルトインのvault暗号化 — AES-256。`npkm vault encrypt` でファイルを暗号化する。実行時に透過的に復号される。外部のシークレットマネージャーは不要。

**「CI/CDはどうなる？」**  
単一バイナリだ。パイプラインに置くだけ。macOS、Linux、Windowsで動く。インストールするランタイムなし。

**「50台のマシンのクラスターはどうなる？」**  
プレイブックに `forks: 50` を設定する。50台全てのホストが並列にプロビジョニングされる。以上。

**「IDEサポートは？」**  
リリースzipにIntelliJプラグインが同梱されている。

---

## Bashスクリプトとansibleの本当のコスト

チームが手動でインフラを管理する毎日、こんな代償を払っている：

- 手動SSHでサーバーに接続して1デプロイあたり約**10分**
- 「なぜserver4はserver1と違うのか」のデバッグに週約**1時間**
- 新しいエンジニアへのbashスクリプト博物館のオンボーディングに四半期あたり約**1日**
- 半分だけ実行されたマイグレーションと「スクリプトはもう実行した？」というSlackメッセージに費やす**無数の時間**

5人のエンジニアチームなら、年間**何週間もの時間**が失われている — 製品ではなく、自動化の管理に。

**NPKMはその時間を返す。**

---

## 今すぐ試す

```bash
# localhostに対して実行 — SSHは不要
npkm playbook.yml

# 新しいプロジェクトをスキャフォールド
npkm init my-infra/

# リリース前に検証
npkm lint my-infra/main.edn

# 本番実行
npkm -i my-infra/inventory.edn my-infra/main.edn
```

インストールウィザードなし。アカウント登録なし。「ウォームアップ」なし。

**インフラだけが、動く。**

---

> *「800行のbashスクリプトを削除して、40行のNPKMプレイブック1つに置き換えた。3ヶ月後、全ての新しいサーバーが2分以内に自分でプロビジョニングされる。チケットなし。ドリフトなし。サプライズなし。」*

---

## NPKMを入手

📦 **ダウンロード:** [github.com/coni-lang/npkm/releases](https://github.com/coni-lang/npkm/releases)  
📖 **ドキュメント:** [NPKM-EXPLAINER_ja.md](./NPKM-EXPLAINER_ja.md)  
🔌 **IntelliJプラグイン:** リリースzipに同梱

**あなたの自動化は、深夜3時に壊れるものであるべきではない。**

NPKMは、それをあなたが信頼するものにする。
