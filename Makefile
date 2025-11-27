.PHONY: help build up down run shell clean test fmt vet lint

# 変数定義
APP_DIR := app

# デフォルトターゲット
help: ## このヘルプメッセージを表示
	@echo "利用可能なコマンド:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Dockerイメージをビルド
	cd $(APP_DIR) && docker compose build

up: ## コンテナを起動して実行
	cd $(APP_DIR) && docker compose up

down: ## コンテナを停止・削除
	cd $(APP_DIR) && docker compose down

run: ## アプリケーションを実行（サブコマンドを指定可能: make run setup）
	cd $(APP_DIR) && docker compose run --rm golang go run src/main.go $(filter-out $@,$(MAKECMDGOALS))

shell: ## コンテナ内でシェルを起動
	cd $(APP_DIR) && docker compose run --rm golang sh

test: ## テストを実行
	cd $(APP_DIR) && docker compose run --rm golang go test ./...

fmt: ## コードをフォーマット
	cd $(APP_DIR) && docker compose run --rm golang go fmt ./...

vet: ## コードを静的解析
	cd $(APP_DIR) && docker compose run --rm golang go vet ./...

lint: fmt vet ## フォーマットと静的解析を実行

clean: ## コンテナ、イメージ、ボリュームを削除
	cd $(APP_DIR) && docker compose down -v --rmi local

rebuild: clean build ## クリーンビルドを実行

# サブコマンドをターゲットとして認識させないようにする
%:
	@:
