let confirmOpen = false

document.addEventListener("DOMContentLoaded", () => {
    document.querySelector("#shorten-button").onclick = () => {
        if (confirmOpen) {
            return
        }
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
                        showConfirmBox(json.url)
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

function showConfirmBox(url) {
    confirmOpen = true

    const html = `
    <div id="confirm-message" class="absolute inset-0 shadow-lg rounded-2xl p-4 bg-white dark:bg-gray-800 w-80 h-80 z-10 m-auto">
        <div class="w-full h-full text-center">
            <div class="flex h-full flex-col justify-between">
                <svg class="h-12 w-12 mt-4 m-auto text-green-500" stroke="currentColor" fill="none" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7">
                    </path>
                </svg>
                <p class="text-gray-600 dark:text-gray-100 text-md py-2 px-6">
                    There you go !
                    <br/>
                    <span class="text-gray-800 dark:text-white font-bold">
                        <a href="${url}">${url}</a>
                    </span>
                    <br/>
                    The link has been pasted in you clipboard
                </p>
                <div class="flex items-center justify-between gap-4 w-full mt-8">
                    <button onclick="closeConfirmBox()" type="button" class="py-2 px-4  bg-indigo-600 hover:bg-indigo-700 focus:ring-indigo-500 focus:ring-offset-indigo-200 text-white w-full transition ease-in duration-200 text-center text-base font-semibold shadow-md focus:outline-none focus:ring-2 focus:ring-offset-2  rounded-lg ">
                        Close
                    </button>
                </div>
            </div>
        </div>
    </div>`

    document.body.insertAdjacentHTML("beforeend", html)

}

function closeConfirmBox() {
    document.querySelector("#confirm-message").remove()
    confirmOpen = false
}