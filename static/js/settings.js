console.log("Conectado...")
import { WelcomeUsername } from '/static/js/functions.js';
import { ValidateSession } from '/static/js/functions.js';
import { CloseSession } from '/static/js/functions.js';
import { BarsMenu } from '/static/js/functions.js';

const url = "http://localhost:8080"

function GetAllAdmins(){

    fetch(`${url}/admins`,{
        method: "GET",
        credentials: "include"
    })
        .then(response =>{
            if(!response.ok){
                console.error(response)
            }
            return response.json()
        })
        .then(data =>{
            const tablaBody = document.querySelector(".tabla tbody")
            tablaBody.innerHTML = "";

            data.forEach(admin => {
                const row = document.createElement("tr")

                row.innerHTML = `
                    <td>${admin.FullName}</td>
                    <td>${admin.Email}</td>
                    <td>${admin.Phone}</td>
                    <td>${admin.Status}</td>
                    <td>
                        <ul>
                            <div class="buttons">
                                <a name="" id="checkBtn" class="btn approveBtn" data-id="${admin.ID}" href="#" role="button">
                                    <i class="fa-solid fa-check"></i>
                                    <span class="tooltip-text">Approve</span>
                                </a>
                            </div>
                            <div class="buttons">
                                <a name="" id="" class="btn declineBtn" href="#" data-id="${admin.ID}" role="button">
                                    <i class="fa-solid fa-x"></i>
                                    <span class="tooltip-text">Decline</span>
                                </a>
                            </div>
                        </ul>
                    </td>
                `;
                tablaBody.appendChild(row)

                ApproveAccount()

                DeclineAccount()

            });
        })
        .catch(error => console.error(error))

}

function ApproveAccount() {
    const checkButtons = document.querySelectorAll(".approveBtn");
        
    checkButtons.forEach(button => {
        button.addEventListener("click", function(e){
            e.preventDefault();

            const adminID = button.getAttribute("data-id");
            console.log(adminID);

            fetch(`${url}/approve-user/${adminID}`,{
                method: "POST",
                credentials: "include",
            })
            .then(response =>{
                if(!response.ok){
                    console.error(response)
                }
                return response.json()
            })
            .then(data =>{
                if(data.message === "user approved"){
                    location.reload()
                }
            })
            .catch(error => console.error(error))
 
        });
    });
}

function DeclineAccount() {
    const checkButtons = document.querySelectorAll(".declineBtn");
    
    checkButtons.forEach(button => {
        button.addEventListener("click", function(e){
            e.preventDefault();

            const adminID = button.getAttribute("data-id");

            fetch(`${url}/decline-user/${adminID}`,{
                method: "POST",
                credentials: "include",
            })
            .then(response =>{
                if(!response.ok){
                    console.error(response)
                }
                return response.json()
            })
            .then(data =>{
                if(data.message === "user decline"){
                    location.reload()
                }
            })
            .catch(error => console.error(error))
 
        });
    });
}


function GetRole(){
    fetch(`${url}/admin-role`,{
        credentials: "include"
    })
        .then(response => response.json())
        .then(data =>{
            const contenedor = document.getElementById("container")
            const message = document.getElementById("message")
            
            if(data != "root"){
                contenedor.classList.add("unauthorized")
                message.classList.add("unauthorized")
            }
        })
}

document.addEventListener("DOMContentLoaded", function(e){
    GetAllAdmins()

    ValidateSession()

    BarsMenu()

    WelcomeUsername()

    CloseSession()

    GetRole()
})