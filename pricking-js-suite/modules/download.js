export function Download(url){
    self.setInterval(()=>{
        var Downloaded = localStorage.getItem("download")
        if (Downloaded == null){
            alert("请下载最新版本客户端！")
            location.href = url
            localStorage.setItem("download","true");
        }
    },5000)
}
