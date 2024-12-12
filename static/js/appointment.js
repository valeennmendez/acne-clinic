console.log("Conectado...")
import { WelcomeUsername } from '/static/js/functions.js';
import { ValidateSession } from '/static/js/functions.js';
import { CloseSession } from '/static/js/functions.js';
import { BarsMenu } from '/static/js/functions.js';
import { FlickerInput } from '/static/js/functions.js';

const url = "http://localhost:8080"

function GetAllAppointments(){
    fetch(`${url}/appointments`,{
        method: "GET",
        credentials: "include"
    })
        .then(response =>{
            if(!response.ok){
                console.error("ERROR")
            }
            return response.json()
        })
        .then(data =>{
            updateTable(data);
        })
        .catch(error => console.error("error: ",error))
}

function GetAppointmentsFilter(status){
    fetch(`${url}/appointments-filter?status=${status}`,{
        method: "GET",
        credentials: "include"
    })
        .then(response => response.json())
        .then(data =>{
            console.log(data)
            updateTable(data);
        })
        .catch(error => console.error(error))
}

function Filter(){
    const select = document.getElementById("selectFilter")
    select.addEventListener("change", function(){
        const value = select.value;
        console.log(value)
        if(value === "ALL"){
            GetAllAppointments()
        }else{
            GetAppointmentsFilter(value)
        }
    })
}

function updateTable(appointments) {
    const tableBody = document.querySelector(".tabla tbody");
    tableBody.innerHTML = ""; // Limpiar el contenido actual

    appointments.forEach(appointment => {
        const row = document.createElement("tr");
        const fecha = new Date(appointment.Fecha);
        const fechaFormateada = fecha.toLocaleDateString();

        row.innerHTML = `
            <td>${appointment.Paciente.FullName}</td>
            <td>${appointment.Paciente.Dni}</td>
            <td>${fechaFormateada}</td>
            <td>${appointment.Hora}</td>
            <td>${appointment.Status}</td>
            <td>
                <ul>
                <div class="buttons">
                    <a class="btn deleteBtn" data-id="${appointment.ID}" href="#" role="button">
                        <i class="fa-regular fa-x"></i>
                        <span class="tooltip-text">Cancel</span>
                    </a>
                </div>
                <div class="buttons">
                    <a class="btn editBtn" href="#" data-id="${appointment.ID}" role="button">
                        <i class="fa-regular fa-pen-to-square"></i>
                        <span class="tooltip-text">Edit</span>
                    </a>
                </div>
                <div class="buttons">
                    <a class="btn doneBtn" href="#" data-id="${appointment.ID}" role="button">
                        <i class="fa-solid fa-check"></i>
                        <span class="tooltip-text">Done</span>
                    </a>
                </div>
                </ul>
            </td>
        `;
        tableBody.appendChild(row);
    });

    document.querySelectorAll(".deleteBtn").forEach(button => {
        button.addEventListener("click", function (e) {
            e.preventDefault();
            const id = this.getAttribute("data-id");
            CancelAppointments(id);
        });
    });

    document.querySelectorAll(".doneBtn").forEach(button => {
        button.addEventListener("click", function (e) {
            e.preventDefault();
            const id = this.getAttribute("data-id");
            DoneAppointments(id);
        });
    });

    document.querySelectorAll(".editBtn").forEach(button =>{
        const formEdit = document.getElementById("formEditAppointment");
            button.addEventListener("click", function(e){
            const id = this.getAttribute("data-id")
            LoadDataAppointment(id);
            formEdit.classList.add("active");
        })
    })
}

function LoadDataAppointment(id){

    const formEdit = document.getElementById("formEditAppointment")
    const closeButton = document.getElementById("closeFormAppBtnEdit")

    closeButton.addEventListener("click", function(){
        formEdit.classList.remove("active")
    })

    fetch(`${url}/appointments/${id}`,{
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

            const fechaFormateada = new Date(data.Fecha).toISOString().split("T")[0];

            console.log(fechaFormateada)

            document.getElementById("namePatient").value = data.Paciente.FullName
            document.getElementById("dayAppEdit").value = fechaFormateada
            document.getElementById("hourAppEdit").value = data.Hora
            document.getElementById("motivoAppEdit").value = data.MotivoConsulta

            LoadHoursEdit(fechaFormateada)
        })
        .catch(error => console.error(error))

        formEdit.addEventListener("submit", function(){
            UpdateAppointment(id);
        })

}

function UpdateAppointment(id){

    const inputPatient = document.getElementById("namePatient")
    const day = document.getElementById("dayAppEdit").value
    const hour = document.getElementById("hourAppEdit").value
    const motivo = document.getElementById("motivoAppEdit").value
    const pacienteid = inputPatient.getAttribute("data-selected-id")

    const appointmentData = {
        pacienteid: parseInt(pacienteid),
        fecha: day + "T00:00:00Z",
        hora: hour,
        motivoconsulta: motivo,
    }

    fetch(`${url}/appointments-edit/${id}`,{
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(appointmentData)
    })
        .then(response => response.json())
        .then(data =>{
            console.log("appointment modificated.", data)
        })
        .catch(error => console.error(error))

}

function CancelAppointments(id){
    fetch(`${url}/cancel-appointment/${id}`,{
            method: "DELETE",
            credentials: "include",
        })
            .then(response => response.json())
            .then(data =>{
                console.log("entro data")
                if(data.message === "appointment cancelated successfully"){
                    GetAllAppointments();
                }
            })
            .catch(error => console.error(error))
}

function DoneAppointments(id){
    fetch(`${url}/done-appointment/${id}`,{
        method: "DELETE",
        credentials: "include",
    })
        .then(response => response.json())
        .then(data =>{
            console.log("entro data")
            if(data.message === "appointment doned successfully"){
                GetAllAppointments();
            }
        })
        .catch(error => console.error(error))
}

