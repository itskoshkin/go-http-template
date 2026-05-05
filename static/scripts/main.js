function goHome() { /* Go to home page */
    window.location.href = '/'; /* Navigate to root */
}

document.addEventListener('DOMContentLoaded', () => { /* Bind handlers after page markup is ready. */
    const homeButton = document.querySelector('.error-btn-home'); /* Find the error page home button. */
    if (!homeButton) { /* Ignore pages without the home button. */
        return; /* Stop before binding the click handler. */
    }
    homeButton.addEventListener('click', goHome); /* Navigate to root when the home button is clicked. */
});
