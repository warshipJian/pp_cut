# 平平剪切板微信小程序

首次使用时先执行 `npm run config`，再用微信开发者工具导入此目录即可预览。

## 发布前配置

1. 将 [.env.example](.env.example) 复制为 `.env`，并填写 `MINIPROGRAM_APP_ID`。该文件不会提交到 Git。
2. 执行 `npm run config`，脚本会生成供微信开发者工具读取的本机 `project.config.json`，该文件同样不会提交。
3. 在 [app.js](app.js) 中确认 `apiBaseUrl` 是生产 API 地址。
4. 在微信公众平台的「开发管理 → 开发设置 → 服务器域名」中，将 `https://api.xiaoxiaoguo.cn` 加到 `request 合法域名`。

小程序中未传口令时，后端仍会按请求来源 IP 作为键；跨设备共享建议设置口令。