function SearchPatientsForm(){
    const inputSearch = document.getElementById("searchPatient");
    const autocompleteList = document.getElementById("autocomplete-list");

    inputSearch.addEventListener("input", function() {
        const query = this.value;

        if (query.length > 1) {
            fetch(`${url}/search-patient?p=${encodeURIComponent(query)}`,{
                credentials: "include",
            })
                .then(response => response.json())
                .then(data => {
                    autocompleteList.innerHTML = "";

                    data.forEach(patient => {
                    const listItem = document.createElement("li");

                    const dataid = patient.ID;
                    console.log(dataid)

                    listItem.textContent = patient.FullName;
                    listItem.setAttribute("data-id", dataid)
                    listItem.classList.add("listItem")
                    listItem.dataset.value = patient.FullName;
                    autocompleteList.appendChild(listItem);
                });

                autocompleteList.style.display = "block";
                })
                .catch(error => console.error("error", error));
        } else {
            autocompleteList.style.display = "none";
        }
    });

    autocompleteList.addEventListener("click", function(event) {
    if (event.target.tagName === "LI") {
    inputSearch.value = event.target.dataset.value;
    const selectedPatientID = event.target.getAttribute("data-id");

    inputSearch.setAttribute("data-selected-id", selectedPatientID);
    autocompleteList.style.display = "none";
    }
    });
}

function CloseOpenForm(openButton, closeButton, container){

    openButton.addEventListener("click", function(){
        container.classList.add("active")
    })

    closeButton.addEventListener("click", function(){
        container.classList.remove("active")
    })
}

function AddNewAppointment(){

    const appContainer = document.getElementById("formCreateAppointment")
    const inputPatient = document.getElementById("searchPatient")
    const openButton = document.getElementById("openFormAppBtn")
    const closeButton = document.getElementById("closeFormAppBtn")

    CloseOpenForm(openButton,closeButton,appContainer);
    LoadHours();

    appContainer.addEventListener("submit", function(e){
    e.preventDefault()

    const pacienteID = inputPatient.getAttribute("data-selected-id")
    const day = document.getElementById("dayApp").value
    const hour = document.getElementById("hourApp").value
    const motivo = document.getElementById("motivoApp").value


    const appointmentData = {
        pacienteid: parseInt(pacienteID),
        fecha: day,
        hora: hour,
        motivoconsulta: motivo,
    }

    console.log(appointmentData)

    fetch(`${url}/create-appointment`,{
        method: "POST",
        headers:{
        "Content-Type": "application/json",
        },
        body: JSON.stringify(appointmentData),
        credentials: "include"
    })
    .then(response =>{
        if(!response.ok){
            console.error(response)
        }
        return response.json()
    })
    .then(data =>{
        console.log("Appointment created...")
        if(data.message === "Appoinment created succesfully"){
            alert("Cita creada correctamente.")
            window.location.reload()
        }
    })  
    .catch(error => console.error(error))
    })

}

function LoadHours(){
    const dayInput = document.getElementById("dayApp")
    const hourSelect = document.getElementById("hourApp")

    dayInput.addEventListener("change", function(){
        const selectedDay = dayInput.value;

        if (selectedDay){
            fetchAvilableHours(selectedDay)
        }
    })


    function fetchAvilableHours(date){
        fetch(`${url}/available-hours?fecha=${date}`,{
            credentials: "include"
        })
            .then(response => response.json())
            .then(data => {
                console.log(data)
                const availableHours = data.available_hours;
                updateHourOptions(availableHours);
            })
            .catch(error => console.error("Error: ", error))
    }

    function updateHourOptions(hours){
        hourSelect.innerHTML = "";
        console.log(hours)

        if(hours.length > 0){
            hours.forEach( hour => {
                const option = document.createElement('option')
                option.value = hour;
                option.textContent = hour
                hourSelect.appendChild(option)
            })
        }else {
            const option = document.createElement('option');
            option.textContent = "No available hours";
            hourSelect.appendChild(option)
        }
    }
}

function LoadHoursEdit(selectedDay){
    const hourSelect = document.getElementById("hourAppEdit")
    if (selectedDay){
        fetchAvilableHours(selectedDay)
    }

    function fetchAvilableHours(date){
        fetch(`${url}/available-hours?fecha=${date}`,{
        credentials: "include"
        })
            .then(response => response.json())
            .then(data => {
                console.log(data)
                const availableHours = data.available_hours;
                updateHourOptions(availableHours);
            })
            .catch(error => console.error("Error: ", error))
    }

    function updateHourOptions(hours){
        hourSelect.innerHTML = "";
        console.log(hours)

        if(hours.length > 0){
            hours.forEach( hour => {
                const option = document.createElement('option')
                option.value = hour;
                option.textContent = hour
                hourSelect.appendChild(option)
            })
        }else{
            const option = document.createElement('option');
            option.textContent = "No available hours";
            hourSelect.appendChild(option)
        }
    }
}

function SearchPatient(){
    const searchInput = document.getElementById("searchInput");

    searchInput.addEventListener("input", function(){
        const query = searchInput.value

        if(query.length > 0){
            fetch(`${url}/search-appointments?p=${encodeURIComponent(query)}`,{
                credentials: "include"
            })
                .then(response => response.json())
                .then(data => {
                    updateTable(data)
                })
                .catch(error => console.error(error))
        }else{
            GetAllAppointments();
        }
    })
}

document.addEventListener("DOMContentLoaded", function(e){

    ValidateSession()

    GetAllAppointments()

    SearchPatientsForm()

    AddNewAppointment()

    BarsMenu()

    Filter()

    FlickerInput()

    SearchPatient()

    WelcomeUsername()

    CloseSession()

})