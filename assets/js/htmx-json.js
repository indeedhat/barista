htmx.defineExtension('json-enc', (() => {
    const buildObject  = (ob, params) => k => {
        const values = params.getAll(k)

        let key = k
        let int = false
        let rest = []

        if (key.endsWith(".int")) {
            int = true
            key = key.slice(0, -4)
        }

        if (key.includes(".")) {
            [ key, ...rest ] = key.split(".")
        }

        if (values.length == 1 || !key.endsWith('[]')) {
            const v = int ? ~~values[0] : values[0]

            if (rest.length) {
                if (key in ob) {
                    ob[key][rest.join(".")] = v
                } else {
                    ob[key] = { [rest.join(".")]: v }
                }
            } else {
                ob[key] = v
            }
            return
        }

        if (key.endsWith('[]')) {
            key = key.slice(0, -2)
        }
        if (!(key in ob)) {
            ob[key] = []
        }

        for (let i in values) {
            const v = int ? ~~values[i] : values[i]

            if (rest.length) {
                if (i < ob[key].length) {
                    ob[key][i][rest.join(".")] = v
                } else {
                    ob[key].push({ [rest.join(".")]: v })
                }
            } else {
                ob[key].push(v)
            }
        }
    }

    return {
        onEvent: function (name, evt) {
            if (name === "htmx:configRequest") {
                evt.detail.headers['Content-Type'] = "application/json"
            }
        },
        encodeParameters : function(xhr, parameters, _) {
            xhr.overrideMimeType('text/json')
            let ob = {}

            parameters.keys().forEach(buildObject(ob, parameters))

            return JSON.stringify(ob)
        }
    }
})())
