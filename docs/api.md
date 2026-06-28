# API

特に指定が無い場合は、

- `Content-type: application/json`
- 日時はISO 8601形式
- 全て `/api` 以下にする
    - `/api/posts` など

ユーザー認証は、NeoShowcase の機能を使う。ヘッダーからtraQのユーザー名を取得する。 https://wiki.trap.jp/services/NeoShowcase#head20 の hard 認証。ログインに失敗した場合は 401。

## `POST /posts`

投稿する。

### request

201

`Content-type: multipart/form-data`

- `image` 画像ファイルのバイナリ (png,jpeg)
- `tags` コンマ（`,`）区切りのタグ。`multipart/form-data` だと配列を送れないので、コンマ区切りで送ってサーバーでsplitする。タグ無しの場合は空文字列。タグにコンマが含まれないかはフロントエンドでの検証が必要

400 

- `image` が画像(png,jpeg)ではない
- `tags` が10個より多い
- `tags` に空文字列もしくは17文字以上のタグが含まれる。

### response

- `id` 投稿のid(uuid)

```json
{
    "id": "031e1831-49eb-4dc2-ae07-e8aa9dc1e483"
}
```

## `GET /posts` 

投稿を新しい順に取得する。

### request

クエリパラメータ

- `before` UUID このidより古い投稿を取得する。指定が無い場合は最も新しい投稿から取得する。
- `limit` number(int) この件数だけ取得する。指定が無い場合は30件取得する。

### response

投稿日時の降順（新しい順）

- `id` id
- `userName` traQのユーザー名
- `tags` タグの配列。タグが無い場合は空配列
- `imageUrl` 画像の一時URL
- `reactions` リアクション。数が0のリアクションは配列に含まれない。リアクションが一つも無い場合は空配列
    - `id` リアクションのid (number)
    - `count` リアクション数
    - `myReaction` 自身がリアクションしたかどうか
- `createdAt` 投稿の作成日時

```json
[
    {
        "id": "031e1831-49eb-4dc2-ae07-e8aa9dc1e483",
        "userName": "ikura-hamu",
        "tags": [
            "四川屋台"
        ],
        "imageUrl": "https://example.com/image",
        "reactions": [
            {
                "id": 1,
                "count": 10,
                "myReaction": true
            }
        ],
        "createdAt": "2026-06-25T13:04:54.744Z"
    }
]
```

画像の URL は、s3 の presigned URL をリクエストごとに発行する。

## `GET /posts/{id}` 

投稿を取得する。

### request

### response

- `id` id
- `userName` traQのユーザー名
- `tags` タグの配列。タグが無い場合は空配列
- `imageUrl` 画像の一時URL
- `reactions` リアクション。数が0のリアクションは配列に含まれない。リアクションが一つも無い場合は空配列
    - `id` リアクションのid (number)
    - `count` リアクション数
    - `myReaction` 自身がリアクションしたかどうか
- `createdAt` 投稿の作成日時

```json
{
    "id": "031e1831-49eb-4dc2-ae07-e8aa9dc1e483",
    "userName": "ikura-hamu",
    "tags": [
        "四川屋台"
    ],
    "imageUrl": "https://example.com/image",
    "reactions": [
        {
            "id": 1,
            "count": 10,
            "myReaction": true
        }
    ],
    "createdAt": "2026-06-25T13:04:54.744Z"
}
```
 
## `DELETE /posts/{id}`

投稿を削除する。自身の投稿のみ削除できる。
投稿に紐づいたリアクションとかタグとの結びつきとかも削除する。

### request

- `id` 投稿のID

### response

- 200
- 403 自身の投稿でない
- 404 該当する投稿が存在しない
 
## `POST /posts/{id}/reactions`

投稿にリアクションを付ける。
実装を簡単にするために、リアクションの id と 絵文字の種類の組み合わせはフロントエンドのみで扱い、バックエンドはidのみを扱う。

### request

```json
{
    "id": 1
}
```

### response

```json
[
    {
        "id": 1,
        "count": 10,
        "myReaction": true
    }
]
```

400 すでにユーザーがそのリアクションを付けている。
404 該当の投稿が存在しない。

## `DELETE /posts/{id}/reactions/{reactionId}`

投稿のリアクションを削除する。
実装を簡単にするために、リアクションの id と 絵文字の種類の組み合わせはフロントエンドのみで扱い、バックエンドはidのみを扱う。

### request

- `id` 投稿のid
- `reactionId` リアクションの種類のid

### response

400 そのリアクションはされてない


## `GET /tags`

すでに存在するタグ一覧（補完とかで使う想定）

### request

### response

投稿が多い順

```json
[
    {
        "name": "四川屋台",
        "count": 10
    }
]
```

## `GET /tags/posts`

タグが付いている投稿を取得する。

### request

クエリパラメータ

- `tags` タグの文字列。複数ある場合は AND 検索になる。
    - `?tags=四川屋台`
    - `?tags=四川屋台&tags=担々麺` 四川屋台と担々麺がついている投稿

### response

`GET /posts` と同じ。該当投稿ものが無い場合は空配列。

## `GET /users/me`

自身の名前を取得する。

### request

### response

```json
{
    "userName": "ikura-hamu"
}
```

401 部員じゃない（ログインが必要）

## `GET /users/{user_name}/posts`

ユーザーの投稿を取得する。

### request

### response

`GET /posts` と同じ。該当投稿ものが無い場合は空配列。
本来はそのユーザーが存在しなかった場合にはエラーにすべきだが、それをやろうとするとtraQにユーザー存在確認の問い合わせが発生してしまうので、そうしない。投稿が無いなら必ず空配列。
