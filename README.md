[![Build Status](https://travis-ci.org/wtks/ir_decoder.svg?branch=master)](https://travis-ci.org/wtks/ir_decoder)

# IR Decoder
IR Decoderは家製協(AEHA)フォーマット用赤外線デコーダーです。

標準入力からmode2形式のpulse, space列を入力として受け取り解析し、フレーム(カスタマーコード以下)のバイト列を出力します。

フォーマットにおける単位時間Tは事前に調べる必要があります。

## 使用例
`mode2 -d /dev/lirc0 | ./ir_decoder -T 445 -h -c -l 20 -b`

## オプション
+ `-T [n]`: 単位時間(マイクロ秒)
+ `-h`: 16進数出力(`-h=false`で無効)。デフォルトでオンです。
+ `-b`: 2進数出力
+ `-c`: 前のフレームとの差分を色付けで表示します。前のフレームとフレームの長さが一致している必要があります。下記の`-skip`,`-l`と組み合わせるとおすすめです。
+ `-skip [hexString]`: 指定したhexStringのバイト列のフレームを無視する。
+ `-l`: 有効なフレームのバイト列の長さ

## ダウンロード
[Github Releases](https://github.com/wtks/ir_decoder/releases)

## ライセンス
MIT License
