# Epformat

Epformat 是一个用于重命名下载的番剧文件名称的工具，支持根据格式化规则重命名单个或批量的文件，同时提供命令行和图形用户界面（GUI）。

## 功能

- **格式化文件名**：通过指定格式规则，将文件名格式化为标准的剧集命名方式。
- **批量重命名**：支持对文件夹内所有文件进行批量重命名。
- **灵活配置**：支持通过命令行参数或配置文件设置标题、季数、集数等信息。
- **GUI 支持**：简单直观的图形界面，支持选择文件/文件夹、设置格式并直接重命名。
- **用户交互**：可以选择是否在重命名前确认操作。

## 使用方法

### 安装

确保已安装 Go 环境和 Fyne 库。然后克隆仓库并构建项目：

```bash
git clone https://github.com/yourusername/epformat.git
cd epformat
go build -o epformat
```

### 命令说明

#### 启动 GUI 模式

```bash
epformat gui
```

启动图形界面后，可以通过以下步骤使用：
1. 点击 "Select" 按钮选择文件或文件夹。
2. 填写 `Format` 字段（如 `S{season:02}E{episode:02} - {title}`）。
3. 可选：填写 `Title`、`Season` 和 `Episode` 字段。
4. 点击 `Apply` 查看预览。
5. 确认无误后点击 `Rename` 按钮完成重命名。

#### 格式化文件名（命令行模式）

```bash
epformat format <name>... [flags]
```

**参数：**
- `<name>`: 要格式化的文件或文件夹名称。

**选项：**
- `-t, --title <title>`: 指定剧集标题。
- `-f, --format <format>`: 指定格式化字符串（默认为 `S{season:02}E{episode:02} {title}`）。
- `-s, --season <season>`: 指定季数。
- `-e, --episode <episode>`: 指定集数。
- `-v, --verbose`: 显示详细信息。

#### 批量重命名

```bash
epformat rename <file/directory>... [flags]
```

**参数：**
- `<file/directory>`: 要重命名的文件或文件夹路径。

**选项：**
- `-t, --title <title>`: 指定剧集标题。
- `-f, --format <format>`: 指定格式化字符串。
- `-s, --season <season>`: 指定季数。
- `-e, --episode <episode>`: 指定集数。
- `-y, --yes`: 不确认直接重命名。

### 示例

#### 启动 GUI

```bash
epformat gui
```

#### 格式化文件名

```bash
epformat format "episode_01.mp4" -t "My Anime" -s 1 -e 1
```

输出：
```
S01E01 - My Anime.mp4
```

#### 批量重命名文件夹中的所有文件

```bash
epformat rename ./downloads -t "My Anime" -s 1
```