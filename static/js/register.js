function verifierComplexiteMotDePasse() {
    var motDePasse = document.getElementById("password").value;
    var indicateur = document.getElementById("password-strength");
    var pastille = document.getElementById("password-color-indicator");

    if (motDePasse.length === 0) {
        indicateur.innerText = "";
        pastille.className = "color-indicator";
        return;
    }

    var complexite = 0;
    if (motDePasse.length >= 8) complexite++;
    if (motDePasse.match(/[A-Z]/)) complexite++;
    if (motDePasse.match(/[a-z]/)) complexite++;
    if (motDePasse.match(/[0-9]/)) complexite++;
    if (motDePasse.match(/[^A-Za-z0-9]/)) complexite++;

    if (complexite > 4) {
        indicateur.innerText = "Fort";
        pastille.className = "color-indicator green";
        indicateur.className = "strength-text strong";
    } else if (complexite >= 2) {
        indicateur.innerText = "Moyen";
        pastille.className = "color-indicator orange";
        indicateur.className = "strength-text medium";
    } else {
        indicateur.innerText = "Faible";
        pastille.className = "color-indicator red";
        indicateur.className = "strength-text weak";
    }
}

function togglePasswordVisibility() {
    var passwordInput = document.getElementById("password");
    var toggleButton = document.getElementById("toggle-password");
    if (passwordInput.type === "password") {
        passwordInput.type = "text";
        toggleButton.innerText = "🙉";
    } else {
        passwordInput.type = "password";
        toggleButton.innerText = "🙈";
    }
}
document.getElementById("password").addEventListener("input", verifierComplexiteMotDePasse);
