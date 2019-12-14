import Vue from 'vue'
import App from './App.vue'
import BaiduMap from 'vue-baidu-map'
import ViewUI from 'view-design';
import 'view-design/dist/styles/iview.css';
import utils from './js/utils.js'
import router from './js/router.js'

import axiosConfig from './js/axiosConfig.js'

Vue.config.productionTip = false
Vue.config.devtools = true;

Vue.use(ViewUI);
Vue.use(utils)
Vue.use(axiosConfig)
Vue.use(BaiduMap, {
  // ak 是在百度地图开发者平台申请的密钥 详见 http://lbsyun.baidu.com/apiconsole/key */
  ak: 'UuWGDj774aFbBhXgtSFMgvKISROhfcAy'
})

new Vue({
  router,
  render: h => h(App),
}).$mount('#app')
