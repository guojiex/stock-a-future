# 收藏数据目录

这个目录用于存储用户的收藏股票数据和分组信息。

## 文件说明

- `favorites.json` - 用户收藏的股票列表
- `groups.json` - 用户创建的分组信息

## 注意事项

- 这个目录下的所有文件都被 `.gitignore` 忽略，不会提交到版本控制
- 应用首次运行时会自动创建这些文件
- 删除这些文件会清空所有收藏数据，请谨慎操作

## 数据格式

### favorites.json
```json
[
  {
    "id": "唯一ID",
    "ts_code": "股票代码",
    "name": "股票名称",
    "start_date": "开始日期",
    "end_date": "结束日期",
    "group_id": "分组ID",
    "sort_order": "排序序号",
    "created_at": "创建时间",
    "updated_at": "更新时间"
  }
]
```

### groups.json
```json
[
  {
    "id": "分组ID",
    "name": "分组名称",
    "color": "分组颜色",
    "sort_order": "排序序号",
    "created_at": "创建时间",
    "updated_at": "更新时间"
  }
]
```
