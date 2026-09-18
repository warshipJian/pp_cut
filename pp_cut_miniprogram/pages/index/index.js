const app = getApp()

const request = (options) => new Promise((resolve, reject) => {
  wx.request({
    ...options,
    success: (res) => res.statusCode >= 200 && res.statusCode < 300 ? resolve(res.data) : reject(new Error(`请求失败（${res.statusCode}）`)),
    fail: () => reject(new Error('网络连接失败，请稍后重试'))
  })
})

Page({
  data: { content: '', password: '', hasContent: false, getting: false, submitting: false },
  onLoad() {
    const draft = wx.getStorageSync('clipboardDraft') || {}
    this.setData({ content: draft.content || '', password: draft.password || '', hasContent: Boolean((draft.content || '').trim()) })
  },
  onUnload() { this.saveDraft() },
  onContentInput(e) { this.setData({ content: e.detail.value, hasContent: Boolean(e.detail.value.trim()) }) },
  onPasswordInput(e) { this.setData({ password: e.detail.value }) },
  saveDraft() { wx.setStorageSync('clipboardDraft', { content: this.data.content, password: this.data.password }) },
  goTv() { this.saveDraft(); wx.navigateTo({ url: '/pages/tv/index' }) },
  async getContent() {
    this.setData({ getting: true })
    try {
      const password = this.data.password.trim()
      const data = await request({ url: `${app.globalData.apiBaseUrl}/c/get`, data: password ? { password } : {} })
      this.setData({ content: data.content || '', hasContent: Boolean((data.content || '').trim()) })
      wx.showToast({ title: '内容获取成功', icon: 'success' })
    } catch (error) { wx.showToast({ title: error.message, icon: 'none' }) }
    finally { this.setData({ getting: false }) }
  },
  async submitContent() {
    const content = this.data.content
    if (!content.trim()) {
      const result = await new Promise(resolve => wx.showModal({ title: '确认操作', content: '空内容会清空远端，是否继续？', success: res => resolve(res.confirm) }))
      if (!result) return
    }
    this.setData({ submitting: true })
    try {
      const password = this.data.password.trim()
      await request({ url: `${app.globalData.apiBaseUrl}/c/submit`, method: 'POST', header: { 'content-type': 'application/json' }, data: password ? { content, password } : { content } })
      this.saveDraft()
      wx.showToast({ title: content ? '提交内容成功' : '远端内容已清空', icon: 'success' })
    } catch (error) { wx.showToast({ title: error.message, icon: 'none' }) }
    finally { this.setData({ submitting: false }) }
  },
  async clearRemote() {
    const result = await new Promise(resolve => wx.showModal({ title: '确认清空', content: '清空内容会把远端内容也清空，是否继续？', confirmColor: '#e74c3c', success: res => resolve(res.confirm) }))
    if (!result) return
    this.setData({ submitting: true })
    try {
      const password = this.data.password.trim()
      await request({ url: `${app.globalData.apiBaseUrl}/c/submit`, method: 'POST', header: { 'content-type': 'application/json' }, data: password ? { content: '', password } : { content: '' } })
      this.setData({ content: '', hasContent: false })
      this.saveDraft()
      wx.showToast({ title: '远端内容已清空', icon: 'success' })
    } catch (error) { wx.showToast({ title: error.message, icon: 'none' }) }
    finally { this.setData({ submitting: false }) }
  },
  copyContent() {
    if (!this.data.content.trim()) return wx.showToast({ title: '文本内容为空，无法复制', icon: 'none' })
    wx.setClipboardData({ data: this.data.content, success: () => wx.showToast({ title: '已复制到剪切板', icon: 'success' }) })
  }
})
