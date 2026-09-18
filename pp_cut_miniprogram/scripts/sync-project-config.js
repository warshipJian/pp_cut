const fs = require('fs')
const path = require('path')

const root = path.resolve(__dirname, '..')
const envPath = path.join(root, '.env')
const templatePath = path.join(root, 'project.config.template.json')
const outputPath = path.join(root, 'project.config.json')

if (!fs.existsSync(envPath)) {
  throw new Error('未找到 .env；请先由 .env.example 创建，并填写 MINIPROGRAM_APP_ID。')
}

const env = Object.fromEntries(
  fs.readFileSync(envPath, 'utf8')
    .split(/\r?\n/)
    .map(line => line.trim())
    .filter(line => line && !line.startsWith('#') && line.includes('='))
    .map(line => {
      const index = line.indexOf('=')
      return [line.slice(0, index).trim(), line.slice(index + 1).trim()]
    })
)

const appId = env.MINIPROGRAM_APP_ID
if (!/^wx[a-zA-Z0-9]{16}$/.test(appId || '')) {
  throw new Error('MINIPROGRAM_APP_ID 格式无效，应为 wx 开头的 18 位 AppID。')
}

const config = JSON.parse(fs.readFileSync(templatePath, 'utf8'))
config.appid = appId
fs.writeFileSync(outputPath, `${JSON.stringify(config, null, 2)}\n`)
console.log('已生成 project.config.json（该文件不会被 Git 提交）。')
