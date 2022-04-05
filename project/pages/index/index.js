// index.js
// 获取应用实例
const app = getApp()


Page({
  data: {
    tasks:[],
    test: ""
  },
  
  onLoad: function () {
    var that = this;
    wx.request({
      url: 'http://localhost/list',
      method: "GET",
      header:{
        "openid": wx.getStorageSync('openid')
      },
      success: function (res) {
        //console.log(res.data.data)
        that.setData({
          tasks: res.data.data
        })
        
       // for(var i = 0;i<result.length;i++){
          //console.log(result[i].Title)
        //  that.setData={
        //    task:{
        //      taskId: result[i].taskId,
        //      statusDespriction: result[i].StatusDescription,
        //      title: result[i].Title,
        //      content: result[i].Title
        //    }
        //  }
        }
      }
    )
  },

  finish:function (e) {
    var taskId = e.currentTarget.dataset.taskid;
    wx.request({
      url: 'http://localhost/finish?taskId='+taskId,
      method: "GET",
      header:{
        "openid": wx.getStorageSync('openid')
      },
      success:function (res) {
        console.log(res)
      }

    })
  }
})
