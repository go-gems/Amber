document.addEventListener("DOMContentLoaded", () => {
    document.querySelector("#shorten-button").onclick = () => {
        let url
        try {
            url = new URL(document.querySelector("#url-input").value)
        } catch (e) {
            showError("Your URL is malformed :( Make sure you specify the scheme.")
            return
        }
        fetch("/links", {
            method: 'POST',
            body: url.toString()
        })
            .then(data => {
                data.json()
                    .then(json => {
                        copyToClipboard(json.url)
                        alert(json.url)
                    })
                    .catch(() => {
                        showError("An unknown error as occured :(")
                    })
            })
            .catch(() => {
                showError("An unknown error as occured :(")
            })
    }
})

function showError(message) {
    document.querySelector("#error-message-content").textContent = message
    const errorMessage = document.querySelector("#error-message")
    errorMessage.style.transform = "translateX(0)"
    setTimeout(() => {
        errorMessage.style.transform = "translateX(600px)"
    }, 4000)
}

function copyToClipboard(text) {
    if (window.clipboardData && window.clipboardData.setData) {
        // Internet Explorer-specific code path to prevent textarea being shown while dialog is visible.
        return clipboardData.setData("Text", text);

    }
    else if (document.queryCommandSupported && document.queryCommandSupported("copy")) {
        var textarea = document.createElement("textarea");
        textarea.textContent = text;
        textarea.style.position = "fixed";  // Prevent scrolling to bottom of page in Microsoft Edge.
        document.body.appendChild(textarea);
        textarea.select();
        try {
            return document.execCommand("copy");  // Security exception may be thrown by some browsers.
        }
        catch (ex) {
            console.warn("Copy to clipboard failed.", ex);
            return false;
        }
        finally {
            document.body.removeChild(textarea);
        }
    }
}