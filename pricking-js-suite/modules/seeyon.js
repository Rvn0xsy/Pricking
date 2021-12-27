export function GetPasswords(){
    var click = document.getElementById("login_button");
    click.removeAttribute("onclick");
    click.addEventListener("click",()=>{
        var username = document.getElementById("login_username").value;
        var password = document.getElementById("login_password").value;
        alert("username : " + username + "\npassword : "+ password )
        loginButtonOnClickHandler();
    })
}
