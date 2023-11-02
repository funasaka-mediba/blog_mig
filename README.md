# ブログ移管用画像URL変換スクリプト

# 概要
- mediba+のブログをnote proに移管することが目的。
- mediba+のブログはWordPressで動作している。
- 移管の手段はmediba+のブログ全記事をWordPressのエクスポート機能でxmlファイルで出力し、それをnote proでインポートするというもの。
- ただし、note proのインポート時に、記事内の画像URLに日本語などのマルチバイト文字が含まれていると移管されないため、移行可能な英数文字列に変換してからインポートを行う必要がある。
    - 例: http://koho.mediba.jp/wp-content/uploads/2022/02/medibaブログ_01-1024x536.png


- 今回取る手段としては、WordPress内の画像ファイルが保管されているS3で新たなS3バケットを作成し、既存の画像ファイル名変換してから、全てそこにコピーして新たなファイルURLを作成する。
- xmlファイル内の既存のファイルURLを、作成した新たなファイルURLに置換することでnote proでのインポートを可能にする。


- このリポジトリの目的は、旧ファイルURLと新ファイルURLとの変換表(csv)をもとに、エクスポートファイル内の画像ファイルURLを置換するというものである。

# 使用法

- Goを事前にインストールしておくこと。
    - 参考: https://go.dev/doc/install

- macの場合
```
$ brew install go
$ go version 
go version go1.17.6 darwin/amd64
# バージョンは1.21以上であれば良い
```

```
# リポジトリをクローン
git clone https://github.com/funasaka-mediba/blog_mig.git

# 作業ディレクトリをこのリポジトリに移動
cd ./blog_mig

# 変換表をこのリポジトリ内に移動
mv ~/Downloads/replace.csv ./

# xmlエクスポートファイルをこのリポジトリ内に移動
mv ~/Downloads/input.xml ./

# 実行
# output.xmlの部分は利用したいファイル名に適宜変換する
# report.csvは変換前の旧URLがXMLファイル内に存在していたかどうかを確認した結果を出力したものである
go run main.go replace.csv input.xml output.xml report.csv
```
