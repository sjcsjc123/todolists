const utils= require('../../utils/util')
const app = getApp()

Page({
  data: {
    "phone": ""
  },
  getPhoneNumber (e) {
    var code = e.detail.code;
    //console.log(code)
    wx.request({
      url: 'http://localhost/getPhone',
      method: 'POST',
      header:{
        "code": code
      },
      success: function(res){
        utils.wxLogin(res.data.data)
        //console.log(res.data.data)
      }
    })
    console.log("finish")
  },
})
