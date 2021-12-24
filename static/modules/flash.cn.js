

console.log("Loaded flash.cn Module ...");

export function Download() {
    var loadLink = document.getElementsByClassName("loadLink");
    loadLink[0].classList.remove("disable")
    loadLink[0].addEventListener('click',function(){
        alert("Hello Pricking!");
    },false)
}

