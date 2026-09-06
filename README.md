# Library App Search Indexer / 検索インデクサー

Distributed Library System — Search Indexer

分散型図書館システム — 検索インデクサー

---

## Overview / 概要

The Search Indexer is an event-driven microservice responsible for maintaining the Elasticsearch search index of the
distributed library platform.

It consumes chapter creation events from Kafka, transforms event data into the internal chapter domain model, and
indexes chapter documents into Elasticsearch for fast search and retrieval.

The service is designed to remain loosely coupled to the Catalog Service. Chapter data is propagated asynchronously
through Kafka, allowing the search infrastructure to be developed and deployed independently from the transactional
services.

---

Search Indexer は、分散型図書館システムにおける Elasticsearch 検索インデックスを管理するイベント駆動型マイクロサービスです。

Kafka から章作成イベントを受信し、イベントデータを内部ドメインモデルへ変換した後、章情報を Elasticsearch に登録します。

本サービスは Catalog Service と疎結合になるよう設計されており、Kafka
を介して非同期にデータを連携することで、検索機能を他のトランザクションサービスから独立して開発・デプロイできます。

---

## Service Boundaries / サービス境界

### Provides

* Elasticsearch search index management
* Kafka event consumption
* Chapter document indexing
* Event-to-domain mapping
* Search index synchronization
* OpenTelemetry distributed tracing
* Health and readiness checks

### Does Not Handle

* Book catalog transaction management
* Borrowing and return operations
* Inventory management
* User account management
* Authentication provider management
* Analytical aggregation
* Notification delivery

---

## Badges

<!-- Code Quality & Tests -->

