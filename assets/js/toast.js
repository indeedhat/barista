(() => {
    document.addEventListener("DOMContentLoaded", () => {
        document.body.addEventListener("triggerToast", e => {
            const $toast = createElement(`
                <div role="alert" class="alert">
                    ${createIcon(e.detail.level)}
                    <span>${e.detail.message}</span>
                </div>
            `)

            let timeout = setTimeout(() => $toast.remove(), 10_000)
            $toast.addEventListener("click", () => {
                clearTimeout(timeout)
                $toast.remove()
            })

            document.getElementById("toast-container").appendChild($toast)
            $toast.classList.add(`alert-${e.detail.level}`)
        })
    });

    const createIcon = (level) => {
        let stroke
        switch (level) {
        case 'success':
            stroke = "M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
            break
        case 'error':
            stroke = "M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
            break
        case 'info':
            stroke = "M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            break
        }
        return `
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="h-6 w-6 shrink-0 stroke-current">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="${stroke}"></path>
            </svg>
        `;
    }

    const createElement = (html) => {
        console.log(html)
        const $wrapper = document.createElement("template")
        $wrapper.className = "hidden"
        $wrapper.innerHTML = html
        return $wrapper.content.firstElementChild
    }
})();
