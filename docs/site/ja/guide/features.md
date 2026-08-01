# 主な機能

`bwsf`は、.envファイルの円滑な管理を、[Bitwarden](https://bitwarden.com/)を用いて、セキュアに行うためのヘルパーコマンドです。

## .envファイルの保存

```bash
# cd /path/to/your/project_root
bwsf push
```

このコマンドは、プロジェクトルートにある.envファイルを、まとめてBitwardenに保存します。

- .env
- .env.local
- .env.staging
- .env.production

これらのファイルは、まとめてBitwardenにアップロードされます。

なお、`.env.local.example`といった、`.example`の文字列を含む設定ファイルは、Bitwardenには保存されません。

## .envファイルのプロジェクトへの適用

```bash
# cd /path/to/your/project_root
bwsf pull
```

Bitwarden側で保存されている、当該プロジェクトの.envファイルを、まとめてBitwardenに保存します。

## Bitwarden上を使ったマルチユーザーでの共有

Bitwarden側では、`dotenvs`というフォルダに保存されます。（これは予約語）です。

`bwsf`をプロジェクトルートで実行した時、そのルートフォルダの名前がプロジェクトネームとなります。

`dotenvs`フォルダを、Bitwarden上で他のユーザーと共有することにより、.envファイルを、複数のメンバーで共有することができます。

詳しくは、[Bitwardenのドキュメント](https://bitwarden.com/resources/)をご覧ください。