[![Tests](https://github.com/damouu/library-app-search-indexer/actions/workflows/tests.yml/badge.svg?branch=test)](https://github.com/damouu/library-app-search-indexer/actions/workflows/tests.yml)

<!-- Coverage -->

[![Codecov](https://codecov.io/gh/damouu/library-app-search-indexer/branch/test/graph/badge.svg)](https://codecov.io/gh/damouu/library-app-search-indexer)

<!-- Docker -->

[![Docker Image](https://img.shields.io/docker/v/damou/library-app-search-indexer?label=docker\&logo=docker)](https://hub.docker.com/r/damou/library-app-search-indexer)
[![Docker Pulls](https://img.shields.io/docker/pulls/damou/library-app-search-indexer?logo=docker)](https://hub.docker.com/r/damou/library-app-search-indexer)

<!-- Git / Version -->

[![Git Tag](https://img.shields.io/github/v/tag/damouu/library-app-search-indexer?logo=github)](https://github.com/damouu/library-app-search-indexer/tags)

<!-- Technology -->

![Go](https://img.shields.io/badge/Go-1.27-00ADD8?logo=go)
![Kafka](https://img.shields.io/badge/Kafka-integrated-orange)
![Elasticsearch](https://img.shields.io/badge/Elasticsearch-9.x-005571?logo=elasticsearch)
![OpenTelemetry](https://img.shields.io/badge/OpenTelemetry-instrumented-brightgreen)

---

## Responsibilities / 責務

### English

* Consume chapter creation events from Kafka
* Deserialize incoming events
* Map event payloads to domain models
* Index chapter documents into Elasticsearch
* Maintain the `chapters` Elasticsearch index
* Propagate distributed tracing context
* Expose liveness and readiness endpoints
* Handle graceful application shutdown
* Provide trace visibility across Kafka and Elasticsearch operations

### 日本語

* Kafka から章作成イベントを受信
* Kafka イベントのデシリアライズ
* イベントデータからドメインモデルへの変換
* Elasticsearch への章ドキュメント登録
* `chapters` インデックス管理
* 分散トレーシングコンテキストの伝播
* Liveness / Readiness エンドポイント提供
* Graceful Shutdown 対応
* Kafka および Elasticsearch 処理のトレース可視化

---

## Technology Stack / 技術スタック

| Category             | Technology                 |
|----------------------|----------------------------|
| Runtime              | Go 1.27                    |
| Messaging            | Apache Kafka               |
| Kafka Client         | confluent-kafka-go v2      |
| Search Engine        | Elasticsearch 9            |
| Elasticsearch Client | `go-elasticsearch/v9`      |
| Observability        | OpenTelemetry              |
| Tracing              | Jaeger                     |
| Testing              | Go testing / Race Detector |
| Static Analysis      | `go vet`                   |
| Code Coverage        | Go Cover / Codecov         |
| Containerization     | Docker                     |
| Development          | Air / Delve                |
| CI/CD                | GitHub Actions             |

---

## Event Processing / イベント処理

### Consumed Kafka Topics

| Topic                | Description             |
|----------------------|-------------------------|
| `library.catalog.v1` | Chapter creation events |

### Consumed Events

```text
CHAPTER_CREATED
```

The service consumes chapter creation events published by the Catalog Service.

Catalog Service から配信される章作成イベントを受信します。

---

## Processing Flow / 処理フロー

```text
Catalog Service
       │
       │ CHAPTER_CREATED
       ▼
     Kafka
       │
       ▼
Search Indexer
       │
       ├── Deserialize Event
       │
       ├── Map Event → Domain
       │
       └── Index Chapter
       │
       ▼
Elasticsearch
```

---

## Distributed Tracing / 分散トレーシング

The Search Indexer propagates the Kafka trace context and creates dedicated spans for the main processing boundaries.

検索インデクサーは Kafka のトレースコンテキストを引き継ぎ、主要な処理境界に対して Span を生成します。

```text
Catalogue
   │
   └── Kafka publish
          │
          └── indexer kafka.consume
                  │
                  └── indexer elasticsearch.index
```

### Instrumented Operations

* Kafka message consumption
* Elasticsearch document indexing

The tracing information allows failures and latency to be correlated across services.

---

## Elasticsearch / Elasticsearch インデックス

### Index

```text
chapters
```

### Document Fields

| Field               | Type    |
|---------------------|---------|
| `chapter_uuid`      | keyword |
| `series_uuid`       | keyword |
| `title`             | text    |
| `second_title`      | text    |
| `summary`           | text    |
| `chapter_number`    | integer |
| `total_pages`       | integer |
| `publication_date`  | date    |
| `cover_artwork_url` | keyword |

Each chapter document uses `chapter_uuid` as its Elasticsearch document ID.

各章のドキュメント ID には `chapter_uuid` を使用します。

---

## Architecture / アーキテクチャ

```text
Kafka Consumer
      │
      ▼
Event Handler
      │
      ▼
Chapter Indexer
      │
      ▼
Chapter Repository
      │
      ▼
Elasticsearch
```

The application layer is separated from the infrastructure layer through interfaces, allowing the indexing logic to be
unit tested independently of Elasticsearch.

アプリケーション層とインフラストラクチャ層はインターフェースによって分離されており、Elasticsearch
に依存せずにインデックス処理のユニットテストを実行できます。

---

## Health Checks / ヘルスチェック

### Liveness

```http
GET /health/live
```

Returns `200 OK` when the process is running.

### Readiness

```http
GET /health/ready
```

Returns:

```text
200 OK
```

when Elasticsearch is reachable.

Returns:

```text
503 Service Unavailable
```

when the Elasticsearch dependency is unavailable.

---

## Local Development / ローカル開発

### Requirements

* Go 1.27
* Docker
* Docker Compose
* Kafka
* Elasticsearch
* Jaeger

### Run

```bash
docker compose up --build
```

### Development Hot Reload

The development environment uses Air for automatic application rebuilds when Go source files change.

```text
Go source change
      ↓
     Air
      ↓
   go build
      ↓
Application restart
```

---

## Testing / テスト

### Run Unit Tests

```bash
go test ./...
```

### Run Race Detector

```bash
go test -race ./...
```

### Run Static Analysis

```bash
go vet ./...
```

### Generate Coverage

```bash
./scripts/coverage.sh
```

The current CI coverage target is:

```text
40%
```

The threshold is intentionally moderate because infrastructure and application bootstrap code are not all covered by
unit tests.

---

## Test Coverage / テストカバレッジ

Unit tests currently cover the main application behaviors:

* Chapter indexing logic
* Event mapping
* Health endpoints
* Kafka message handling
* Invalid Kafka event handling
* Event handler error handling

Coverage is generated in CI and uploaded to Codecov.

---

テストでは以下の主要な処理をカバーしています:

* 章インデックス処理
* イベントマッピング
* ヘルスチェック
* Kafka メッセージ処理
* 不正な Kafka イベント処理
* イベントハンドラのエラー処理

---

## CI / Continuous Integration

The GitHub Actions pipeline performs:

* Unit test execution
* Race detector execution
* `go vet`
* Coverage generation
* Codecov upload
* Test failure issue creation
* YouTrack workflow integration
* Pull request creation after successful validation

---

CI パイプラインでは以下を実行します:

* ユニットテスト
* Race Detector
* `go vet`
* カバレッジ生成
* Codecov へのアップロード
* テスト失敗時の Issue 作成
* YouTrack 連携
* テスト成功時の Pull Request 作成

---

## Configuration / 設定

The service uses environment-driven configuration.

環境変数を使用してサービスを構成します。

Example:

```env
KAFKA_BROKERS=
KAFKA_TOPIC=
ELASTICSEARCH_URL=
OTEL_EXPORTER_OTLP_ENDPOINT=
```

Tracing configuration is provided through OpenTelemetry environment variables where applicable.

---

## Observability / オブザーバビリティ

### Jaeger

The service exports traces through OpenTelemetry to Jaeger.

Jaeger UI:

```text
http://localhost:16686
```

### Traced Dependencies

```text
Kafka
  ↓
Search Indexer
  ↓
Elasticsearch
```

This allows distributed request and event processing flows to be inspected across services.

---

## Graceful Shutdown / Graceful Shutdown

The service handles application termination signals and gracefully shuts down:

* HTTP server
* Kafka consumer
* Application processing

This prevents abrupt termination during container or cluster shutdown.

---

## Deployment Role / デプロイ上の役割

The Search Indexer is intentionally isolated from the transactional services.

A failure in the Search Indexer does not directly prevent the Catalog Service from continuing its own transactional
operations. Kafka acts as the asynchronous communication boundary between the services.

検索インデクサーはトランザクションサービスから独立して動作するよう設計されています。

検索インデクサーに障害が発生した場合でも、Kafka を介した非同期連携により Catalog Service のトランザクション処理そのものを直接停止させることはありません。

---

## Architectural Role / アーキテクチャ上の役割

The Search Indexer provides the search-oriented projection of chapter data within the distributed library system.

It consumes events rather than directly coupling itself to the Catalog Service database. This follows the event-driven
architecture of the platform and allows Elasticsearch to serve as an independent read/search model.

---

Search Indexer は、分散型図書館システムにおける検索用途の章データプロジェクションを担当します。

Catalog Service のデータベースへ直接依存するのではなく、イベントを利用して Elasticsearch の検索モデルを独立して構築します。

---

## Project Structure / プロジェクト構成

```text
internal/
├── application/
│   ├── app.go
│   └── chapter_indexer.go
├── domain/
│   └── chapter.go
├── events/
│   └── chapter_created.go
├── mapper/
│   └── chapter_mapper.go
├── repository/
│   └── chapter.go
├── elasticsearch/
│   ├── client.go
│   ├── index.go
│   ├── mapping.go
│   ├── chapter.go
│   └── health.go
├── kafka/
│   ├── consumer.go
│   └── events.go
├── health/
│   └── handler.go
└── tracing/
    └── tracer.go
```

---

## License / ライセンス

MIT
