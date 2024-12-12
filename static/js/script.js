console.log("Conectado...")
import { WelcomeUsername } from '/static/js/functions.js';
import { ValidateSession } from '/static/js/functions.js';
import { CloseSession } from '/static/js/functions.js';
import { BarsMenu } from '/static/js/functions.js';

const url = "http://localhost:8080"


function CountPatients(){
    const contadorPatients = document.getElementById("totalPatients")
        
    if(contadorPatients){
        fetch(`${url}/total-patients`,{
            credentials: "include",
        })
        .then(response =>{
            if(!response.ok){
                console.error(response.json())
            }
            return response.json()
        })
        .then(data =>{
            contadorPatients.innerHTML = data.total
        })
        .catch(error => console.error(error))
    }
}

function AppointmentToday(){
    const cantCitas = document.getElementById("totalCitas")

    if(cantCitas){
        fetch(`${url}/appointment-today`,{
            credentials: "include",
        })
        .then(response =>{
            if(!response.ok){
                console.error(response.json())
            }
            return response.json()
        })
        .then(data =>{
            cantCitas.innerHTML = data.count
        })
        .catch(error => console.error(error))
    }
}

function AppointmentWeek(){
    const cantAppointments = document.getElementById("appointmentsWeek")

    fetch(`${url}/appointments-week`,{
        credentials: "include",
    })
        .then(response =>{
            if(!response.ok){
                console.error(response.json())
            }
            return response.json()
        })
        .then(data =>{
            cantAppointments.innerHTML = data.count
        })
        .catch(error => console.error(error))
}

function NextAppointments(){
    fetch(`${url}/next-appointments`,{
        credentials: "include",
    })
        .then(response =>{
            if(!response.ok){
                console.error(response)
            }
            return response.json()
        })
        .then(data => {
            console.log(data)
            const tablaBody = document.querySelector(".tabla tbody")
            tablaBody.innerHTML = ""

            data.forEach(patient =>{
                const row = document.createElement("tr")

                const date = new Date(patient.Fecha)
                const dateParsed = date.toLocaleDateString()
                
                row.innerHTML = `
                    <td>${patient.Paciente.FullName}</td>
                     <td>${patient.Paciente.Email}</td>
                     <td>${patient.Paciente.Dni}</td>
                     <td>${dateParsed}</td>
                     <td>${patient.Hora}</td>
              
                `;
                tablaBody.appendChild(row)
            })
        })
        .catch(error => console.error("Error", error))
}

document.addEventListener("DOMContentLoaded",function(e){
    e.preventDefault();

    ValidateSession()

    CountPatients()

    AppointmentToday()

    CloseSession()

    AppointmentWeek()

    BarsMenu()

    WelcomeUsername()

    NextAppointments()
})


