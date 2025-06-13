document.addEventListener("DOMContentLoaded", () => {
    document.body.addEventListener("triggerToast", e => {
        const $toast = document.createElement("div")
        $toast.className = `alert alert-${e.detail.level}`

        const $body = document.createElement("span")
        $body.textContent = e.detail.message

        appendIcon($toast, e.detail.level)
        $toast.appendChild($body)

        let timeout = setTimeout(() => $toast.remove(), 10_000)
        $toast.addEventListener("click", () => {
            clearTimeout(timeout)
            $toast.remove()
        })

        document.getElementById("toast-container").appendChild($toast)
    })
});

const appendIcon = ($el, level) => {
    const $svg = document.createElement("svg")
    $svg.setAttribute("xmlns", "http://www.w3.org/2000/svg")
    $svg.setAttribute("fill", "none")
    $svg.setAttribute("viewBox", "0 0 24 24")
    $svg.className = "h-6 w-6 shrink-0 stroke-current"

    const $stroke = document.createElement("path")
    $stroke.setAttribute("stroke-linecap", "round")
    $stroke.setAttribute("stroke-linejoin", "round")
    $stroke.setAttribute("stroke-width", "2")
    $svg.appendChild($stroke)

    switch (level) {
    case 'success':
        $stroke.setAttribute("d", "M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z")
        break
    case 'error':
        $stroke.setAttribute("d", "M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z")
        break
    case 'info':
        $stroke.setAttribute("d", "M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z")
        break
    }

    $el.appendChild($svg)
}
