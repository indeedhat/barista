document.addEventListener("DOMContentLoaded", () => {
    document.body.addEventListener("triggerToast", e => {
        let timeout

        const $toast = document.createElement("article")
        $toast.className = `alert alert-${e.detail.level}`
        $toast.addEventListener("click", () => {
            clearTimeout(timeout)
            $toast.remove()
        })

        const $body = document.createElement("span")
        $body.textContent = e.detail.message
        $toast.appendChild($body)

        document.getElementById("toast-container").appendChild($toast)

        timeout = setTimeout(() => $toast.remove(), 10_000)
    })
});
