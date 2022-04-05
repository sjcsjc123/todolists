const formatTime = date => {
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()
  const hour = date.getHours()
  const minute = date.getMinutes()
  const second = date.getSeconds()

  return `${[year, month, day].map(formatNumber).join('/')} ${[hour, minute, second].map(formatNumber).join(':')}`
}

const formatNumber = n => {
  n = n.toString()
  return n[1] ? n : `0${n}`
}

module.exports = {
  formatTime,
  wxLogin,
}
//用户授权获取手机号后，登录微信，存入数据库
function wxLogin(phone) {
  wx.login({
    success: res => {
      // 发送 res.code 到后台换取 openId, sessionKey, unionId
      var code = res.code
      //console.info(code)
      wx.request({
        url: 'http://localhost/oauth',
        method: "post",
        data: {
          "code": code,
          "phone": phone
        },
        header: {
          "content-type": "application/json"
        },
        success: function (res) {
          console.log(res.data.data)
          //console.log("openid:"+res.data.data.openid)
          wx.setStorageSync('openid', res.data.data)
          wx.redirectTo({
            url: '/pages/index/index',
          })
        }
      })
    }
  })
}
