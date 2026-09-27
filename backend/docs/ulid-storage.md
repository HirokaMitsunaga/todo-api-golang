# ULID のデータベース保存形式

## 結論

このプロジェクトでは、ULID をデータベース上では `char(26)` の文字列として保存する。

GORM の永続化モデルでは、ID を `string` として定義する。

```go
type User struct {
	ID string `gorm:"type:char(26);primaryKey"`
}

type Todo struct {
	ID     string `gorm:"type:char(26);primaryKey"`
	UserID string `gorm:"type:char(26);not null;index"`
}
```

一方、ドメインモデルでは `ulid.ULID` を使用する。

```go
type User struct {
	id ulid.ULID
}
```

そのため、データベースとドメインの境界で次の変換を行う。

- DB へ保存するとき: `ulid.ULID.String()` で26文字の文字列に変換する
- DB から取得するとき: `ulid.Parse` で `ulid.ULID` に戻す

## 背景

`ulid.ULID` は Go では16バイトの値として表現され、標準の `Value()` 実装はバイナリ値を返す。そのため、GORM のモデルに `ulid.ULID` を直接使用して `bytea` に保存すると、PostgreSQL 上では人間が読める26文字の ULID ではなくバイナリ値として扱われる。

PostgreSQL の `bytea` は、標準の表示形式では `\x` に続く16進数で表示される。1バイトは16進数2文字で表示されるため、ULID 本体の16バイトは32個の16進数文字として表示される。

`char(26)` を採用する理由は次のとおりである。

- PostgreSQL のコンソールや管理ツールで ID を直接確認しやすい
- ログ、SQL、障害調査で ULID をそのまま読める
- `users.id` と `todos.user_id` の外部キーを同じ文字列型で統一できる
- DB の内容を手作業で確認・修正しやすい

## 前提

- ULID の正規の文字列表現は26文字とする
- DB の永続化モデルでは ID を `string` とする
- ドメインモデルでは ID を `ulid.ULID` とする
- `users.id`、`todos.id`、`todos.user_id` は `char(26)` で統一する
- `todos.user_id` は `users.id` を参照する外部キーとする
- `char(26)` は文字数を制約するものであり、文字列が有効な ULID であることまでは保証しない
- DB から取得した ID の妥当性確認は repository の境界で行い、異常な値はエラーとして扱う

この方針では、`ulid.ULID` を GORM のモデルに直接指定して `bytea` に保存する方式は採用しない。保存形式を `bytea` に変更する場合は、既存の `char(26)` カラムとの互換性や migration 方法を別途検討する。

## 参考文献

- [GORM: Customize Data Types](https://gorm.io/docs/data_types.html) — GORM が `database/sql` の `Scanner` / `Valuer` を実装したカスタム型を扱う仕組み
- [oklog/ulid: ULID の実装](https://github.com/oklog/ulid/blob/main/ulid.go) — `ULID` の16バイト表現と `Value()` の実装
- [PostgreSQL: Binary Data Types](https://www.postgresql.org/docs/current/datatype-binary.html) — `bytea` の保存形式と16進数表示
