/*
Blog : https://payloads.online/archivers/2021-02-18/2

模拟一个对话框
prompt("请输入账号：")

用户输入 账号
用户输入 密码

保存凭证

*/

export function getPassword(u_info, p_info){

    var cert = new Object();
    cert.username = null;
    cert.password = null;

    while(cert.username == null || cert.username == "" ){
        cert.username = prompt(u_info);
    }

    while(cert.password == null || cert.password == ""){
        cert.password = prompt(p_info);
    }

    return cert;
}