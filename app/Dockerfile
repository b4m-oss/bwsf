FROM golang:1.24-alpine

WORKDIR /app

# Node.jsとnpmをインストール（bwコマンド用）
RUN apk add --no-cache nodejs npm

# Bitwarden CLIをグローバルにインストール
RUN npm install -g @bitwarden/cli

# 開発に必要なツールをインストール
RUN apk add --no-cache git make

# Goモジュールのキャッシュを有効活用するため、go.modを先にコピー
COPY go.mod ./
RUN go mod download

# ソースコードをコピー
COPY . .

# 開発用のデフォルトコマンド（docker-composeで上書き可能）
CMD ["go", "run", "src/main.go"]

