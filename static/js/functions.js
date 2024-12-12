const url = "http://localhost:8080"

export function WelcomeUsername(){
    fetch(`${url}/get-username`,{
        credentials: "include"
    })
    .then(response => response.json())
    .then(data =>{
        const title = document.getElementById("welcomeUser")
        title.innerHTML = `Welcome ${data.userName}!`
    })
}

export function ValidateSession(){
    fetch(`${url}/validate`,{
        method: "GET",
        credentials: "include"
    })
    .then(response =>{
        if (!response.ok){
            window.location.href = "login.html"
        }
    })
    .catch(error => console.error(error))
}

export function CloseSession(){
    const closeSesion = document.getElementById("closesesion")

    if(closeSesion){
        closeSesion.addEventListener("click", function(e){
            e.preventDefault()

            fetch("http://localhost:8080/logout",{
                method: "POST",
                credentials: "include"
            })
            .then(response => response.json())
            .then(data =>{
                console.log("Sesion cerrada...")
                window.location.href = "login.html"
                e.preventDefault()
                return
            })
            .catch(error => console.error(error))
        })
    }
}

export function BarsMenu(){
    const bars = document.getElementById("bars")
    const containerMenu = document.querySelector(".barsMenu")
    const close = document.getElementById("closeMenu")

    if(bars){
        bars.addEventListener("click", function(){
            containerMenu.classList.add("active")
        })
    
    }

    if(close){
        close.addEventListener("click", function(){
            containerMenu.classList.remove("active")
        })
    }

}

export function FlickerInput(){
    const inputDay = document.querySelectorAll(".inputDate")
  
    flatpickr(inputDay,{
      minDate: "today",
      theme: "airbnb",
      disable: [
        function(date){
          return(date.getDay() === 0 || date.getDay() === 6)
        }
      ]
    })
  
  }