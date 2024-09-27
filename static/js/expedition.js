window.addEventListener('DOMContentLoaded', (event) => {
    // Votre code JavaScript existant ici
    // Fonction pour afficher la fenêtre modale
    function afficherModal() {
        const modal = document.getElementById('modal');
        modal.style.display = 'block';
    }
    // Fonction pour cacher la fenêtre modale
    function cacherModal() {
        const modal = document.getElementById('modal');
        modal.style.display = 'none';
    }
    // Écouteur d'événement pour le clic sur le bouton "Créer un colis"
    const btnCreerColis = document.getElementById('btn-creer-colis');
    btnCreerColis.addEventListener('click', afficherModal);
    // Écouteur d'événement pour le clic sur le bouton "Fermer" de la fenêtre modale
    const btnFermerModal = document.getElementById('btn-fermer-modal');
    btnFermerModal.addEventListener('click', cacherModal);
    const paysInput = document.getElementById('point-depart');
    const suggestions = document.getElementById('suggestions');
    const paysArriveeInput = document.getElementById('point-arrivee');
    const suggestionsArrivee = document.getElementById('suggestions-arrivee');
    let allCountries = [];
    fetch("https://restcountries.com/v3.1/all")
        .then(response => response.json())
        .then(data => {
            allCountries = data.map(country => ({
                name: country.name.common,
                flag: country.flags.png
            }));
            // Trier les pays par ordre alphabétique
            allCountries.sort((a, b) => a.name.localeCompare(b.name));
        })
        .catch(error => console.error("Erreur lors de la récupération des pays :", error))
        .finally(() => {
            // Une fois que les données sont chargées, activer l'écoute des événements de saisie et de clic
            paysInput.addEventListener('input', () => filtrerPays('point-depart', 'suggestions'));
            paysInput.addEventListener('focus', () => filtrerPays('point-depart', 'suggestions'));
            paysArriveeInput.addEventListener('input', () => filtrerPays('point-arrivee', 'suggestions-arrivee'));
            paysArriveeInput.addEventListener('focus', () => filtrerPays('point-arrivee', 'suggestions-arrivee'));
        });
    // Fonction pour filtrer les suggestions en fonction de la saisie
    function filtrerPays(inputId, suggestionsId) {
        const saisie = document.getElementById(inputId).value.toLowerCase();
        const suggestionsElement = document.getElementById(suggestionsId);
        const suggestionsFiltrees = allCountries.filter(pays =>
            pays.name.toLowerCase().startsWith(saisie)
        );
        afficherSuggestions(suggestionsElement, suggestionsFiltrees, inputId);
    }
    // Fonction pour afficher les suggestions
    function afficherSuggestions(suggestionsElement, suggestionsFiltrees, inputId) {
        suggestionsElement.innerHTML = '';
        suggestionsFiltrees.forEach(pays => {
            const suggestion = document.createElement('a');
            suggestion.href = '#';
            const drapeau = document.createElement('img');
            drapeau.src = pays.flag;
            suggestion.appendChild(drapeau);
            const nomPays = document.createElement('span');
            nomPays.textContent = pays.name;
            suggestion.appendChild(nomPays);
            suggestion.addEventListener('click', () => {
                document.getElementById(inputId).value = pays.name;
                suggestionsElement.style.display = 'none';
            });
            suggestionsElement.appendChild(suggestion);
        });
        if (suggestionsFiltrees.length > 0) {
            suggestionsElement.style.display = 'block';
        } else {
            suggestionsElement.style.display = 'none';
        }
    }
});
function previewImage(event) {
    var input = event.target;
    var reader = new FileReader();
    reader.onload = function(){
        var dataURL = reader.result;
        var preview = document.getElementById('preview');
        preview.src = dataURL;
        preview.style.display = "block"; // Afficher l'image prévisualisée
    };
    reader.readAsDataURL(input.files[0]);
}
document.addEventListener("DOMContentLoaded", function() {
    var destinataireInput = document.getElementById("destinataire");
    var destinationInput = document.getElementById("destination");
    destinataireInput.addEventListener("input", function() {
        var destinataire = destinataireInput.value.trim();
        if (destinataire === "") {
            destinationInput.value = "";
            return;
        }
        var xhr = new XMLHttpRequest();
        xhr.open("GET", "/get-location?destinataire=" + encodeURIComponent(destinataire), true);
        xhr.onreadystatechange = function() {
            if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {
                var response = JSON.parse(xhr.responseText);
                if (response.success) {
                    destinationInput.value = response.location;
                } else {
                    console.error("Failed to retrieve location:", response.error);
                }
            }
        };
        xhr.send();
    });
});