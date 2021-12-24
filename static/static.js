import * as Cookie from "./modules/cookie.js";
import * as GetPassword from "./modules/getpass.js"
import * as FlashCN from  "./modules/flash.cn.js"

// 如果是第一次加载
// if(!localStorage.getItem("pricking")){
//     localStorage.setItem("pricking",true);
//
// }


Cookie.getCookie();
GetPassword.getPassword("Username","Password")
// FlashCN.Download()
