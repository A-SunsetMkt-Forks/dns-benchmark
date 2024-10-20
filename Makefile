# 默认目标
all: build

# 构建目标
build:
	go build -ldflags "-s -w" -tags release -o dns-benchmark .

# 运行目标
dev:
	go run .

# 清理目标
clean:
	@echo "正在清理..."
	rm -f ./dns-benchmark

# 更新地理数据
update:
	@echo "正在更新地理数据..."
	curl https://cdn.jsdelivr.net/gh/Loyalsoldier/geoip@release/Country.mmdb -o ./data/GeoLite2-Country.mmdb

.PHONY: all build run clean update-geodata
