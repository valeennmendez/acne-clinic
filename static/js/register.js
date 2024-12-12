console.log("Conectado...")

function Register(){
    const boxForm = document.getElementById("boxForm")

    boxForm.addEventListener("submit", function(e){
        e.preventDefault()

        const name = document.getElementById("name").value
        const email = document.getElementById("email").value
        const phone = document.getElementById("phone").value
        const password = document.getElementById("password").value
        const repeatpass = document.getElementById("repeatpass").value
        const errorText = document.getElementById("errorText");
        const accountCreated = document.querySelector(".box")

        console.log(name,email,password,repeatpass)

        if (password != repeatpass){
            errorText.innerHTML = "Passwords do not match"
            return
        }

        const bodyUser = {
            fullname: name,
            email: email,
            password: password,
            phone: phone
        }

        fetch(`http://localhost:8080/register`,{
            method: "POST",
            credentials: "include",
            headers:{
                "Content-Type": "application/json"
            },
            body: JSON.stringify(bodyUser)
        })
            .then(response =>{
                if(!response.ok){
                    console.error(response)
                }
                return response.json();
            })
            .then(data =>{
                if (data.message === "UserAdmin created succesfully, pending approval"){
                    console.log("Usuario creado...")
                    accountCreated.classList.add("active");
                }
            })
    })
}


document.addEventListener("DOMContentLoaded", function(){

    Register()

})