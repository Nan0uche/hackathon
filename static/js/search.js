function redirect() {
    var trackingId = document.getElementById("trackingId").value;
    if (/^[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{4}$/.test(trackingId)) {
        window.location.href = "suivi?id=" + trackingId;
    } else {
        alert("Veuillez entrer un ID de colis valide au format XXXX-XXXX-XXXX.");
    }
}