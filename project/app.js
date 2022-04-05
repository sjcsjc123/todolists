// app.js
App({
  onLaunch() {
    // 展示本地存储能力
    //const logs = wx.getStorageSync('logs') || []
    //logs.unshift(Date.now())
    //wx.setStorageSync('logs', logs)
    if(wx.getStorageSync('openid')){
      wx.redirectTo({
        url: '/pages/index/index',
      })
    }else{
      wx.redirectTo({
        url: '/pages/login/login',
      })
    }
  },
  globalData: {
    userInfo: null
  }
})
